// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import { Test, console } from "forge-std/Test.sol";
import { OptimisticZKGame } from "src/dispute/OptimisticZKGame.sol";
import { DisputeGameFactory } from "src/dispute/DisputeGameFactory.sol";

contract OptimisticZKGameTest is Test {

    // function setUp() public {
    //     super.enableFaultProofs();
    //     super.setUp();
    //     DisputeGameFactory factory = new DisputeGameFactory();
    //     factory.initialize(address(this));

    //     OptimisticZKGame gameImpl = new OptimisticZKGame({
    //         _factory: address(factory),
    //         _gameType: GameType.OPTIMISTIC_ZK,
    //         _maxGameDuration: 6 days,
    //         _maxProposerDuration: 3 days,
    //         _weth:
    //         _originBlock: 0,
    //         _maxBlockGap: 52 weeks / 2
    //     });


    //     factory.setImplementation(GameType.OPTIMISTIC_ZK, address(gameImpl));
    // }
}
