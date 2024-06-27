// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import { IDisputeGame } from "./interfaces/IDisputeGame.sol";
import { IBlockDisputeGame }  from "./interfaces/IBlockDisputeGame.sol";
import { IDelayedWETH } from "./interfaces/IDelayedWETH.sol";
import { IInitializable } from "./interfaces/IInitializable.sol";
import { IOptimisticZKGame } from "./interfaces/IOptimisticZKGame.sol";
import { IDisputeGameFactory } from "./interfaces/IDisputeGameFactory.sol";

import "src/dispute/lib/Types.sol";

import { Clone } from "@solady/utils/Clone.sol";
import { SP1Verifier } from "@sp1-contracts/SP1Verifier.sol";

contract OptimisticZKGame is IOptimisticZKGame, Clone, SP1Verifier {
    using LibClock for Clock;
    using LibDuration for Duration;
    using LibTimestamp for Timestamp;

    ////////////////////////////////////////////////////////////////
    //                         State Vars                         //
    ////////////////////////////////////////////////////////////////

    /// @notice An ID used for the proposer's initial bond, since it doesn't fall into a specific challenge.
    uint constant GLOBAL_CHALLENGE_ID = type(uint64).max;

    /// @notice The initial bond required to start a game.
    uint constant INITIAL_BOND = 1 ether;

    /// @notice The bond required for each bisection.
    uint constant BISECTION_BOND = 0.1 ether;

    /// @notice The DisputeGameFactory contract.
    IDisputeGameFactory immutable FACTORY;

    /// @notice The game type ID.
    GameType immutable GAME_TYPE;

    /// @notice The maximum duration of the entire game.
    /// @dev If the game has not resolved against the propoer by this time, the proposer wins.
    Duration immutable MAX_GAME_DURATION;

    /// @notice The maximum duration of the sum of the proposer's turns in this game.
    /// @dev If the proposer's clock runs out, the challenger wins the game.
    Duration immutable MAX_PROPOSER_DURATION;

    /// @notice The DelayedWETH contract.
    /// @dev Used as a temporary replacement for WETH to clawback bonds in the case of a fraudulent game.
    IDelayedWETH immutable WETH;

    /// @notice The block number at which Ecotone was deployed.
    /// @dev Blocks before this number cannot be used as a reference block for proving against.
    uint immutable ECOTONE_ORIGIN_BLOCK;

    /// @notice The largest allowed gap between the reference block and the block being proved.
    uint immutable MAX_BLOCK_GAP;

    /// @notice The verification key used by the SP1Verifier contract.
    bytes32 immutable VKEY;

    /// @notice The timestamp at which the game was created.
    Timestamp public createdAt;

    /// @notice The timestamp at which the game was resolved.
    Timestamp public resolvedAt;

    /// @notice The previously confirmed root that the game is starting from.
    OutputRoot anchorStateRoot;

    /// @notice The status of the game (in progress, challenger wins, or defender wins).
    GameStatus public status;

    /// @notice The challenges being played on the game.
    Challenge[] challenges;

    /// @notice Funds that have are pending withdrawal from the DelayedWETH contract.
    mapping(address => uint) credits;

    ////////////////////////////////////////////////////////////////
    //                           SETUP                            //
    ////////////////////////////////////////////////////////////////

    /// @param _factory The DisputeGameFactory contract.
    /// @param _gameType The game type ID.
    /// @param _maxGameDuration The maximum duration of the entire game.
    /// @param _maxProposerDuration The maximum duration of the sum of the proposer's turns in this game.
    /// @param _weth The DelayedWETH contract.
    /// @param _originBlock The earliest block number that can be proved against.
    /// @param _maxBlockGap The largest allowed gap between the reference block and the block being proved.
    /// @param _vkey The verification key used by the SP1Verifier contract.
    constructor(
        address _factory,
        GameType _gameType,
        Duration _maxGameDuration,
        Duration _maxProposerDuration,
        IDelayedWETH _weth,
        uint _originBlock,
        uint _maxBlockGap,
        bytes32 _vkey
    ) {
        // The challenger's time will be equal to _maxGameDuration - _maxProposerDuration, so ensure that it is equal to the proposer's time.
        if (_maxGameDuration.raw() != _maxProposerDuration.raw() * 2) revert InvalidDurations();

        // Set all the immutable values in the implementation contract.
        FACTORY = IDisputeGameFactory(_factory);
        GAME_TYPE = _gameType;
        MAX_GAME_DURATION = _maxGameDuration;
        MAX_PROPOSER_DURATION = _maxProposerDuration;
        WETH = _weth;
        ECOTONE_ORIGIN_BLOCK = _originBlock;
        MAX_BLOCK_GAP = _maxBlockGap;
        VKEY = _vkey;
    }

    /// @inheritdoc IInitializable
    function initialize() external payable {
        if (msg.sender != address(FACTORY)) revert Unauthorized();

        // Query the factory to retrieve a game that has already been settled in favor of the defender.
        (,, IDisputeGame prevGame) = FACTORY.gameAtIndex(startingRootGameIndex());
        if (prevGame.status() != GameStatus.DEFENDER_WINS) revert PreviousGameNotResolved();

        // Pull the L2 block number from the previous game.
        uint prevL2BlockNumber = IBlockDisputeGame(address(prevGame)).l2BlockNumber();

        // Set the anchorStateRoot to the previous game's root.
        anchorStateRoot = OutputRoot({ root: Hash.wrap(prevGame.rootClaim().raw()), l2BlockNumber: prevL2BlockNumber });

        // Only allow proving against blocks since Ecotone.
        // This is because the state transition function used in the final ZK proof only computed Ecotone blocks.
        if (anchorStateRoot.l2BlockNumber < ECOTONE_ORIGIN_BLOCK) revert InvalidBlockNumber();

        // Require that the block being proven comes after the reference block.
        if (anchorStateRoot.l2BlockNumber >= l2BlockNumber()) revert InvalidBlockNumber();

        // Require that the gap between the reference block and the block being proven is less than MAX_BLOCK_GAP.
        // This is used to avoid attacks where such a large gap is used that the game cannot be completed.
        if (anchorStateRoot.l2BlockNumber + MAX_BLOCK_GAP < l2BlockNumber()) revert InvalidBlockNumber();

        // Deposit the initial bond into the contract.
        _depositBond(GLOBAL_CHALLENGE_ID);

        // Set the game's creation timestamp to the time of initialization.
        createdAt = Timestamp.wrap(uint64(block.timestamp));
    }

    ////////////////////////////////////////////////////////////////
    //                         BISECTION                          //
    ////////////////////////////////////////////////////////////////

    function createNewChallenge(uint _blockNum, bytes32 _proposedRoot) public payable {
        // The current length will be the ID after a new challenge is pushed to the array.
        uint challengeId = challenges.length;

        // Don't allow proposer to challenge themselves.
        // This is needed so that we can trust that the proposer address winning means the root is valid.
        if (msg.sender == gameCreator()) revert ProposerIsChallenger();

        // Create a new challenge.
        challenges.push(Challenge({
            challenger: msg.sender,
            proposerClock: LibClock.wrap(Duration.wrap(0), Timestamp.wrap(0)),
            left: IntermediateRootClaim({
                outputRoot: anchorStateRoot,
                status: IntermediateClaimStatus.ACCEPTED
            }),
            right: IntermediateRootClaim({
                outputRoot: OutputRoot({ root: Hash.wrap(rootClaim().raw()), l2BlockNumber: l2BlockNumber() }),
                status: IntermediateClaimStatus.ACCEPTED
            }),
            current: OutputRoot({ root: Hash.wrap(bytes32(0)), l2BlockNumber: 0 }),
            totalBonds: 0,
            resolved: false,
            l1Head: blockhash(block.number - 1)
        }));

        // Perform the first split on this new challenge.
        split(challengeId, _blockNum, _proposedRoot);
    }

    /// @param _blockNum The block number that the proposer is proposing a root for.
    /// @param _proposedRoot The root that the proposer is proposing for the given block number.
    function split(uint _challengeId, uint _blockNum, bytes32 _proposedRoot) public payable {
        // Use the _challengeId to access the correct challenge struct.
        if (_challengeId >= challenges.length) revert InvalidChallengeId();
        Challenge memory challenge = challenges[_challengeId];

        // Only allow the challenge creator to act on this challenge.
        if (msg.sender != challenge.challenger) revert Unauthorized();

        // Only allow the challenger to call this function if it is their turn (ie the proposer timestamp is unset).
        if (challenge.proposerClock.timestamp().raw() != 0) revert WrongTurn();

        // Require this to be called at the bisection point between the left and right roots.
        // Note that it is required for this to be an argument to avoid a frontrunning attack where `nextSplitBlock()`
        // is shifted so that the challenger attests to the proposed root for the wrong block.
        if (_blockNum != nextSplitBlock(_challengeId)) revert InvalidBlockNumber();

        // If the left and right roots are already adjacent, `nextSplitBlock()` will return
        // the left root. In this case, there is nothing to split and the game is ready to prove.
        if (_blockNum == challenge.left.outputRoot.l2BlockNumber) revert ReadyToProve();

        // Deposit the bond.
        _depositBond(_challengeId);

        // Create a new bisection point for the proposal to respond to.
        challenge.current = OutputRoot({
            l2BlockNumber: _blockNum,
            root: Hash.wrap(_proposedRoot)
        });

        // Keep the proposer's new duration the same, and set the timestamp to now to start their turn.
        challenge.proposerClock = LibClock.wrap(challenge.proposerClock.duration(), Timestamp.wrap(uint64(block.timestamp)));

        // Update the challenge in storage.
        challenges[_challengeId] = challenge;
    }

    /// @param _challengeId The ID of the challenge to respond to.
    /// @param _accepted Whether the proposer accepts the proposed root or rejects it.
    function respondToSplit(uint _challengeId, bool _accepted) public payable {
        // Only allow the proposer to respond to splits.
        if (msg.sender != gameCreator()) revert Unauthorized();

        // Use the _challengeId to access the correct challenge struct.
        if (_challengeId >= challenges.length) revert InvalidChallengeId();
        Challenge memory challenge = challenges[_challengeId];

        // Only allow the proposer to call this function if it is their turn (ie their timestamp is set).
        if (challenge.proposerClock.timestamp().raw() == 0) revert WrongTurn();

        // Deposit the bond.
        _depositBond(_challengeId);

        // If accepted, move the left root to the proposed root.
        if (_accepted) {
            challenge.left = IntermediateRootClaim({
                outputRoot: challenge.current,
                status: IntermediateClaimStatus.ACCEPTED
            });

        // If disputed, move the right root to the proposed root.
        } else {
            challenge.right = IntermediateRootClaim({
                outputRoot: challenge.current,
                status: IntermediateClaimStatus.CHALLENGED
            });
        }

        // Calculate the proposer's new duration by incorporating any time since the latest timestamp.
        // Set the proposer's timestamp to 0 to stop their turn.
        Duration newDuration = totalClockDuration(challenge.proposerClock);
        challenge.proposerClock = LibClock.wrap(newDuration, Timestamp.wrap(0));

        // Update the challenge in storage.
        challenges[_challengeId] = challenge;
    }

    ////////////////////////////////////////////////////////////////
    //                        RESOLUTION                          //
    ////////////////////////////////////////////////////////////////

    function proveStep(uint _challengeId, bytes memory _proofBytes, PublicValuesStruct memory _publicValues) public payable {
        // Use the _challengeId to access the correct challenge struct.
        if (_challengeId >= challenges.length) revert InvalidChallengeId();
        Challenge memory challenge = challenges[_challengeId];

        // Only allow the challenge creator to act on this challenge.
        if (msg.sender != challenge.challenger) revert Unauthorized();

        // Don't allow this to be called if (a) challenger has already won a different proof or (b) game has already resolved for proposer.
        if (status != GameStatus.IN_PROGRESS) revert GameAlreadyResolved();

        // Require that the left and right roots are adjacent and ready to prove.
        if (challenge.left.outputRoot.l2BlockNumber + 1 != challenge.right.outputRoot.l2BlockNumber) revert NotReadyToProve();

        // Validate public values passed to the verifier...

        // 1) The real left root of the game matches the passed l2PreRoot.
        if (challenge.left.outputRoot.root.raw() != _publicValues.l2PreRoot) revert InvalidPublicInput();

        // 2) The real L1 block hash at challenge time matches the passed L1 block hash.
        if (challenge.l1Head != _publicValues.l1Head) revert InvalidPublicInput();

        // 3) The real L2 block number being proven matches the passed L2 block number.
        if (challenge.right.outputRoot.l2BlockNumber != _publicValues.l2BlockNumber) revert InvalidPublicInput();

        if (challenge.right.status == IntermediateClaimStatus.CHALLENGED) {
            // 4a) If the right root has been challenged by the proposer, the challenger must prove that we CAN transition from left to right.
            // Therefore, prove that the real right root matches the passed l2PostRoot.
            if (challenge.right.outputRoot.root.raw() != _publicValues.l2PostRoot) revert InvalidPublicInput();
            verifyProof(VKEY, abi.encode(_publicValues), _proofBytes);
        } else {
            // 4b) If the right root is ACCEPTED, it means nothing has been challenged.
            // The proposer is claiming that left (proposed block minus 1) DOES transition to right (proposed block).
            // Therefore, the challenger must prove a block with an l2PostRoot that does NOT match the right root in the game.
            if (challenge.right.outputRoot.root.raw() == _publicValues.l2PostRoot) revert InvalidPublicInput();
            verifyProof(VKEY, abi.encode(_publicValues), _proofBytes);
        }

        // Once the proof has been completed, resolve the game.
        // The challenger who proves this step gets their own challenge's bond, plus the proposer's global bond.
        uint[] memory _challengeIds = new uint[](2);
        _challengeIds[0] = GLOBAL_CHALLENGE_ID;
        _challengeIds[1] = _challengeId;
        _resolveInternal(_challengeIds, msg.sender);
    }

    /// @notice This function is called to end the game in favor of the proposer if they were not successfully challenged.
    /// @dev Unlike the Fraud Proof Game, this can't be called to settle to CHALLENGER_WINS, only DEFENDER_WINS.
    function resolve() public returns (GameStatus status_) {
        // We can only resolve an in progress game where the full game clock has run out.
        // If the proposer clock had run out or a ZK proof succeeded, the status would have been updated.
        // This implies we are ready for the proposer to win the game.
        if (status != GameStatus.IN_PROGRESS) revert GameAlreadyResolved();
        if (block.timestamp < createdAt.raw() + MAX_GAME_DURATION.raw()) revert ClockNotExpired();

        // The proposer wins all games, so create an array of all of them (including global).
        uint[] memory cIds = new uint[](challenges.length + 1);
        for (uint i; i < challenges.length; i++) {
            cIds[i] = i;
        }
        cIds[challenges.length] = GLOBAL_CHALLENGE_ID;

        // Resolve all challenges in favor of the proposer, and update the game status.
        _resolveInternal(cIds, gameCreator());

        // Return the resulting game status.
        return status;
    }

    /// @notice This function is called by the challenger to end the game if (a) the proposer's clock has run out
    ///         or (b) another challenge has already resolved in the challenger's favor.
    function resolveChallenge(uint _challengeId) public {
        // Use the _challengeId to access the correct challenge struct.
        if (_challengeId >= challenges.length) revert InvalidChallengeId();
        Challenge memory challenge = challenges[_challengeId];

        // Do not allow the same challenge to be resolved twice.
        if (challenge.resolved) revert ChallengeAlreadyResolved();

        // If this is the first challenge in which the proposer's clock has run out,
        // resolve the game in favor of the challenger (and also reward them with the global bond).
        if (
            totalClockDuration(challenge.proposerClock).raw() > MAX_PROPOSER_DURATION.raw() &&
            status == GameStatus.IN_PROGRESS
        ) {
            uint[] memory _challengeIds = new uint[](2);
            _challengeIds[0] = GLOBAL_CHALLENGE_ID;
            _challengeIds[1] = _challengeId;
            _resolveInternal(_challengeIds, challenge.challenger);

        // If another challenge has resolved for the challenger, automatically allow the game to
        // resolve in challenger's favor.
        } else if (status == GameStatus.CHALLENGER_WINS) {
            uint[] memory _challengeIds = new uint[](1);
            _challengeIds[0] = _challengeId;
            _resolveInternal(_challengeIds, challenge.challenger);

        // If the game is not resolved for the challenger and the clock is not expired, revert.
        } else {
            revert ChallengeCantBeResolved();
        }
    }

    function _resolveInternal(uint[] memory _challengeIds, address _recipient) internal {
        // Calculate the total amount to distribute to the recipient by summing the bonds of all challengeIds.
        uint amountToDistribute;
        for (uint i; i < _challengeIds.length; i++) {
            uint cId = _challengeIds[i];
            amountToDistribute += _calculateBonds(cId);

            // Mark each challenge as resolved so it can't be called again.
            if (cId != GLOBAL_CHALLENGE_ID) challenges[cId].resolved = true;
        }

        // Request a delayed withdrawal from DelayedWETH and add withdrawal credits to the recipient.
        credits[_recipient] += amountToDistribute;
        WETH.unlock(_recipient, amountToDistribute);

        // If this is the first challenge being resolved, update the game status and resolvedAt timestamp.
        if (status == GameStatus.IN_PROGRESS) {
            status = _recipient == gameCreator() ? GameStatus.DEFENDER_WINS : GameStatus.CHALLENGER_WINS;
            resolvedAt = Timestamp.wrap(uint64(block.timestamp));
        }
    }

    ////////////////////////////////////////////////////////////////
    //                           BONDS                            //
    ////////////////////////////////////////////////////////////////

    // Deposit the bond into the DelayedWETH contract and increment the total bonds for the given challenge.
    function _depositBond(uint _challengeId) internal {
        if (msg.value != getRequiredBond(_challengeId)) revert WrongBondAmount();
        if (_challengeId != GLOBAL_CHALLENGE_ID) challenges[_challengeId].totalBonds += msg.value;
        WETH.deposit{value: msg.value}();
    }

    // Calculate the amount that should be distributed to the winner.
    function _calculateBonds(uint _challengeId) internal view returns (uint) {
        if (_challengeId == GLOBAL_CHALLENGE_ID) return getRequiredBond(GLOBAL_CHALLENGE_ID);

        return challenges[_challengeId].totalBonds;
    }

    /// @param _recipient The user to claim the ETH for.
    function claimCredit(address _recipient) external {
        // Cache the number of credits the user has and reset their local balance to 0.
        uint256 recipientCredit = credits[_recipient];
        credits[_recipient] = 0;

        // Request the funds from the DelayedWETH contract.
        // Note: This will revert if it has not been DELAY_SECONDS since WETH.unlock() was called.
        WETH.withdraw(_recipient, recipientCredit);

        // Forward the funds to the recipient.
        (bool success,) = _recipient.call{ value: recipientCredit }(hex"");
        if (!success) revert TransferFailed();
    }

    ////////////////////////////////////////////////////////////////
    //                           VIEWS                            //
    ////////////////////////////////////////////////////////////////

    /// @param _challengeId The ID of the challenge to get the next split block for.
    /// @return blockNum The block number that the proposer should propose a root for next.
    function nextSplitBlock(uint _challengeId) public view returns (uint blockNum) {
        Challenge memory challenge = challenges[_challengeId];

        // If it's the proposer's turn, it returns the L2 block that they currently need to evaluate.
        if (challenge.proposerClock.timestamp().raw() != 0) return challenge.current.l2BlockNumber;

        // If it's the challenger's turn, it bisects left and right to return the block they should next split on.
        return (challenge.left.outputRoot.l2BlockNumber + challenge.right.outputRoot.l2BlockNumber) / 2;
    }

    /// @param _clock The clock to calculate the total current duration for.
    /// @return duration_ The total duration that the clock has been running for (including pending time on the current turn).
    function totalClockDuration(Clock _clock) public view returns (Duration duration_) {
        // If timestamp == 0, it's not the user's turn, so just return the duration.
        if (_clock.timestamp().raw() == 0) {
            duration_ = _clock.duration();

        // Otherwise, add the time on the current turn to the duration to get the total.
        } else {
            uint timeOnCurrentMove = block.timestamp - _clock.timestamp().raw();
            duration_ = Duration.wrap(uint64(_clock.duration().raw() + timeOnCurrentMove));
        }
    }

    /// @param _challengeId The ID of the challenge to get the required bond for.
    /// @return requiredBond_ The amount of ETH required to bond for the given challenge.
    /// @dev It is important that the sum of totalBonds on a challenge is enough to justify the ZK work.
    ///      Otherwise, the proposer could frontrun the proof to claim the larger bond back.
    function getRequiredBond(uint _challengeId) public pure returns (uint requiredBond_) {
        requiredBond_ = _challengeId == GLOBAL_CHALLENGE_ID ? INITIAL_BOND : BISECTION_BOND;
    }

    /// @inheritdoc IDisputeGame
    function gameType() public view override returns (GameType gameType_) {
        gameType_ = GAME_TYPE;
    }

    /// @inheritdoc IDisputeGame
    function gameCreator() public pure returns (address creator_) {
        creator_ = _getArgAddress(0x00);
    }

    /// @inheritdoc IDisputeGame
    function rootClaim() public pure returns (Claim rootClaim_) {
        rootClaim_ = Claim.wrap(_getArgBytes32(0x14));
    }

    /// @inheritdoc IDisputeGame
    function l1Head() public pure returns (Hash l1Head_) {
        l1Head_ = Hash.wrap(_getArgBytes32(0x34));
    }

    /// @inheritdoc IDisputeGame
    function extraData() public pure returns (bytes memory extraData_) {
        // The extra data starts at the second word within the cwia calldata and
        // is 32 bytes long.
        extraData_ = _getArgBytes(0x54, 0x60);
    }

    /// @inheritdoc IDisputeGame
    function gameData() external view returns (GameType gameType_, Claim rootClaim_, bytes memory extraData_) {
        gameType_ = gameType();
        rootClaim_ = rootClaim();
        extraData_ = extraData();
    }


    /// @return startingRootGameIndex_ The index of the validated game we are proving from.
    function startingRootGameIndex() public pure returns (uint256 startingRootGameIndex_) {
        startingRootGameIndex_ = _getArgUint256(0x74);
    }

    /// @inheritdoc IBlockDisputeGame
    function l2BlockNumber() public pure returns (uint256 l2BlockNumber_) {
        l2BlockNumber_ = _getArgUint256(0x54);
    }


    /// @notice Starting output root and block number of the game.
    /// @return startingRoot_ The root that the game claiming is proving.
    /// @return l2BlockNumber_ The block number that the game claiming is proving.
    function startingOutputRoot() external view returns (Hash startingRoot_, uint256 l2BlockNumber_) {
        startingRoot_ = Hash.wrap(anchorStateRoot.root.raw());
        l2BlockNumber_ = anchorStateRoot.l2BlockNumber;
    }

    /// @notice Only the starting block number of the game.
    /// @return startingBlockNumber_ The block number that the game claiming is proving.
    function startingBlockNumber() external view returns (uint256 startingBlockNumber_) {
        startingBlockNumber_ = anchorStateRoot.l2BlockNumber;
    }

    /// @notice Only the starting output root of the game.
    /// @return startingRootHash_ The root that the game claiming is proving.
    function startingRootHash() external view returns (Hash startingRootHash_) {
        startingRootHash_ = Hash.wrap(anchorStateRoot.root.raw());
    }
}
