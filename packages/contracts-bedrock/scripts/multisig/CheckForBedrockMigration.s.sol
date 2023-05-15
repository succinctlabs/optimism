// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { console2 } from "forge-std/console2.sol";
import { Script } from "forge-std/Script.sol";
import { StdAssertions } from "forge-std/StdAssertions.sol";

/**
 * @title BedrockMigrationChecker
 * @notice A script to check safety of multisig operations for Bedrock.
 *         The usage is as follows:
 *         $ forge script scripts/CheckForBedrockMigration.s.sol \
 *             --rpc-url $ETH_RPC_URL
 */

contract BedrockMigrationChecker is Script, StdAssertions {

    struct ContractSet {
        // Please keep these sorted by name.
        address AddressManager;
        address L1CrossDomainMessengerImpl;
        address L1CrossDomainMessengerProxy;
        address L1ERC721BridgeImpl;
        address L1ERC721BridgeProxy;
        address L1ProxyAdmin;
        address L1StandardBridgeImpl;
        address L1StandardBridgeProxy;
        address L1UpgradeKey;
        address L2OutputOracleImpl;
        address L2OutputOracleProxy;
        address OptimismMintableERC20FactoryImpl;
        address OptimismMintableERC20FactoryProxy;
        address OptimismPortalImpl;
        address OptimismPortalProxy;
        address PortalSender;
        address SystemConfigProxy;
        address SystemDictatorImpl;
        address SystemDictatorProxy;
    }

    /**
     * @notice The entrypoint function.
     */
    function run() external {
        string memory deploymentJsonDir = vm.envString("DEPLOYMENT_JSON_DIR");
        console2.log("DEPLOYMENT_JSON_DIR = %s", deploymentJsonDir);
        ContractSet memory contracts = getContracts(deploymentJsonDir);
        checkAddressManager(contracts);
        checkL1CrossDomainMessengerImpl(contracts);
        checkL1CrossDomainMessengerProxy(contracts);
        checkL1ERC721BridgeImpl(contracts);
        checkL1ERC721BridgeProxy(contracts);
        checkL1ProxyAdmin(contracts);
        checkL1StandardBridgeImpl(contracts);
        checkL1StandardBridgeProxy(contracts);
        checkL1UpgradeKey(contracts);
        checkL2OutputOracleImpl(contracts);
        checkL2OutputOracleProxy(contracts);
        checkOptimismMintableERC20FactoryImpl(contracts);
        checkOptimismMintableERC20FactoryProxy(contracts);
        checkOptimismPortalImpl(contracts);
        checkOptimismPortalProxy(contracts);
        checkPortalSender(contracts);
        checkSystemConfigProxy(contracts);
        checkSystemDictatorImpl(contracts);
        checkSystemDictatorProxy(contracts);
    }

    function checkAddressManager(ContractSet memory contracts) internal {
        console2.log("Checking AddressManager %s", contracts.AddressManager);
        checkAddressIsExpected(contracts.L1UpgradeKey, contracts.AddressManager, "owner()");
    }

    function checkL1CrossDomainMessengerImpl(ContractSet memory contracts) internal {
        console2.log("Checking L1CrossDomainMessenger %s", contracts.L1CrossDomainMessengerImpl);
    }

    function checkL1CrossDomainMessengerProxy(ContractSet memory contracts) internal {
        console2.log("Checking L1CrossDomainMessengerProxy %s", contracts.L1CrossDomainMessengerProxy);
        checkAddressIsExpected(contracts.L1UpgradeKey, contracts.L1CrossDomainMessengerProxy, "owner()");
        checkAddressIsExpected(contracts.AddressManager, contracts.L1CrossDomainMessengerProxy, "libAddressManager()");
    }

    function checkL1ERC721BridgeImpl(ContractSet memory contracts) internal {
        console2.log("Checking L1ERC721Bridge %s", contracts.L1ERC721BridgeImpl);
    }

    function checkL1ERC721BridgeProxy(ContractSet memory contracts) internal {
        console2.log("Checking L1ERC721BridgeProxy %s", contracts.L1ERC721BridgeProxy);
        checkAddressIsExpected(contracts.L1UpgradeKey, contracts.L1ERC721BridgeProxy, "admin()");
        checkAddressIsExpected(contracts.L1CrossDomainMessengerProxy, contracts.L1ERC721BridgeProxy, "messenger()");
    }

    function checkL1ProxyAdmin(ContractSet memory contracts) internal {
        console2.log("Checking L1ProxyAdmin %s", contracts.L1ProxyAdmin);
        checkAddressIsExpected(contracts.L1UpgradeKey, contracts.L1ProxyAdmin, "owner()");
    }

    function checkL1StandardBridgeImpl(ContractSet memory contracts) internal {
        console2.log("Checking L1StandardBridge %s", contracts.L1StandardBridgeImpl);
    }

    function checkL1StandardBridgeProxy(ContractSet memory contracts) internal {
        console2.log("Checking L1StandardBridgeProxy %s", contracts.L1StandardBridgeProxy);
        checkAddressIsExpected(contracts.L1UpgradeKey, contracts.L1StandardBridgeProxy, "getOwner()");
        checkAddressIsExpected(contracts.L1CrossDomainMessengerProxy, contracts.L1StandardBridgeProxy, "messenger()");
    }

    function checkL1UpgradeKey(ContractSet memory contracts) internal {
        console2.log("Checking L1UpgradeKeyAddress %s", contracts.L1UpgradeKey);
    }

    function checkL2OutputOracleImpl(ContractSet memory contracts) internal {
        console2.log("Checking L2OutputOracle %s", contracts.L2OutputOracleImpl);
    }

    function checkL2OutputOracleProxy(ContractSet memory contracts) internal {
        console2.log("Checking L2OutputOracleProxy %s", contracts.L2OutputOracleProxy);
        checkAddressIsExpected(contracts.L1ProxyAdmin, contracts.L2OutputOracleProxy, "admin()");
    }

    function checkOptimismMintableERC20FactoryImpl(ContractSet memory contracts) internal {
        console2.log("Checking OptimismMintableERC20Factory %s", contracts.OptimismMintableERC20FactoryImpl);
    }

    function checkOptimismMintableERC20FactoryProxy(ContractSet memory contracts) internal {
        console2.log("Checking OptimismMintableERC20FactoryProxy %s", contracts.OptimismMintableERC20FactoryProxy);
        checkAddressIsExpected(contracts.L1ProxyAdmin, contracts.OptimismMintableERC20FactoryProxy, "admin()");
        checkAddressIsExpected(contracts.L1StandardBridgeProxy, contracts.OptimismMintableERC20FactoryProxy, "BRIDGE()");
    }

    function checkOptimismPortalImpl(ContractSet memory contracts) internal {
        console2.log("Checking OptimismPortal %s", contracts.OptimismPortalImpl);
    }

    function checkOptimismPortalProxy(ContractSet memory contracts) internal {
        console2.log("Checking OptimismPortalProxy %s", contracts.OptimismPortalProxy);
        checkAddressIsExpected(contracts.L1ProxyAdmin, contracts.OptimismPortalProxy, "admin()");
        checkAddressIsExpected(contracts.L2OutputOracleProxy, contracts.OptimismPortalProxy, "L2_ORACLE()");
    }

    function checkPortalSender(ContractSet memory contracts) internal {
        console2.log("Checking PortalSender %s", contracts.PortalSender);
        checkAddressIsExpected(contracts.OptimismPortalProxy, contracts.PortalSender, "PORTAL()");
    }

    function checkSystemConfigProxy(ContractSet memory contracts) internal {
        console2.log("Checking SystemConfigProxy %s", contracts.SystemConfigProxy);
    }

    function checkSystemDictatorImpl(ContractSet memory contracts) internal {
        console2.log("Checking SystemDictator %s", contracts.SystemDictatorImpl);
    }

    function checkSystemDictatorProxy(ContractSet memory contracts) internal {
        console2.log("Checking SystemDictatorProxy %s", contracts.SystemDictatorProxy);
        checkAddressIsExpected(contracts.SystemDictatorImpl, contracts.SystemDictatorProxy, "implementation()");
        checkAddressIsExpected(contracts.L1UpgradeKey, contracts.SystemDictatorProxy, "owner()");
        checkAddressIsExpected(contracts.L1UpgradeKey, contracts.SystemDictatorProxy, "admin()");
    }

    function checkAddressIsExpected(address expectedAddr, address contractAddr, string memory signature) internal {
        address actual = getAddressFromCall(contractAddr, signature);
        if (expectedAddr != actual) {
            console2.log("  !! Error: %s != %s.%s, ", expectedAddr, contractAddr, signature);
            console2.log("           which is %s", actual);
        } else {
            console2.log("  -- Success: %s == %s.%s.", expectedAddr, contractAddr, signature);
        }
    }

    function getAddressFromCall(address contractAddr, string memory signature) internal returns (address) {
        vm.prank(address(0));
        (bool success, bytes memory addrBytes) = contractAddr.staticcall(abi.encodeWithSignature(signature));
        if (!success) {
            console2.log("  !! Error calling %s.%s", contractAddr, signature);
            return address(0);
        }
        return abi.decode(addrBytes, (address));
    }

    function getContracts(string memory deploymentJsonDir) internal returns (ContractSet memory) {
        return ContractSet({
            AddressManager: getAddressFromJson(string.concat(deploymentJsonDir, "/Lib_AddressManager.json")),
                    L1CrossDomainMessengerImpl: getAddressFromJson(string.concat(deploymentJsonDir, "/L1CrossDomainMessenger.json")),
                    L1CrossDomainMessengerProxy: getAddressFromJson(string.concat(deploymentJsonDir, "/Proxy__OVM_L1CrossDomainMessenger.json")),
                L1ERC721BridgeImpl: getAddressFromJson(string.concat(deploymentJsonDir, "/L1ERC721Bridge.json")),
                L1ERC721BridgeProxy: getAddressFromJson(string.concat(deploymentJsonDir, "/L1ERC721BridgeProxy.json")),
                L1ProxyAdmin: getAddressFromJson(string.concat(deploymentJsonDir, "/ProxyAdmin.json")),
                L1StandardBridgeImpl: getAddressFromJson(string.concat(deploymentJsonDir, "/L1StandardBridge.json")),
                L1StandardBridgeProxy: getAddressFromJson(string.concat(deploymentJsonDir, "/Proxy__OVM_L1StandardBridge.json")),
                L1UpgradeKey: vm.envAddress("L1_UPGRADE_KEY"),
                L2OutputOracleImpl: getAddressFromJson(string.concat(deploymentJsonDir, "/L2OutputOracle.json")),
                L2OutputOracleProxy: getAddressFromJson(string.concat(deploymentJsonDir, "/L2OutputOracleProxy.json")),
                OptimismMintableERC20FactoryImpl: getAddressFromJson(string.concat(deploymentJsonDir, "/OptimismMintableERC20Factory.json")),
                OptimismMintableERC20FactoryProxy: getAddressFromJson(string.concat(deploymentJsonDir, "/OptimismMintableERC20FactoryProxy.json")),
                OptimismPortalImpl: getAddressFromJson(string.concat(deploymentJsonDir, "/OptimismPortal.json")),
                OptimismPortalProxy: getAddressFromJson(string.concat(deploymentJsonDir, "/OptimismPortalProxy.json")),
                PortalSender: getAddressFromJson(string.concat(deploymentJsonDir, "/PortalSender.json")),
                SystemConfigProxy: getAddressFromJson(string.concat(deploymentJsonDir, "/SystemConfigProxy.json")),
                SystemDictatorImpl: getAddressFromJson(string.concat(deploymentJsonDir, "/SystemDictator.json")),
                SystemDictatorProxy: getAddressFromJson(string.concat(deploymentJsonDir, "/SystemDictatorProxy.json"))
            });
    }

    function getAddressFromJson(string memory jsonPath) internal returns (address) {
        string memory json = vm.readFile(jsonPath);
        return vm.parseJsonAddress(json, ".address");
    }

}
