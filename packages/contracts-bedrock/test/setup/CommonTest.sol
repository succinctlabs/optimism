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
    bytes proof = hex"fedc1fcc0c401ce68e913792251428e02475772a12db38da45b78f0aaf011a5396fff95b2fbdab6a9fdef15ad8ea58d25e61c8389ff6c6c69a18ff0d30d185fbff7a18e60a6162c634f294b6a46f1ddf269d01437c1d56ec8df35f6d75c490496c031db0040d0149bdb878e9d23753e22105bdcbe6e07571f1617d11c72f5b4f8677d1ba27b05dcb17f68477f2c8ecc72e4c33bb6b85536bed9c466ea68cc99c1e1138e20afd82348fe5180bdbac47e15e121d8f02b2ce22f275b4370e1e0f4a3b65026b0ddf4a1762ec35baa5ad6c92e18c52e06743e1d7b3e105bc5fa4b697612a0bc5258058591da838a15ad3e1c9830b7d5492f86698d68e32d7e8e25e0e050f50351d5b59450932060fa80b6777bba2361f68fc1c3d7661525215c6b5c23a9de2fd2e336640c93e2db350b3549d922873c69d04bcf0090bafe60a8700b56d344d9a0874a88772c9285e4b26eaa3f9c9466e26ae7185879a1ffdb9e30bb080c2ca672c02024943a8ca4cd68b8cb3f2605086a9bcfcf56d080a9277beec07bde0756023ebd7570791c7fb70637f0e9f6109a87c7fe832b2ed91a91dc3c24808f2cc4d196f7b757b6232bfa1bca4b2f5e22487196cf059524be780ee740c2e560e46bf25825045ef51a3ebd67ae727eea5b0e3b10aa1d3a3cc5b4abad6216e803ed20d25f1a77a210ce37f2561a8d909ed597c60285d8f6a405d80320e397a58b54cb420a34b6ebb711c50cb6c4e1d5936e454d265a2b2c3b8aae1632aeffc645352d5154d388031fc6c44a7321ad19053c9d4850e9f85247fc2cfb4c344075e9305a51f33637795a47f7ae0a75be1078d9f42b24b4732a9f1a57df127a5dc6892d18d24137f8ffdb079102030ac3b3da044cf63d23a754e1c4ff71543fd3e49b2c2ec2660c3e111f0354e1b2430da3d1084cb9ae583ac6a1595b95f936be744f06fd819663d6c7b1402a717fff9993806c7de8fa449440e6a321774db6521432babfc2b5d3d4d259c9d0ebb86df54904e8745551e9915bd74107d4fb2e9a1ac37380619fa037d302f248b558520ba08bdcebe22cb912b0e3a94d208e63023cac7c616252b1671b2c81eba47da35e240edb8d5e721d7dbfcc8d88a848eb3c958f6749117d8d678f238ed9881677060908bac0be91d81999d7f2cf40152e3aef3a3a39c04b646b0f8c42b848b41f7fc16ab29a5b38f83c498ce0449f912a1ba7a315719";

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

    // /// @dev Helper function to propose an output.
    // function proposeAnotherOutput() public {
    //     bytes32 proposedOutput2 = keccak256(abi.encode());
    //     uint256 nextBlockNumber = l2OutputOracle.nextBlockNumber();
    //     uint256 nextOutputIndex = l2OutputOracle.nextOutputIndex();
    //     warpToProposeTime(nextBlockNumber);
    //     uint256 proposedNumber = l2OutputOracle.latestBlockNumber();

    //     uint256 submissionInterval = deploy.cfg().l2OutputOracleSubmissionInterval();
    //     // Ensure the submissionInterval is enforced
    //     assertEq(nextBlockNumber, proposedNumber + submissionInterval);

    //     vm.roll(nextBlockNumber + 1);

    //     vm.expectEmit(true, true, true, true);
    //     emit OutputProposed(proposedOutput2, nextOutputIndex, nextBlockNumber, block.timestamp);

    //     address proposer = deploy.cfg().l2OutputOracleProposer();
    //     vm.prank(proposer);
    //     l2OutputOracle.proposeL2Output(proposedOutput2, nextBlockNumber, 0, 0);
    // }

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

    function enableZK(address _verifierGateway) public {
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
    }
}
