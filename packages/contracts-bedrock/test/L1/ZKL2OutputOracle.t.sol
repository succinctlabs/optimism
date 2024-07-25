// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

// Testing utilities
import { stdError, console } from "forge-std/Test.sol";
import { CommonTest } from "test/setup/CommonTest.sol";
import { NextImpl } from "test/mocks/NextImpl.sol";
import { EIP1967Helper } from "test/mocks/EIP1967Helper.sol";

// Libraries
import { Types } from "src/libraries/Types.sol";
import { Constants } from "src/libraries/Constants.sol";

// Target contract dependencies
import { Proxy } from "src/universal/Proxy.sol";

// Target contract
import { ZKL2OutputOracle } from "src/L1/ZKL2OutputOracle.sol";
import { SP1VerifierGateway } from "@sp1-contracts/src/SP1VerifierGateway.sol";
import { SP1Verifier } from "@sp1-contracts/src/v1.0.1/SP1Verifier.sol";

contract ZKL2OutputOracle_Init is CommonTest {
    bytes32 constant ZK_VKEY = 0x003de3cfc15f7b7e2844f33380b8dde65e0cc65de4f7a27e8b3422d376d982f4;
    bytes32 constant STARTING_L2_OUTPUT_ROOT = 0xd1e578c114d50dbb4431e81f737481b7d07204a35d5968c4b911ec55ba038ed6;


    function setUp() public virtual override {
        super.enableZK(ZK_VKEY, STARTING_L2_OUTPUT_ROOT);
        super.setUp();

        address sp1Verifier = address(new SP1Verifier());

        SP1VerifierGateway verifierGateway = new SP1VerifierGateway(address(this));
        verifierGateway.addRoute(sp1Verifier);

        vm.prank(zkL2OutputOracle.PROPOSER());
        zkL2OutputOracle.upgradeVerifier(SP1VerifierGateway(address(verifierGateway)));
    }
}

contract ZKL2OutputOracle_constructor_Test is ZKL2OutputOracle_Init {
    /// @dev Tests that the proxy is initialized with the correct values.
    function test_initialize_succeeds() external {
        address proposer = deploy.cfg().l2OutputOracleProposer();
        address challenger = deploy.cfg().l2OutputOracleChallenger();
        uint256 submissionInterval = deploy.cfg().l2OutputOracleSubmissionInterval();
        uint256 startingBlockNumber = deploy.cfg().l2OutputOracleStartingBlockNumber();
        uint256 startingTimestamp = deploy.cfg().l2OutputOracleStartingTimestamp();
        uint256 l2BlockTime = deploy.cfg().l2BlockTime();
        uint256 finalizationPeriodSeconds = deploy.cfg().finalizationPeriodSeconds();

        assertEq(zkL2OutputOracle.SUBMISSION_INTERVAL(), submissionInterval);
        assertEq(zkL2OutputOracle.submissionInterval(), submissionInterval);
        assertEq(zkL2OutputOracle.L2_BLOCK_TIME(), l2BlockTime);
        assertEq(zkL2OutputOracle.l2BlockTime(), l2BlockTime);
        assertEq(zkL2OutputOracle.latestBlockNumber(), startingBlockNumber);
        assertEq(zkL2OutputOracle.startingBlockNumber(), startingBlockNumber);
        assertEq(zkL2OutputOracle.startingTimestamp(), startingTimestamp);
        assertEq(zkL2OutputOracle.finalizationPeriodSeconds(), finalizationPeriodSeconds);
        assertEq(zkL2OutputOracle.PROPOSER(), proposer);
        assertEq(zkL2OutputOracle.proposer(), proposer);
        assertEq(zkL2OutputOracle.CHALLENGER(), challenger);
        assertEq(zkL2OutputOracle.FINALIZATION_PERIOD_SECONDS(), finalizationPeriodSeconds);
        assertEq(zkL2OutputOracle.challenger(), challenger);
    }
}

