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

        bytes memory proof = hex"fedc1fcc0c401ce68e913792251428e02475772a12db38da45b78f0aaf011a5396fff95b2fbdab6a9fdef15ad8ea58d25e61c8389ff6c6c69a18ff0d30d185fbff7a18e60a6162c634f294b6a46f1ddf269d01437c1d56ec8df35f6d75c490496c031db0040d0149bdb878e9d23753e22105bdcbe6e07571f1617d11c72f5b4f8677d1ba27b05dcb17f68477f2c8ecc72e4c33bb6b85536bed9c466ea68cc99c1e1138e20afd82348fe5180bdbac47e15e121d8f02b2ce22f275b4370e1e0f4a3b65026b0ddf4a1762ec35baa5ad6c92e18c52e06743e1d7b3e105bc5fa4b697612a0bc5258058591da838a15ad3e1c9830b7d5492f86698d68e32d7e8e25e0e050f50351d5b59450932060fa80b6777bba2361f68fc1c3d7661525215c6b5c23a9de2fd2e336640c93e2db350b3549d922873c69d04bcf0090bafe60a8700b56d344d9a0874a88772c9285e4b26eaa3f9c9466e26ae7185879a1ffdb9e30bb080c2ca672c02024943a8ca4cd68b8cb3f2605086a9bcfcf56d080a9277beec07bde0756023ebd7570791c7fb70637f0e9f6109a87c7fe832b2ed91a91dc3c24808f2cc4d196f7b757b6232bfa1bca4b2f5e22487196cf059524be780ee740c2e560e46bf25825045ef51a3ebd67ae727eea5b0e3b10aa1d3a3cc5b4abad6216e803ed20d25f1a77a210ce37f2561a8d909ed597c60285d8f6a405d80320e397a58b54cb420a34b6ebb711c50cb6c4e1d5936e454d265a2b2c3b8aae1632aeffc645352d5154d388031fc6c44a7321ad19053c9d4850e9f85247fc2cfb4c344075e9305a51f33637795a47f7ae0a75be1078d9f42b24b4732a9f1a57df127a5dc6892d18d24137f8ffdb079102030ac3b3da044cf63d23a754e1c4ff71543fd3e49b2c2ec2660c3e111f0354e1b2430da3d1084cb9ae583ac6a1595b95f936be744f06fd819663d6c7b1402a717fff9993806c7de8fa449440e6a321774db6521432babfc2b5d3d4d259c9d0ebb86df54904e8745551e9915bd74107d4fb2e9a1ac37380619fa037d302f248b558520ba08bdcebe22cb912b0e3a94d208e63023cac7c616252b1671b2c81eba47da35e240edb8d5e721d7dbfcc8d88a848eb3c958f6749117d8d678f238ed9881677060908bac0be91d81999d7f2cf40152e3aef3a3a39c04b646b0f8c42b848b41f7fc16ab29a5b38f83c498ce0449f912a1ba7a315719";

        super.enableZK(address(verifierGateway), proof);
        super.setUp();
    }
}

contract L2OutputOracle_Test is L2OutputOracle_Init {
    function test_startingL2OutputOracleValues_succeeds() external {
        Types.OutputProposal memory initialProposal = l2OutputOracle.getL2Output(0);
        assertEq(initialProposal.outputRoot, bytes32(deploy.cfg().l2OutputOracleStartingOutputRoot()));
        assertEq(initialProposal.timestamp, deploy.cfg().l2OutputOracleStartingTimestamp());
        assertEq(initialProposal.l2BlockNumber, deploy.cfg().l2OutputOracleStartingBlockNumber());
    }

    function test_proposeWithZKProof_succeeds() external {
        vm.warp(block.timestamp + (l2OutputOracle.SUBMISSION_INTERVAL() * l2OutputOracle.L2_BLOCK_TIME()));

        uint l1BlockNum = deploy.cfg().l2OutputOracleStartingBlockNumber();
        l2OutputOracle.checkpointBlockHash(l1BlockNum, l1Head);

        uint l2BlockNum = l1BlockNum + l2OutputOracle.SUBMISSION_INTERVAL();
        bytes32 lastOutputRoot = l2OutputOracle.getL2Output(l2OutputOracle.latestOutputIndex()).outputRoot;

        L2OutputOracle.PublicValuesStruct memory publicValues = L2OutputOracle.PublicValuesStruct({
            l1Head: l1Head,
            l2PreRoot: lastOutputRoot,
            claimRoot: claimedOutputRoot,
            claimBlockNum: l2BlockNum,
            chainId: 901
        });

        vm.prank(l2OutputOracle.PROPOSER());
        l2OutputOracle.proposeL2Output(claimedOutputRoot, l2BlockNum, l1Head, l1BlockNum, proof);

        assertEq(l2OutputOracle.getL2Output(1).outputRoot, claimedOutputRoot);
    }
}
