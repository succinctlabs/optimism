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
        super.setUp();

        address sp1Verifier = address(new SP1Verifier());
        SP1VerifierGateway verifierGateway = new SP1VerifierGateway();
        verifierGateway.addRoute(sp1Verifier);

        vm.startPrank(l2OutputOracle.PROPOSER());
        l2OutputOracle.setVerifierGateway(address(verifierGateway));

        vm.setBlockhash(deploy.cfg().l2OutputOracleStartingBlockNumber(), l1Head);
        l2OutputOracle.checkpointBlockHash(deploy.cfg().l2OutputOracleStartingBlockNumber(), l1Head);

        l2OutputOracle.setVKey(zkVKey);
        l2OutputOracle.setInitialOutputRoot(startingOutputRoot);
        vm.stopPrank();
    }
}

contract L2OutputOracle_Test is L2OutputOracle_Init {
    function test_startingL2OutputOracleValues_succeeds() external {
        Types.OutputProposal memory initialProposal = l2OutputOracle.getL2Output(0);
        assertEq(initialProposal.outputRoot, startingOutputRoot);
        assertEq(initialProposal.timestamp, deploy.cfg().l2OutputOracleStartingTimestamp());
        assertEq(initialProposal.l2BlockNumber, deploy.cfg().l2OutputOracleStartingBlockNumber());
    }

    function test_proposeWithZKProof_succeeds() external {
        vm.warp(block.timestamp + (l2OutputOracle.SUBMISSION_INTERVAL() * l2OutputOracle.L2_BLOCK_TIME()));

        uint l1BlockNum = deploy.cfg().l2OutputOracleStartingBlockNumber();
        uint l2BlockNum = l1BlockNum + l2OutputOracle.SUBMISSION_INTERVAL();
        bytes32 lastOutputRoot = l2OutputOracle.getL2Output(l2OutputOracle.latestOutputIndex()).outputRoot;

        L2OutputOracle.PublicValuesStruct memory publicValues = L2OutputOracle.PublicValuesStruct({
            l1Head: l1Head,
            l2PreRoot: lastOutputRoot,
            claimRoot: claimedOutputRoot,
            claimBlockNum: l2BlockNum,
            chainId: deploy.cfg().l2ChainID()
        });

        vm.prank(l2OutputOracle.PROPOSER());
        l2OutputOracle.proposeL2Output(claimedOutputRoot, l2BlockNum, l1Head, l1BlockNum, proof);

        assertEq(l2OutputOracle.getL2Output(1).outputRoot, claimedOutputRoot);
    }
}
