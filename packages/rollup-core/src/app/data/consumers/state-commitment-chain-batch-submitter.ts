/* External Imports */
import { getLogger, logError, ScheduledTask } from '@eth-optimism/core-utils'
import { Contract } from 'ethers'

/* Internal Imports */
import {
  BatchSubmissionStatus,
  L2DataService,
  StateCommitmentBatchSubmission,
} from '../../../types/data'
import { TransactionReceipt, TransactionResponse } from 'ethers/providers'

const log = getLogger('state-commitment-chain-batch-submitter')

/**
 * Polls the DB for L2 batches ready to send to L1 and submits them.
 */
export class StateCommitmentChainBatchSubmitter extends ScheduledTask {
  constructor(
    private readonly dataService: L2DataService,
    private readonly stateCommitmentChain: Contract,
    periodMilliseconds = 10_000
  ) {
    super(periodMilliseconds)
  }

  /**
   * @inheritDoc
   *
   * Submits L2 batches from L2 Transactions in the DB whenever there is a batch that is ready.
   */
  public async runTask(): Promise<boolean> {
    let stateBatch: StateCommitmentBatchSubmission
    try {
      stateBatch = await this.dataService.getNextStateCommitmentBatchToSubmit()
    } catch (e) {
      logError(
        log,
        `Error fetching state root batch for L1 submission! Continuing...`,
        e
      )
      return false
    }

    if (!stateBatch) {
      log.debug(`No state batches ready for submission.`)
      return false
    }

    if (stateBatch.status !== BatchSubmissionStatus.QUEUED) {
      const msg = `Received state commitment batch to finalize in ${
        stateBatch.status
      } instead of ${
        BatchSubmissionStatus.SENT
      }. Batch Submission: ${JSON.stringify(stateBatch)}.`
      log.error(msg)
      throw msg
    }

    const txHash: string = await this.buildAndSendRollupBatchTransaction(
      stateBatch
    )
    if (!txHash) {
      return false
    }
    return this.waitForProofThatTransactionSucceeded(txHash, stateBatch)
  }

  /**
   * Builds and sends a Rollup State Root Batch transaction to L1, returning its tx hash.
   *
   * @param stateRootBatch The state root batch to send to L1.
   * @returns The L1 tx hash.
   */
  private async buildAndSendRollupBatchTransaction(
    stateRootBatch: StateCommitmentBatchSubmission
  ): Promise<string> {
    let txHash: string
    try {
      const stateRoots: string[] = stateRootBatch.stateRoots

      log.debug(
        `Appending state root batch number: ${stateRootBatch.batchNumber} with ${stateRoots.length} state roots.`
      )

      const txRes: TransactionResponse = await this.stateCommitmentChain.appendStateBatch(
        stateRoots
      )
      log.debug(
        `State Root batch ${stateRootBatch.batchNumber} appended with at least one confirmation! Tx Hash: ${txRes.hash}`
      )
      txHash = txRes.hash
    } catch (e) {
      logError(
        log,
        `Error submitting State Root batch ${stateRootBatch.batchNumber} to state commitment chain!`,
        e
      )
      return undefined
    }

    return txHash
  }

  /**
   * Waits for a confirm to indicate that the transaction did not fail.
   *
   * @param txHash The tx hash to wait for.
   * @param stateRootBatch The rollup batch in question.
   * @returns true if the tx was successful and false otherwise.
   */
  private async waitForProofThatTransactionSucceeded(
    txHash: string,
    stateRootBatch: StateCommitmentBatchSubmission
  ): Promise<boolean> {
    try {
      const receipt: TransactionReceipt = await this.stateCommitmentChain.provider.waitForTransaction(
        txHash,
        1
      )
      if (!receipt.status) {
        log.error(
          `Error submitting State Root batch # ${stateRootBatch.batchNumber} to L1!. Batch: ${stateRootBatch}`
        )
        return false
      }
    } catch (e) {
      logError(
        log,
        `Error submitting State Root batch # ${stateRootBatch.batchNumber} to L1!. Batch: ${stateRootBatch}`,
        e
      )
      return false
    }

    try {
      log.debug(
        `Marking State Root batch ${stateRootBatch.batchNumber} submitted`
      )
      await this.dataService.markStateRootBatchSubmittedToL1(
        stateRootBatch.batchNumber,
        txHash
      )
    } catch (e) {
      logError(
        log,
        `Error marking State Root batch ${stateRootBatch.batchNumber} as submitted!`,
        e
      )
      // TODO: Should we return here? Don't want to resubmit, so I think we should update the DB
      return false
    }
    return true
  }
}
