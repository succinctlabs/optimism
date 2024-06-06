// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import { IDisputeGame } from "./IDisputeGame.sol";
import "src/dispute/lib/Types.sol";

/// @title IBlockDisputeGame
/// @notice The interface for a game meant to resolve an L2 Block.
interface IBlockDisputeGame is IDisputeGame {
    /// @notice The l2BlockNumber of the disputed output root in the `L2OutputOracle`.
    /// @return l2BlockNumber_ The block number that the game claiming is proving.
    function l2BlockNumber() external view returns (uint256 l2BlockNumber_);
}
