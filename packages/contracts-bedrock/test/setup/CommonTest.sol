// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Test } from "forge-std/Test.sol";
import { Setup } from "test/setup/Setup.sol";
import { Events } from "test/setup/Events.sol";
import { FFIInterface } from "test/setup/FFIInterface.sol";
import { Constants } from "src/libraries/Constants.sol";
import "scripts/DeployConfig.s.sol";

/// @title CommonTest
/// @dev An extenstion to `Test` that sets up the optimism smart contracts.
contract CommonTest is Test, Setup, Events {
    address alice;
    address bob;

    bytes32 constant nonZeroHash = keccak256(abi.encode("NON_ZERO"));

    FFIInterface constant ffi = FFIInterface(address(uint160(uint256(keccak256(abi.encode("optimism.ffi"))))));

    bool usePlasmaOverride;
    bool useFaultProofs;
    address customGasToken;
    bool useInteropOverride;
    bool useZK;
    address zkVerifierGateway;
    bytes proof;
    bytes32 l1Head;
    bytes32 claimedOutputRoot;

    function setUp() public virtual override {
        alice = makeAddr("alice");
        bob = makeAddr("bob");
        vm.deal(alice, 10000 ether);
        vm.deal(bob, 10000 ether);

        Setup.setUp();

        // Override the config after the deploy script initialized the config
        if (usePlasmaOverride) {
            deploy.cfg().setUsePlasma(true);
        }
        if (useFaultProofs) {
            deploy.cfg().setUseFaultProofs(true);
        }
        if (customGasToken != address(0)) {
            deploy.cfg().setUseCustomGasToken(customGasToken);
        }
        if (useInteropOverride) {
            deploy.cfg().setUseInterop(true);
        }
        if (useZK) {
            deploy.cfg().setVerifierGateway(zkVerifierGateway);
        }

        vm.etch(address(ffi), vm.getDeployedCode("FFIInterface.sol:FFIInterface"));
        vm.label(address(ffi), "FFIInterface");

        // Exclude contracts for the invariant tests
        excludeContract(address(ffi));
        excludeContract(address(deploy));
        excludeContract(address(deploy.cfg()));

        // Make sure the base fee is non zero
        vm.fee(1 gwei);

        // Set sane initialize block numbers
        vm.warp(deploy.cfg().l2OutputOracleStartingTimestamp() + 1);
        vm.roll(deploy.cfg().l2OutputOracleStartingBlockNumber() + 1);

        // Deploy L1
        Setup.L1();
        // Deploy L2
        Setup.L2();

        // Set the block hash and cache in for ZK L2OutputOracle tests
        vm.setBlockhash(deploy.cfg().l2OutputOracleStartingBlockNumber(), l1Head);
        l2OutputOracle.checkpointBlockHash(deploy.cfg().l2OutputOracleStartingBlockNumber(), l1Head);
    }

    /// @dev Helper function that wraps `TransactionDeposited` event.
    ///      The magic `0` is the version.
    function emitTransactionDeposited(
        address _from,
        address _to,
        uint256 _mint,
        uint256 _value,
        uint64 _gasLimit,
        bool _isCreation,
        bytes memory _data
    )
        internal
    {
        emit TransactionDeposited(_from, _to, 0, abi.encodePacked(_mint, _value, _gasLimit, _isCreation, _data));
    }

    // @dev Advance the evm's time to meet the L2OutputOracle's requirements for proposeL2Output
    function warpToProposeTime(uint256 _nextBlockNumber) public {
        vm.warp(l2OutputOracle.computeL2Timestamp(_nextBlockNumber) + 1);
    }

    /// @dev Helper function to propose an output.
    function proposeAnotherOutput() public {
        bytes32 proposedOutput2 = claimedOutputRoot;
        uint256 nextBlockNumber = deploy.cfg().l2OutputOracleStartingBlockNumber() + l2OutputOracle.SUBMISSION_INTERVAL();
        uint256 nextOutputIndex = l2OutputOracle.nextOutputIndex();
        warpToProposeTime(nextBlockNumber);
        uint256 proposedNumber = l2OutputOracle.latestBlockNumber();

        uint256 submissionInterval = deploy.cfg().l2OutputOracleSubmissionInterval();
        // Ensure the submissionInterval is enforced
        assertEq(nextBlockNumber, proposedNumber + submissionInterval);

        vm.roll(nextBlockNumber + 1);

        vm.expectEmit(true, true, true, true);
        emit OutputProposed(proposedOutput2, nextOutputIndex, nextBlockNumber, block.timestamp);

        address proposer = deploy.cfg().l2OutputOracleProposer();
        vm.prank(proposer);
        l2OutputOracle.proposeL2Output(proposedOutput2, nextBlockNumber, l1Head, deploy.cfg().l2OutputOracleStartingBlockNumber(), proof);
    }

    function enableFaultProofs() public {
        // Check if the system has already been deployed, based off of the heuristic that alice and bob have not been
        // set by the `setUp` function yet.
        if (!(alice == address(0) && bob == address(0))) {
            revert("CommonTest: Cannot enable fault proofs after deployment. Consider overriding `setUp`.");
        }

        useFaultProofs = true;
    }

    function enablePlasma() public {
        // Check if the system has already been deployed, based off of the heuristic that alice and bob have not been
        // set by the `setUp` function yet.
        if (!(alice == address(0) && bob == address(0))) {
            revert("CommonTest: Cannot enable plasma after deployment. Consider overriding `setUp`.");
        }

        usePlasmaOverride = true;
    }

    function enableCustomGasToken(address _token) public {
        // Check if the system has already been deployed, based off of the heuristic that alice and bob have not been
        // set by the `setUp` function yet.
        if (!(alice == address(0) && bob == address(0))) {
            revert("CommonTest: Cannot enable custom gas token after deployment. Consider overriding `setUp`.");
        }
        require(_token != Constants.ETHER);

        customGasToken = _token;
    }

    function enableInterop() public {
        // Check if the system has already been deployed, based off of the heuristic that alice and bob have not been
        // set by the `setUp` function yet.
        if (!(alice == address(0) && bob == address(0))) {
            revert("CommonTest: Cannot enable interop after deployment. Consider overriding `setUp`.");
        }

        useInteropOverride = true;
    }

    function enableZK(address _verifierGateway, bytes memory _proof) public {
        // Check if the system has already been deployed, based off of the heuristic that alice and bob have not been
        // set by the `setUp` function yet.
        if (!(alice == address(0) && bob == address(0))) {
            revert("CommonTest: Cannot enable ZK after deployment. Consider overriding `setUp`.");
        }
        // Check if the system already has Fault Proofs on, in which case ZK cannot be added.
        if (useFaultProofs) {
            revert("ZK and Fault Proofs are mutually exclusive");
        }

        useZK = true;
        zkVerifierGateway = _verifierGateway;
        proof = _proof;
        l1Head = 0xb6bd7b941cd0b2098671484897d070508c2d94ad417b484cb30f68edad011578;
        claimedOutputRoot = 0x83de1383c4b775a69042d4320cedbef37b0f62b7aec7186bb8fd0a2cebbb8073;
    }
}
