// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import { IBlockDisputeGame } from "./IBlockDisputeGame.sol";
import "src/dispute/lib/Types.sol";

interface IOptimisticZKGame is IBlockDisputeGame {
    enum IntermediateClaimStatus {
        NONE,
        ACCEPTED,
        CHALLENGED
    }

    struct IntermediateRootClaim {
        OutputRoot outputRoot;
        IntermediateClaimStatus status;
    }

    struct Challenge {
        address challenger;
        Clock proposerClock;
        IntermediateRootClaim left;
        IntermediateRootClaim right;
        OutputRoot current;
        uint totalBonds;
        bool resolved;
        bytes32 l1Head;
    }

    struct PublicValuesStruct {
        bytes32 l1Head;
        bytes32 l2PreRoot;
        bytes32 l2PostRoot;
        uint256 l2BlockNumber;
        uint256 chainId;
    }

    error Unauthorized();
    error ProposerIsChallenger();
    error InvalidDurations();
    error PreviousGameNotResolved();
    error InvalidBlockNumber();
    error WrongTurn();
    error ReadyToProve();
    error NotReadyToProve();
    error InvalidChallengeId();
    error GameAlreadyResolved();
    error ChallengeAlreadyResolved();
    error ChallengeCantBeResolved();
    error WrongBondAmount();
    error TransferFailed();
    error ClockNotExpired();
    error InvalidPublicInput();
}
