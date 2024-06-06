// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import { IDisputeGame } from "./interfaces/IDisputeGame.sol";
import { IBlockDisputeGame }  from "./interfaces/IBlockDisputeGame.sol";
import { IDelayedWETH } from "./interfaces/IDelayedWETH.sol";
import { IInitializable } from "./interfaces/IInitializable.sol";
import { IPureZKGame } from "./interfaces/IPureZKGame.sol";
import { IDisputeGameFactory } from "./interfaces/IDisputeGameFactory.sol";

import "src/dispute/lib/Types.sol";

import { Clone } from "@solady/utils/Clone.sol";
import { SP1Verifier } from "@sp1-contracts/SP1Verifier.sol";

contract PureZKGame is IPureZKGame, Clone, SP1Verifier {

    ////////////////////////////////////////////////////////////////
    //                         State Vars                         //
    ////////////////////////////////////////////////////////////////

    /// @notice The DisputeGameFactory contract.
    IDisputeGameFactory immutable FACTORY;

    /// @notice The game type ID.
    GameType immutable GAME_TYPE;

    /// @notice The block number at which Ecotone was deployed.
    /// @dev Blocks before this number cannot be used as a reference block for proving against.
    uint immutable ECOTONE_ORIGIN_BLOCK;

    /// @notice The largest allowed gap between the reference block and the block being proved.
    /// @dev This is limited by the ZK proof's ability to cover a large number of blocks.
    uint immutable MAX_BLOCK_GAP;

    /// @notice The verification key used by the SP1Verifier contract.
    bytes32 immutable VKEY;

    /// @notice The timestamp at which the game was created.
    Timestamp public createdAt;

    /// @notice The timestamp at which the game was resolved.
    Timestamp public resolvedAt;

    /// @notice The status of the game (in progress, challenger wins, or defender wins).
    GameStatus public status;

    ////////////////////////////////////////////////////////////////
    //                           SETUP                            //
    ////////////////////////////////////////////////////////////////

    /// @param _factory The DisputeGameFactory contract.
    /// @param _gameType The game type ID.
    /// @param _originBlock The earliest block number that can be proved against.
    /// @param _maxBlockGap The largest allowed gap between the reference block and the block being proved.
    /// @param _vkey The verification key used by the SP1Verifier contract.
    constructor(
        address _factory,
        GameType _gameType,
        uint _originBlock,
        uint _maxBlockGap,
        bytes32 _vkey
    ) {
        // Set all the immutable values in the implementation contract.
        FACTORY = IDisputeGameFactory(_factory);
        GAME_TYPE = _gameType;
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

        // Only allow proving against blocks since Ecotone.
        // This is because the state transition function used in the final ZK proof only computed Ecotone blocks.
        if (prevL2BlockNumber < ECOTONE_ORIGIN_BLOCK) revert InvalidBlockNumber();

        // Require that the block being proven comes after the reference block.
        if (prevL2BlockNumber >= l2BlockNumber()) revert InvalidBlockNumber();

        // Require that the gap between the reference block and the block being proven is less than MAX_BLOCK_GAP.
        // This is used to avoid attacks where such a large gap is used that the game cannot be completed.
        if (prevL2BlockNumber + MAX_BLOCK_GAP < l2BlockNumber()) revert InvalidBlockNumber();

        // Load the Public Values Struct from calldata to validate.
        PublicValuesStruct memory _publicValues = publicValues();

        // Validation #1: The anchor state root is passed as the l2PreRoot.
        if (prevGame.rootClaim().raw() != _publicValues.l2PreRoot) revert InvalidRoot();

        // Validation #2: The claimed root of is passed as the l2PostRoot.
        if (rootClaim() != _publicValues.l2PostRoot) revert InvalidRoot();

        // Validation #3: The real L1 block root matches the passed l1Root.
        // @todo verify relevant L1 block root matches _publicValues.l1Root

        // Validation #4: The real commitment to the blob matches the passed blobKzgCommitment.
        // @todo access correct kzg commitment to verify against _publicValues.blobKzgCommitment?

        // Use the SP1 Verifier to verify the transition function from prevL2BlockNumber to l2BlockNumber.
        verifyProof(VKEY, abi.encode(_publicValues), proofBytes());

        // Set the game's status to resolved in favor of the proposer.
        status = GameStatus.DEFENDER_WINS;

        // Set the game's creation timestamp to the time of initialization.
        createdAt = Timestamp.wrap(uint64(block.timestamp));

        // Set the game's resolved timestamp to the time of initialization.
        resolvedAt = Timestamp.wrap(uint64(block.timestamp));
    }

    ////////////////////////////////////////////////////////////////
    //                           VIEWS                            //
    ////////////////////////////////////////////////////////////////

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

    /// @inheritdoc IBlockDisputeGame
    function l2BlockNumber() public pure returns (uint256 l2BlockNumber_) {
        l2BlockNumber_ = _getArgUint256(0x54);
    }

    /// @return startingRootGameIndex_ The index of the validated game we are proving from.
    function startingRootGameIndex() public pure returns (uint256 startingRootGameIndex_) {
        startingRootGameIndex_ = _getArgUint256(0x74);
    }

    /// @return publicValues_ The public values for the SP1 Verifier proof.
    function publicValues() public pure returns (PublicValuesStruct memory publicValues_) {
        publicValues_.l2PreRoot = _getArgBytes32(0x94);
        publicValues_.l2PostRoot = _getArgBytes32(0xB4);
        publicValues_.l1Root = _getArgBytes32(0xD4);
        publicValues_.blobKzgCommitment = _getArgBytes32(0xF4);
    }

    /// @return proofBytes_ The bytes for the SP1 Verifier proof.
    function proofBytes() public pure returns (bytes memory proofBytes_) {
        // @todo confirm constant size for this proof?
        proofBytes_ = _getArgBytes(0x114, 0x118);
    }
}