contract ZKL2OutputOracle_getter_Test is ZKL2OutputOracle_Init {
    function test_initialL2OutputOracleValues_succeeds() external {
        Types.OutputProposal memory initialProposal = zkL2OutputOracle.getL2Output(0);
        assertEq(initialProposal.outputRoot, bytes32(deploy.cfg().l2OutputOracleStartingOutputRoot()));
        assertEq(initialProposal.timestamp, deploy.cfg().l2OutputOracleStartingTimestamp());
        assertEq(initialProposal.l2BlockNumber, deploy.cfg().l2OutputOracleStartingBlockNumber());
    }

    function test_proposeWithZKProof() external {
        vm.warp(block.timestamp + 1800);

        uint L1_BLOCK_NUM = block.number - 1;
        bytes32 L1_HEAD = 0xb6bd7b941cd0b2098671484897d070508c2d94ad417b484cb30f68edad011578;
        vm.setBlockhash(L1_BLOCK_NUM, L1_HEAD);
        zkL2OutputOracle.checkpointBlockHash(L1_BLOCK_NUM, L1_HEAD);

        uint L2_BLOCK_NUM = 123164174;
        bytes32 CLAIMED_OUTPUT_ROOT = 0x83de1383c4b775a69042d4320cedbef37b0f62b7aec7186bb8fd0a2cebbb8073;
        bytes32 LAST_OUTPUT_ROOT = zkL2OutputOracle.getL2Output(zkL2OutputOracle.latestOutputIndex()).outputRoot;

        ZKL2OutputOracle.PublicValuesStruct memory publicValues = ZKL2OutputOracle.PublicValuesStruct({
            l1Head: L1_HEAD,
            l2PreRoot: LAST_OUTPUT_ROOT,
            claimRoot: CLAIMED_OUTPUT_ROOT,
            claimBlockNum: L2_BLOCK_NUM,
            chainId: 901
        });

        bytes memory proof = hex"fedc1fcc0c401ce68e913792251428e02475772a12db38da45b78f0aaf011a5396fff95b2fbdab6a9fdef15ad8ea58d25e61c8389ff6c6c69a18ff0d30d185fbff7a18e60a6162c634f294b6a46f1ddf269d01437c1d56ec8df35f6d75c490496c031db0040d0149bdb878e9d23753e22105bdcbe6e07571f1617d11c72f5b4f8677d1ba27b05dcb17f68477f2c8ecc72e4c33bb6b85536bed9c466ea68cc99c1e1138e20afd82348fe5180bdbac47e15e121d8f02b2ce22f275b4370e1e0f4a3b65026b0ddf4a1762ec35baa5ad6c92e18c52e06743e1d7b3e105bc5fa4b697612a0bc5258058591da838a15ad3e1c9830b7d5492f86698d68e32d7e8e25e0e050f50351d5b59450932060fa80b6777bba2361f68fc1c3d7661525215c6b5c23a9de2fd2e336640c93e2db350b3549d922873c69d04bcf0090bafe60a8700b56d344d9a0874a88772c9285e4b26eaa3f9c9466e26ae7185879a1ffdb9e30bb080c2ca672c02024943a8ca4cd68b8cb3f2605086a9bcfcf56d080a9277beec07bde0756023ebd7570791c7fb70637f0e9f6109a87c7fe832b2ed91a91dc3c24808f2cc4d196f7b757b6232bfa1bca4b2f5e22487196cf059524be780ee740c2e560e46bf25825045ef51a3ebd67ae727eea5b0e3b10aa1d3a3cc5b4abad6216e803ed20d25f1a77a210ce37f2561a8d909ed597c60285d8f6a405d80320e397a58b54cb420a34b6ebb711c50cb6c4e1d5936e454d265a2b2c3b8aae1632aeffc645352d5154d388031fc6c44a7321ad19053c9d4850e9f85247fc2cfb4c344075e9305a51f33637795a47f7ae0a75be1078d9f42b24b4732a9f1a57df127a5dc6892d18d24137f8ffdb079102030ac3b3da044cf63d23a754e1c4ff71543fd3e49b2c2ec2660c3e111f0354e1b2430da3d1084cb9ae583ac6a1595b95f936be744f06fd819663d6c7b1402a717fff9993806c7de8fa449440e6a321774db6521432babfc2b5d3d4d259c9d0ebb86df54904e8745551e9915bd74107d4fb2e9a1ac37380619fa037d302f248b558520ba08bdcebe22cb912b0e3a94d208e63023cac7c616252b1671b2c81eba47da35e240edb8d5e721d7dbfcc8d88a848eb3c958f6749117d8d678f238ed9881677060908bac0be91d81999d7f2cf40152e3aef3a3a39c04b646b0f8c42b848b41f7fc16ab29a5b38f83c498ce0449f912a1ba7a315719";
        vm.prank(zkL2OutputOracle.PROPOSER());
        zkL2OutputOracle.proposeL2Output(CLAIMED_OUTPUT_ROOT, L2_BLOCK_NUM, L1_HEAD, L1_BLOCK_NUM, proof);

        assertEq(zkL2OutputOracle.getL2Output(1).outputRoot, CLAIMED_OUTPUT_ROOT);
    }
//     bytes32 proposedOutput1 = keccak256(abi.encode(1));

//     /// @dev Tests that `latestBlockNumber` returns the correct value.
//     function test_latestBlockNumber_succeeds() external {
//         uint256 proposedNumber = l2OutputOracle.nextBlockNumber();

//         // Roll to after the block number we'll propose
//         warpToProposeTime(proposedNumber);
//         vm.prank(deploy.cfg().l2OutputOracleProposer());
//         l2OutputOracle.proposeL2Output(proposedOutput1, proposedNumber, 0, 0);
//         assertEq(l2OutputOracle.latestBlockNumber(), proposedNumber);
//     }

//     /// @dev Tests that `getL2Output` returns the correct value.
//     function test_getL2Output_succeeds() external {
//         uint256 nextBlockNumber = l2OutputOracle.nextBlockNumber();
//         uint256 nextOutputIndex = l2OutputOracle.nextOutputIndex();
//         warpToProposeTime(nextBlockNumber);
//         vm.prank(deploy.cfg().l2OutputOracleProposer());
//         l2OutputOracle.proposeL2Output(proposedOutput1, nextBlockNumber, 0, 0);

//         Types.OutputProposal memory proposal = l2OutputOracle.getL2Output(nextOutputIndex);
//         assertEq(proposal.outputRoot, proposedOutput1);
//         assertEq(proposal.timestamp, block.timestamp);

//         // The block number is larger than the latest proposed output:
//         vm.expectRevert(stdError.indexOOBError);
//         l2OutputOracle.getL2Output(nextOutputIndex + 1);
//     }

//     /// @dev Tests that `getL2OutputIndexAfter` returns the correct value
//     ///      when the input is the exact block number of the proposal.
//     function test_getL2OutputIndexAfter_sameBlock_succeeds() external {
//         bytes32 output1 = keccak256(abi.encode(1));
//         uint256 nextBlockNumber1 = l2OutputOracle.nextBlockNumber();
//         warpToProposeTime(nextBlockNumber1);
//         vm.prank(deploy.cfg().l2OutputOracleProposer());
//         l2OutputOracle.proposeL2Output(output1, nextBlockNumber1, 0, 0);

//         // Querying with exact same block as proposed returns the proposal.
//         uint256 index1 = l2OutputOracle.getL2OutputIndexAfter(nextBlockNumber1);
//         assertEq(index1, 0);
//     }

//     /// @dev Tests that `getL2OutputIndexAfter` returns the correct value
//     ///      when the input is the previous block number of the proposal.
//     function test_getL2OutputIndexAfter_previousBlock_succeeds() external {
//         bytes32 output1 = keccak256(abi.encode(1));
//         uint256 nextBlockNumber1 = l2OutputOracle.nextBlockNumber();
//         warpToProposeTime(nextBlockNumber1);
//         vm.prank(deploy.cfg().l2OutputOracleProposer());
//         l2OutputOracle.proposeL2Output(output1, nextBlockNumber1, 0, 0);

//         // Querying with previous block returns the proposal too.
//         uint256 index1 = l2OutputOracle.getL2OutputIndexAfter(nextBlockNumber1 - 1);
//         assertEq(index1, 0);
//     }

//     /// @dev Tests that `getL2OutputIndexAfter` returns the correct value.
//     function test_getL2OutputIndexAfter_multipleOutputsExist_succeeds() external {
//         bytes32 output1 = keccak256(abi.encode(1));
//         uint256 nextBlockNumber1 = l2OutputOracle.nextBlockNumber();
//         warpToProposeTime(nextBlockNumber1);
//         vm.prank(deploy.cfg().l2OutputOracleProposer());
//         l2OutputOracle.proposeL2Output(output1, nextBlockNumber1, 0, 0);

//         bytes32 output2 = keccak256(abi.encode(2));
//         uint256 nextBlockNumber2 = l2OutputOracle.nextBlockNumber();
//         warpToProposeTime(nextBlockNumber2);
//         vm.prank(deploy.cfg().l2OutputOracleProposer());
//         l2OutputOracle.proposeL2Output(output2, nextBlockNumber2, 0, 0);

//         bytes32 output3 = keccak256(abi.encode(3));
//         uint256 nextBlockNumber3 = l2OutputOracle.nextBlockNumber();
//         warpToProposeTime(nextBlockNumber3);
//         vm.prank(deploy.cfg().l2OutputOracleProposer());
//         l2OutputOracle.proposeL2Output(output3, nextBlockNumber3, 0, 0);

//         bytes32 output4 = keccak256(abi.encode(4));
//         uint256 nextBlockNumber4 = l2OutputOracle.nextBlockNumber();
//         warpToProposeTime(nextBlockNumber4);
//         vm.prank(deploy.cfg().l2OutputOracleProposer());
//         l2OutputOracle.proposeL2Output(output4, nextBlockNumber4, 0, 0);

//         // Querying with a block number between the first and second proposal
//         uint256 index1 = l2OutputOracle.getL2OutputIndexAfter(nextBlockNumber1 + 1);
//         assertEq(index1, 1);

//         // Querying with a block number between the second and third proposal
//         uint256 index2 = l2OutputOracle.getL2OutputIndexAfter(nextBlockNumber2 + 1);
//         assertEq(index2, 2);

//         // Querying with a block number between the third and fourth proposal
//         uint256 index3 = l2OutputOracle.getL2OutputIndexAfter(nextBlockNumber3 + 1);
//         assertEq(index3, 3);
//     }

//     /// @dev Tests that `getL2OutputIndexAfter` reverts when no output exists.
//     function test_getL2OutputIndexAfter_noOutputsExis_reverts() external {
//         vm.expectRevert("L2OutputOracle: cannot get output as no outputs have been proposed yet");
//         l2OutputOracle.getL2OutputIndexAfter(0);
//     }

//     /// @dev Tests that `nextBlockNumber` returns the correct value.
//     function test_nextBlockNumber_succeeds() external view {
//         assertEq(
//             l2OutputOracle.nextBlockNumber(),
//             // The return value should match this arithmetic
//             l2OutputOracle.latestBlockNumber() + l2OutputOracle.SUBMISSION_INTERVAL()
//         );
//     }

//     /// @dev Tests that `computeL2Timestamp` returns the correct value.
//     function test_computeL2Timestamp_succeeds() external {
//         uint256 startingBlockNumber = deploy.cfg().l2OutputOracleStartingBlockNumber();
//         uint256 startingTimestamp = deploy.cfg().l2OutputOracleStartingTimestamp();
//         uint256 l2BlockTime = deploy.cfg().l2BlockTime();

//         // reverts if timestamp is too low
//         vm.expectRevert(stdError.arithmeticError);
//         l2OutputOracle.computeL2Timestamp(startingBlockNumber - 1);

//         // check timestamp for the very first block
//         assertEq(l2OutputOracle.computeL2Timestamp(startingBlockNumber), startingTimestamp);

//         // check timestamp for the first block after the starting block
//         assertEq(l2OutputOracle.computeL2Timestamp(startingBlockNumber + 1), startingTimestamp + l2BlockTime);

//         // check timestamp for some other block number
//         assertEq(
//             l2OutputOracle.computeL2Timestamp(startingBlockNumber + 96024), startingTimestamp + l2BlockTime * 96024
//         );
//     }
// }

// contract L2OutputOracle_proposeL2Output_Test is CommonTest {
//     /// @dev Test that `proposeL2Output` succeeds for a valid input
//     ///      and when a block hash and number are not specified.
//     function test_proposeL2Output_proposeAnotherOutput_succeeds() public {
//         proposeAnotherOutput();
//     }

//     /// @dev Tests that `proposeL2Output` succeeds when given valid input and
//     ///      when a block hash and number are specified for reorg protection.
//     function test_proposeWithBlockhashAndHeight_succeeds() external {
//         // Get the number and hash of a previous block in the chain
//         uint256 prevL1BlockNumber = block.number - 1;
//         bytes32 prevL1BlockHash = blockhash(prevL1BlockNumber);

//         uint256 nextBlockNumber = l2OutputOracle.nextBlockNumber();
//         warpToProposeTime(nextBlockNumber);
//         vm.prank(deploy.cfg().l2OutputOracleProposer());
//         l2OutputOracle.proposeL2Output(nonZeroHash, nextBlockNumber, prevL1BlockHash, prevL1BlockNumber);
//     }

//     /// @dev Tests that `proposeL2Output` reverts when called by a party
//     ///      that is not the proposer.
//     function test_proposeL2Output_notProposer_reverts() external {
//         uint256 nextBlockNumber = l2OutputOracle.nextBlockNumber();
//         warpToProposeTime(nextBlockNumber);

//         vm.prank(address(128));
//         vm.expectRevert("L2OutputOracle: only the proposer address can propose new outputs");
//         l2OutputOracle.proposeL2Output(nonZeroHash, nextBlockNumber, 0, 0);
//     }

//     /// @dev Tests that `proposeL2Output` reverts when given a zero blockhash.
//     function test_proposeL2Output_emptyOutput_reverts() external {
//         bytes32 outputToPropose = bytes32(0);
//         uint256 nextBlockNumber = l2OutputOracle.nextBlockNumber();
//         warpToProposeTime(nextBlockNumber);
//         vm.prank(deploy.cfg().l2OutputOracleProposer());
//         vm.expectRevert("L2OutputOracle: L2 output proposal cannot be the zero hash");
//         l2OutputOracle.proposeL2Output(outputToPropose, nextBlockNumber, 0, 0);
//     }

//     /// @dev Tests that `proposeL2Output` reverts when given a block number
//     ///      that does not match the next expected block number.
//     function test_proposeL2Output_unexpectedBlockNumber_reverts() external {
//         uint256 nextBlockNumber = l2OutputOracle.nextBlockNumber();
//         warpToProposeTime(nextBlockNumber);
//         vm.prank(deploy.cfg().l2OutputOracleProposer());
//         vm.expectRevert("L2OutputOracle: block number must be equal to next expected block number");
//         l2OutputOracle.proposeL2Output(nonZeroHash, nextBlockNumber - 1, 0, 0);
//     }

//     /// @dev Tests that `proposeL2Output` reverts when given a block number
//     ///      that has a timestamp in the future.
//     function test_proposeL2Output_futureTimetamp_reverts() external {
//         uint256 nextBlockNumber = l2OutputOracle.nextBlockNumber();
//         uint256 nextTimestamp = l2OutputOracle.computeL2Timestamp(nextBlockNumber);
//         vm.warp(nextTimestamp);
//         vm.prank(deploy.cfg().l2OutputOracleProposer());
//         vm.expectRevert("L2OutputOracle: cannot propose L2 output in the future");
//         l2OutputOracle.proposeL2Output(nonZeroHash, nextBlockNumber, 0, 0);
//     }

//     /// @dev Tests that `proposeL2Output` reverts when given a block number
//     ///      whose hash does not match the given block hash.
//     function test_proposeL2Output_wrongFork_reverts() external {
//         uint256 nextBlockNumber = l2OutputOracle.nextBlockNumber();
//         warpToProposeTime(nextBlockNumber);
//         vm.prank(deploy.cfg().l2OutputOracleProposer());
//         vm.expectRevert("L2OutputOracle: block hash does not match the hash at the expected height");
//         l2OutputOracle.proposeL2Output(nonZeroHash, nextBlockNumber, bytes32(uint256(0x01)), block.number);
//     }

//     /// @dev Tests that `proposeL2Output` reverts when given a block number
//     ///      whose block hash does not match the given block hash.
//     function test_proposeL2Output_unmatchedBlockhash_reverts() external {
//         // Move ahead to block 100 so that we can reference historical blocks
//         vm.roll(100);

//         // Get the number and hash of a previous block in the chain
//         uint256 l1BlockNumber = block.number - 1;
//         bytes32 l1BlockHash = blockhash(l1BlockNumber);

//         uint256 nextBlockNumber = l2OutputOracle.nextBlockNumber();
//         warpToProposeTime(nextBlockNumber);
//         vm.prank(deploy.cfg().l2OutputOracleProposer());

//         // This will fail when foundry no longer returns zerod block hashes
//         vm.expectRevert("L2OutputOracle: block hash does not match the hash at the expected height");
//         l2OutputOracle.proposeL2Output(nonZeroHash, nextBlockNumber, l1BlockHash, l1BlockNumber - 1);
//     }
// }

// contract L2OutputOracle_deleteOutputs_Test is CommonTest {
//     /// @dev Tests that `deleteL2Outputs` succeeds for a single output.
//     function test_deleteOutputs_singleOutput_succeeds() external {
//         proposeAnotherOutput();
//         proposeAnotherOutput();

//         uint256 latestBlockNumber = l2OutputOracle.latestBlockNumber();
//         uint256 latestOutputIndex = l2OutputOracle.latestOutputIndex();
//         Types.OutputProposal memory newLatestOutput = l2OutputOracle.getL2Output(latestOutputIndex - 1);

//         vm.prank(l2OutputOracle.CHALLENGER());
//         vm.prank(l2OutputOracle.challenger());
//         vm.expectEmit(true, true, false, false);
//         emit OutputsDeleted(latestOutputIndex + 1, latestOutputIndex);
//         l2OutputOracle.deleteL2Outputs(latestOutputIndex);

//         // validate latestBlockNumber has been reduced
//         uint256 latestBlockNumberAfter = l2OutputOracle.latestBlockNumber();
//         uint256 latestOutputIndexAfter = l2OutputOracle.latestOutputIndex();
//         uint256 submissionInterval = deploy.cfg().l2OutputOracleSubmissionInterval();
//         assertEq(latestBlockNumber - submissionInterval, latestBlockNumberAfter);

//         // validate that the new latest output is as expected.
//         Types.OutputProposal memory proposal = l2OutputOracle.getL2Output(latestOutputIndexAfter);
//         assertEq(newLatestOutput.outputRoot, proposal.outputRoot);
//         assertEq(newLatestOutput.timestamp, proposal.timestamp);
//     }

//     /// @dev Tests that `deleteL2Outputs` succeeds for multiple outputs.
//     function test_deleteOutputs_multipleOutputs_succeeds() external {
//         proposeAnotherOutput();
//         proposeAnotherOutput();
//         proposeAnotherOutput();
//         proposeAnotherOutput();

//         uint256 latestBlockNumber = l2OutputOracle.latestBlockNumber();
//         uint256 latestOutputIndex = l2OutputOracle.latestOutputIndex();
//         Types.OutputProposal memory newLatestOutput = l2OutputOracle.getL2Output(latestOutputIndex - 3);

//         vm.prank(l2OutputOracle.CHALLENGER());
//         vm.prank(l2OutputOracle.challenger());
//         vm.expectEmit(true, true, false, false);
//         emit OutputsDeleted(latestOutputIndex + 1, latestOutputIndex - 2);
//         l2OutputOracle.deleteL2Outputs(latestOutputIndex - 2);

//         // validate latestBlockNumber has been reduced
//         uint256 latestBlockNumberAfter = l2OutputOracle.latestBlockNumber();
//         uint256 latestOutputIndexAfter = l2OutputOracle.latestOutputIndex();
//         uint256 submissionInterval = deploy.cfg().l2OutputOracleSubmissionInterval();
//         assertEq(latestBlockNumber - submissionInterval * 3, latestBlockNumberAfter);

//         // validate that the new latest output is as expected.
//         Types.OutputProposal memory proposal = l2OutputOracle.getL2Output(latestOutputIndexAfter);
//         assertEq(newLatestOutput.outputRoot, proposal.outputRoot);
//         assertEq(newLatestOutput.timestamp, proposal.timestamp);
//     }

//     /// @dev Tests that `deleteL2Outputs` reverts when not called by the challenger.
//     function test_deleteL2Outputs_ifNotChallenger_reverts() external {
//         uint256 latestBlockNumber = l2OutputOracle.latestBlockNumber();

//         vm.expectRevert("L2OutputOracle: only the challenger address can delete outputs");
//         l2OutputOracle.deleteL2Outputs(latestBlockNumber);
//     }

//     /// @dev Tests that `deleteL2Outputs` reverts for a non-existant output index.
//     function test_deleteL2Outputs_nonExistent_reverts() external {
//         proposeAnotherOutput();

//         uint256 latestBlockNumber = l2OutputOracle.latestBlockNumber();

//         vm.prank(l2OutputOracle.CHALLENGER());
//         vm.prank(l2OutputOracle.challenger());
//         vm.expectRevert("L2OutputOracle: cannot delete outputs after the latest output index");
//         l2OutputOracle.deleteL2Outputs(latestBlockNumber + 1);
//     }

//     /// @dev Tests that `deleteL2Outputs` reverts when trying to delete outputs
//     ///      after the latest output index.
//     function test_deleteL2Outputs_afterLatest_reverts() external {
//         proposeAnotherOutput();
//         proposeAnotherOutput();
//         proposeAnotherOutput();

//         // Delete the latest two outputs
//         uint256 latestOutputIndex = l2OutputOracle.latestOutputIndex();
//         vm.prank(l2OutputOracle.CHALLENGER());
//         vm.prank(l2OutputOracle.challenger());
//         l2OutputOracle.deleteL2Outputs(latestOutputIndex - 2);

//         // Now try to delete the same output again
//         vm.prank(l2OutputOracle.CHALLENGER());
//         vm.prank(l2OutputOracle.challenger());
//         vm.expectRevert("L2OutputOracle: cannot delete outputs after the latest output index");
//         l2OutputOracle.deleteL2Outputs(latestOutputIndex - 2);
//     }

//     /// @dev Tests that `deleteL2Outputs` reverts for finalized outputs.
//     function test_deleteL2Outputs_finalized_reverts() external {
//         proposeAnotherOutput();

//         // Warp past the finalization period + 1 second
//         vm.warp(block.timestamp + l2OutputOracle.FINALIZATION_PERIOD_SECONDS() + 1);

//         uint256 latestOutputIndex = l2OutputOracle.latestOutputIndex();

//         // Try to delete a finalized output
//         vm.prank(l2OutputOracle.CHALLENGER());
//         vm.prank(l2OutputOracle.challenger());
//         vm.expectRevert("L2OutputOracle: cannot delete outputs that have already been finalized");
//         l2OutputOracle.deleteL2Outputs(latestOutputIndex);
//     }
// }

// contract L2OutputOracleUpgradeable_Test is CommonTest {
//     /// @dev Tests that the proxy can be successfully upgraded.
//     function test_upgrading_succeeds() external {
//         Proxy proxy = Proxy(deploy.mustGetAddress("L2OutputOracleProxy"));
//         // Check an unused slot before upgrading.
//         bytes32 slot21Before = vm.load(address(l2OutputOracle), bytes32(uint256(21)));
//         assertEq(bytes32(0), slot21Before);

//         NextImpl nextImpl = new NextImpl();
//         vm.startPrank(EIP1967Helper.getAdmin(address(proxy)));
//         // Reviewer note: the NextImpl() still uses reinitializer. If we want to remove that, we'll need to use a
//         //   two step upgrade with the Storage lib.
//         proxy.upgradeToAndCall(address(nextImpl), abi.encodeWithSelector(NextImpl.initialize.selector, 2));
//         assertEq(proxy.implementation(), address(nextImpl));

//         // Verify that the NextImpl contract initialized its values according as expected
//         bytes32 slot21After = vm.load(address(l2OutputOracle), bytes32(uint256(21)));
//         bytes32 slot21Expected = NextImpl(address(l2OutputOracle)).slot21Init();
//         assertEq(slot21Expected, slot21After);
//     }

//     /// @dev Tests that initialize reverts if the submissionInterval is zero.
//     function test_initialize_submissionInterval_reverts() external {
//         // Reset the initialized field in the 0th storage slot
//         // so that initialize can be called again.
//         vm.store(address(l2OutputOracle), bytes32(uint256(0)), bytes32(uint256(0)));

//         uint256 l2BlockTime = deploy.cfg().l2BlockTime();
//         uint256 startingBlockNumber = deploy.cfg().l2OutputOracleStartingBlockNumber();
//         uint256 startingTimestamp = deploy.cfg().l2OutputOracleStartingTimestamp();
//         address proposer = deploy.cfg().l2OutputOracleProposer();
//         address challenger = deploy.cfg().l2OutputOracleChallenger();
//         uint256 finalizationPeriodSeconds = deploy.cfg().finalizationPeriodSeconds();

//         vm.expectRevert("L2OutputOracle: submission interval must be greater than 0");
//         l2OutputOracle.initialize({
//             _submissionInterval: 0,
//             _l2BlockTime: l2BlockTime,
//             _startingBlockNumber: startingBlockNumber,
//             _startingTimestamp: startingTimestamp,
//             _proposer: proposer,
//             _challenger: challenger,
//             _finalizationPeriodSeconds: finalizationPeriodSeconds
//         });
//     }

//     /// @dev Tests that initialize reverts if the l2BlockTime is invalid.
//     function test_initialize_l2BlockTimeZero_reverts() external {
//         // Reset the initialized field in the 0th storage slot
//         // so that initialize can be called again.
//         vm.store(address(l2OutputOracle), bytes32(uint256(0)), bytes32(uint256(0)));

//         uint256 submissionInterval = deploy.cfg().l2OutputOracleSubmissionInterval();
//         uint256 startingBlockNumber = deploy.cfg().l2OutputOracleStartingBlockNumber();
//         uint256 startingTimestamp = deploy.cfg().l2OutputOracleStartingTimestamp();
//         address proposer = deploy.cfg().l2OutputOracleProposer();
//         address challenger = deploy.cfg().l2OutputOracleChallenger();
//         uint256 finalizationPeriodSeconds = deploy.cfg().finalizationPeriodSeconds();

//         vm.expectRevert("L2OutputOracle: L2 block time must be greater than 0");
//         l2OutputOracle.initialize({
//             _submissionInterval: submissionInterval,
//             _l2BlockTime: 0,
//             _startingBlockNumber: startingBlockNumber,
//             _startingTimestamp: startingTimestamp,
//             _proposer: proposer,
//             _challenger: challenger,
//             _finalizationPeriodSeconds: finalizationPeriodSeconds
//         });
//     }

//     /// @dev Tests that initialize reverts if the starting timestamp is invalid.
//     function test_initialize_badTimestamp_reverts() external {
//         // Reset the initialized field in the 0th storage slot
//         // so that initialize can be called again.
//         vm.store(address(l2OutputOracle), bytes32(uint256(0)), bytes32(uint256(0)));

//         uint256 submissionInterval = deploy.cfg().l2OutputOracleSubmissionInterval();
//         uint256 l2BlockTime = deploy.cfg().l2BlockTime();
//         uint256 startingBlockNumber = deploy.cfg().l2OutputOracleStartingBlockNumber();
//         address proposer = deploy.cfg().l2OutputOracleProposer();
//         address challenger = deploy.cfg().l2OutputOracleChallenger();
//         uint256 finalizationPeriodSeconds = deploy.cfg().finalizationPeriodSeconds();

//         vm.expectRevert("L2OutputOracle: starting L2 timestamp must be less than current time");
//         l2OutputOracle.initialize({
//             _submissionInterval: submissionInterval,
//             _l2BlockTime: l2BlockTime,
//             _startingBlockNumber: startingBlockNumber,
//             _startingTimestamp: block.timestamp + 1,
//             _proposer: proposer,
//             _challenger: challenger,
//             _finalizationPeriodSeconds: finalizationPeriodSeconds
//         });
//     }
}
