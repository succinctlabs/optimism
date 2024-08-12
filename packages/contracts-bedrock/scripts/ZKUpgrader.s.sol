// // SPDX-License-Identifier: MIT
// pragma solidity ^0.8.0;

// import { VmSafe } from "forge-std/Vm.sol";
// import { Script } from "forge-std/Script.sol";

// import { console2 as console } from "forge-std/console2.sol";
// import { stdJson } from "forge-std/StdJson.sol";

// import { GnosisSafe as Safe } from "safe-contracts/GnosisSafe.sol";
// import { OwnerManager } from "safe-contracts/base/OwnerManager.sol";
// import { GnosisSafeProxyFactory as SafeProxyFactory } from "safe-contracts/proxies/GnosisSafeProxyFactory.sol";
// import { Enum as SafeOps } from "safe-contracts/common/Enum.sol";

// import { Deployer } from "scripts/deploy/Deployer.sol";

// import { ProxyAdmin } from "src/universal/ProxyAdmin.sol";
// import { AddressManager } from "src/legacy/AddressManager.sol";
// import { Proxy } from "src/universal/Proxy.sol";
// import { L1StandardBridge } from "src/L1/L1StandardBridge.sol";
// import { StandardBridge } from "src/universal/StandardBridge.sol";
// import { OptimismPortal } from "src/L1/OptimismPortal.sol";
// import { OptimismPortal2 } from "src/L1/OptimismPortal2.sol";
// import { OptimismPortalInterop } from "src/L1/OptimismPortalInterop.sol";
// import { L1ChugSplashProxy } from "src/legacy/L1ChugSplashProxy.sol";
// import { ResolvedDelegateProxy } from "src/legacy/ResolvedDelegateProxy.sol";
// import { L1CrossDomainMessenger } from "src/L1/L1CrossDomainMessenger.sol";
// import { L2OutputOracle } from "src/L1/L2OutputOracle.sol";
// import { OptimismMintableERC20Factory } from "src/universal/OptimismMintableERC20Factory.sol";
// import { SuperchainConfig } from "src/L1/SuperchainConfig.sol";
// import { SystemConfig } from "src/L1/SystemConfig.sol";
// import { SystemConfigInterop } from "src/L1/SystemConfigInterop.sol";
// import { ResourceMetering } from "src/L1/ResourceMetering.sol";
// import { DataAvailabilityChallenge } from "src/L1/DataAvailabilityChallenge.sol";
// import { Constants } from "src/libraries/Constants.sol";
// import { DisputeGameFactory } from "src/dispute/DisputeGameFactory.sol";
// import { FaultDisputeGame } from "src/dispute/FaultDisputeGame.sol";
// import { PermissionedDisputeGame } from "src/dispute/PermissionedDisputeGame.sol";
// import { DelayedWETH } from "src/dispute/weth/DelayedWETH.sol";
// import { AnchorStateRegistry } from "src/dispute/AnchorStateRegistry.sol";
// import { PreimageOracle } from "src/cannon/PreimageOracle.sol";
// import { MIPS } from "src/cannon/MIPS.sol";
// import { L1ERC721Bridge } from "src/L1/L1ERC721Bridge.sol";
// import { ProtocolVersions, ProtocolVersion } from "src/L1/ProtocolVersions.sol";
// import { StorageSetter } from "src/universal/StorageSetter.sol";
// import { Predeploys } from "src/libraries/Predeploys.sol";
// import { Chains } from "scripts/Chains.sol";
// import { Config } from "scripts/Config.sol";

// import { IBigStepper } from "src/dispute/interfaces/IBigStepper.sol";
// import { IPreimageOracle } from "src/cannon/interfaces/IPreimageOracle.sol";
// import { AlphabetVM } from "test/mocks/AlphabetVM.sol";
// import "src/dispute/lib/Types.sol";
// import { ChainAssertions } from "scripts/ChainAssertions.sol";
// import { Types } from "scripts/Types.sol";
// import { LibStateDiff } from "scripts/libraries/LibStateDiff.sol";
// import { EIP1967Helper } from "test/mocks/EIP1967Helper.sol";
// import { ForgeArtifacts } from "scripts/ForgeArtifacts.sol";
// import { Process } from "scripts/libraries/Process.sol";

// contract ZKUpgrader is Deployer {
//     ////////////////////////////////////////////////////////////////
//     //                        Modifiers                           //
//     ////////////////////////////////////////////////////////////////

//     /// @notice Modifier that wraps a function in broadcasting.
//     modifier broadcast() {
//         vm.startBroadcast(msg.sender);
//         _;
//         vm.stopBroadcast();
//     }

//     ////////////////////////////////////////////////////////////////
//     //                        Functions                           //
//     ////////////////////////////////////////////////////////////////

//     function upgradeToZK() public broadcast {
//         // deploy an implementation of the ZK L2OO
//         L2OutputOracle zkL2OutputOracleImpl = new L2OutputOracle();

//         // get the address of the L2OutputOracleProxy
//         address l2OutputOracleAddress = mustGetAddress("L2OutputOracleProxy");

//         // require that we have permission to upgrade the contract
//         require(Proxy(l2OutputOracleAddress).admin() == msg.sender, "ZKUpgrader: not admin");

//         if (L2OutputOracle(l2OutputOracleAddress).startingBlockNumber() == 0) {
//             // it's a fresh deployment and hasn't been initialized
//             Proxy(l2OutputOracleAddress).upgradeToAndCall(
//                 address(zkL2OutputOracleImpl),
//                 abi.encodeCall(L2OutputOracle.initialize, (address(0), address(0)))
//             );
//         } else {
//             // upgrading existing contract, everything is already set
//             Proxy(l2OutputOracleAddress).upgradeTo(address(zkL2ooImpl));
//         }
//     }
// }
