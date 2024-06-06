// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import { IBlockDisputeGame } from "./IBlockDisputeGame.sol";
import "src/dispute/lib/Types.sol";

interface IPureZKGame is IBlockDisputeGame {
    struct PublicValuesStruct {
        bytes32 l1Root;
        bytes32 l2PreRoot;
        bytes32 l2PostRoot;
        bytes32 blobKzgCommitment;
    }

    // error Unauthorized();
    // error ProposerIsChallenger();
    // error InvalidDurations();
    // error PreviousGameNotResolved();
    // error InvalidBlockNumber();
    // error WrongTurn();
    // error ReadyToProve();
    // error NotReadyToProve();
    // error InvalidChallengeId();
    // error GameAlreadyResolved();
    // error ChallengeAlreadyResolved();
    // error ChallengeCantBeResolved();
    // error WrongBondAmount();
    // error TransferFailed();
    // error ClockNotExpired();
    // error InvalidRoot();
    // error InvalidBlobCommitment();
}
