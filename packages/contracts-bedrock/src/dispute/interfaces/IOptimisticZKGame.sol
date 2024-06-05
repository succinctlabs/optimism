// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import { IDisputeGame } from "./IDisputeGame.sol";
import "src/dispute/lib/Types.sol";

interface IOptimisticZKGame is IDisputeGame {
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
    error InvalidRoot();
}
