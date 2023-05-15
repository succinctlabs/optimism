// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { console } from "forge-std/console.sol";
import { Script } from "forge-std/Script.sol";

/**
 * @title BedrockMigrationChecker
 * @notice A script to check safety of multisig operations for Bedrock.
 *         The usage is as follows:
 *         $ forge script scripts/CheckForBedrockMigration.s.sol \
 *             --rpc-url $ETH_RPC_URL
 */

contract BedrockMigrationChecker is Script {
    /**
     * @notice The entrypoint function.
     */
    function run() external {
        console.log("hello world");
    }
}
