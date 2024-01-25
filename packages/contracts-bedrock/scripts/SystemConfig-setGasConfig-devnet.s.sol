// SDPX-License-Identifier: MIT
pragma solidity ^0.8.0;

import { Script, console2 as console } from "forge-std/Script.sol";
import { StdAssertions } from "forge-std/StdAssertions.sol";

import { SystemConfig } from "src/L1/SystemConfig.sol";

contract SystemConfig_SetGasConfig is Script, StdAssertions {
    /// @notice Devnet system config.
    SystemConfig internal constant SYS_CFG = SystemConfig(0x78c9876A3621b97cC9B3C13ad0F35091dB49E8E3);

    uint256 internal constant EXPECTED_SCALAR_CURRENT = 1_000_000;
    uint256 internal constant EXPECTED_OVERHEAD_CURRENT = 2_100;

    uint256 internal constant DEVNET_ECOTONE_GAS_CONFIG = 0x010000000000000000000000000000000000000000000000000d273000001db0;

    // Entrypoint.
    function run() public {
        // Fork the devnet state.
        vm.createSelectFork(vm.envString("ETH_RPC_URL"));

        assertEq(SYS_CFG.scalar(), EXPECTED_SCALAR_CURRENT);
        assertEq(SYS_CFG.overhead(), EXPECTED_OVERHEAD_CURRENT);

        // Set the gas config.
        vm.broadcast();
        SYS_CFG.setGasConfig(0, DEVNET_ECOTONE_GAS_CONFIG);

        console.log("New scalar: %d", SYS_CFG.scalar());
        console.log("New overhead: %d", SYS_CFG.overhead());
    }
}
