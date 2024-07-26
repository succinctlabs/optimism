// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

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
import { L2OutputOracle } from "src/L1/L2OutputOracle.sol";
import { SP1VerifierGateway } from "@sp1-contracts/src/SP1VerifierGateway.sol";
import { SP1Verifier } from "@sp1-contracts/src/v1.0.1/SP1Verifier.sol";

contract L2OutputOracle_Init is CommonTest {
    function setUp() public virtual override {
        address sp1Verifier = address(new SP1Verifier());
        SP1VerifierGateway verifierGateway = new SP1VerifierGateway();
        verifierGateway.addRoute(sp1Verifier);

        super.enableZK(address(verifierGateway));
        super.setUp();
    }
}

contract L2OutputOracle_Test is L2OutputOracle_Init {
    function test_startingL2OutputOracleValues_succeeds() external {
        Types.OutputProposal memory initialProposal = L2OutputOracle.getL2Output(0);
        assertEq(initialProposal.outputRoot, bytes32(deploy.cfg().l2OutputOracleStartingOutputRoot()));
        assertEq(initialProposal.timestamp, deploy.cfg().l2OutputOracleStartingTimestamp());
        assertEq(initialProposal.l2BlockNumber, deploy.cfg().l2OutputOracleStartingBlockNumber());
    }

    function test_proposeWithZKProof_succeeds() external {
        vm.warp(block.timestamp + (l2OutputOracle.SUBMISSION_INTERVAL() * l2OutputOracle.L2_BLOCK_TIME()));

        uint L1_BLOCK_NUM = block.number - 1;
        bytes32 L1_HEAD = 0xb6bd7b941cd0b2098671484897d070508c2d94ad417b484cb30f68edad011578;
        vm.setBlockhash(L1_BLOCK_NUM, L1_HEAD);
        l2OutputOracle.checkpointBlockHash(L1_BLOCK_NUM, L1_HEAD);

        uint L2_BLOCK_NUM = L1_BLOCK_NUM + l2OutputOracle.SUBMISSION_INTERVAL();
        bytes32 CLAIMED_OUTPUT_ROOT = 0x83de1383c4b775a69042d4320cedbef37b0f62b7aec7186bb8fd0a2cebbb8073;
        bytes32 LAST_OUTPUT_ROOT = l2OutputOracle.getL2Output(l2OutputOracle.latestOutputIndex()).outputRoot;

        L2OutputOracle.PublicValuesStruct memory publicValues = L2OutputOracle.PublicValuesStruct({
            l1Head: L1_HEAD,
            l2PreRoot: LAST_OUTPUT_ROOT,
            claimRoot: CLAIMED_OUTPUT_ROOT,
            claimBlockNum: L2_BLOCK_NUM,
            chainId: 901
        });

        vm.prank(L2OutputOracle.PROPOSER());
        l2OutputOracle.proposeL2Output(CLAIMED_OUTPUT_ROOT, L2_BLOCK_NUM, L1_HEAD, L1_BLOCK_NUM, proof);

        assertEq(l2OutputOracle.getL2Output(1).outputRoot, CLAIMED_OUTPUT_ROOT);
    }
}
