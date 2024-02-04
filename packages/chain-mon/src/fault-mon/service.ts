import {
  BaseServiceV2,
  StandardOptions,
  ExpressRouter,
  Gauge,
  validators,
  waitForProvider,
} from '@eth-optimism/common-ts'
import {
  BedrockOutputData,
  getChainId,
  sleep,
  toRpcHexString,
} from '@eth-optimism/core-utils'
import { getOEContract, DEFAULT_L2_CONTRACT_ADDRESSES } from '@eth-optimism/sdk'
import { Provider } from '@ethersproject/abstract-provider'
import { config } from 'dotenv'
import { Contract, ethers } from 'ethers'
import dateformat from 'dateformat'

import { version } from '../../package.json'
import { findFirstUnfinalizedOutputIndex, findOutputForIndex } from './helpers'

type Options = {
  l1RpcProvider: Provider
  l2RpcProvider: Provider
  startOutputIndex: number
  optimismPortalAddress: string
  l2OutputOracleAddress: string
}

type Metrics = {
  highestOutputIndex: Gauge
  isCurrentlyMismatched: Gauge
  nodeConnectionFailures: Gauge
}

type State = {
  faultProofWindow: number
  optimismPortal: Contract
  optimismPortal2: Contract
  l2OutputOracle: Contract
  currentOutputIndex: number
  diverged: boolean
}

export class FaultDetector extends BaseServiceV2<Options, Metrics, State> {
  constructor(options?: Partial<Options & StandardOptions>) {
    super({
      version,
      name: 'fault-detector',
      loop: true,
      options: {
        loopIntervalMs: 1000,
        ...options,
      },
      optionsSpec: {
        l1RpcProvider: {
          validator: validators.provider,
          desc: 'Provider for interacting with L1',
        },
        l2RpcProvider: {
          validator: validators.provider,
          desc: 'Provider for interacting with L2',
        },
        startOutputIndex: {
          validator: validators.num,
          default: -1,
          desc: 'The L2 height to start from',
          public: true,
        },
        optimismPortalAddress: {
          validator: validators.str,
          default: null,
          desc: 'Address of the OptimismPortal proxy contract on L1',
          public: true,
        },
        l2OutputOracleAddress: {
          validator: validators.str,
          default: null,
          desc: 'Address of the L2OutputOracle proxy contract on L1',
          public: true,
        },
      },
      metricsSpec: {
        highestOutputIndex: {
          type: Gauge,
          desc: 'Highest output indices (checked and known)',
          labels: ['type'],
        },
        isCurrentlyMismatched: {
          type: Gauge,
          desc: '0 if state is ok, 1 if state is mismatched',
        },
        nodeConnectionFailures: {
          type: Gauge,
          desc: 'Number of times node connection has failed',
          labels: ['layer', 'section'],
        },
      },
    })
  }

  async init(): Promise<void> {
    // Connect to L1.
    await waitForProvider(this.options.l1RpcProvider, {
      logger: this.logger,
      name: 'L1',
    })

    // Connect to L2.
    await waitForProvider(this.options.l2RpcProvider, {
      logger: this.logger,
      name: 'L2',
    })

    // Need L2 chain ID to resolve contract addresses.
    const l2ChainId = await getChainId(this.options.l2RpcProvider)

    this.state.optimismPortal = getOEContract('OptimismPortal', l2ChainId, {
      signerOrProvider: this.options.l1RpcProvider,
      address: this.options.optimismPortalAddress,
    })

    this.state.optimismPortal2 = getOEContract('OptimismPortal2', l2ChainId, {
      signerOrProvider: this.options.l1RpcProvider,
      address: this.state.optimismPortal.address,
    })

    this.state.l2OutputOracle = getOEContract('L2OutputOracle', l2ChainId, {
      signerOrProvider: this.options.l1RpcProvider,
      address: this.options.l2OutputOracleAddress,
    })

    // We use this a lot, a bit cleaner to pull out to the top level of the state object.
    this.state.faultProofWindow =
      this.state.l2OutputOracle.FINALIZATION_PERIOD_SECONDS()

    // Figure out where to start syncing from.
    if (this.options.startOutputIndex === -1) {
      this.logger.info('finding appropriate starting unfinalized output')
      const firstUnfinalized = await findFirstUnfinalizedOutputIndex(
        this.state.l2OutputOracle,
        this.state.faultProofWindow,
        this.logger
      )

      // We may not have an unfinalized outputs in the case where no outputs have been submitted
      // for the entire duration of the FAULT_PROOF_WINDOW. We generally do not expect this to
      // happen on mainnet, but it happens on testnets because the FAULT_PROOF_WINDOW is short.
      if (firstUnfinalized === undefined) {
        this.logger.info('no unfinalized outputs found, skipping all outputs')
        const totalOutputs = await this.state.l2OutputOracle.nextOutputIndex()
        this.state.currentOutputIndex = totalOutputs.toNumber() - 1
      } else {
        this.state.currentOutputIndex = firstUnfinalized
      }
    } else {
      this.state.currentOutputIndex = this.options.startOutputIndex
    }

    // Not diverged by default.
    this.state.diverged = false
    this.metrics.isCurrentlyMismatched.set(0)

    // Log the initial state.
    this.logger.info('initial state', {
      faultProofWindow: this.state.faultProofWindow,
      outputIndex: this.state.currentOutputIndex,
    })
  }

  async routes(router: ExpressRouter): Promise<void> {
    router.get('/status', async (req, res) => {
      return res.status(200).json({
        ok: !this.state.diverged,
      })
    })
  }

  async main(): Promise<void> {
    const startMs = Date.now()

    let latestOutputIndex: number
    try {
      const totalOutputs = await this.state.l2OutputOracle.nextOutputIndex()
      latestOutputIndex = totalOutputs.toNumber() - 1
    } catch (err) {
      this.logger.error('failed to query total # of outputs', {
        error: err,
        node: 'l1',
        section: 'nextOutputIndex',
      })
      this.metrics.nodeConnectionFailures.inc({
        layer: 'l1',
        section: 'nextOutputIndex',
      })
      await sleep(15000)
      return
    }

    if (this.state.currentOutputIndex > latestOutputIndex) {
      this.logger.info('output index is ahead of the oracle. waiting...', {
        outputIndex: this.state.currentOutputIndex,
        latestOutputIndex,
      })
      await sleep(15000)
      return
    }

    this.metrics.highestOutputIndex.set({ type: 'known' }, latestOutputIndex)
    this.logger.info('checking output', {
      outputIndex: this.state.currentOutputIndex,
      latestOutputIndex,
    })

    let outputData: BedrockOutputData
    try {
      outputData = await findOutputForIndex(
        this.state.l2OutputOracle,
        this.state.currentOutputIndex,
        this.logger
      )
    } catch (err) {
      this.logger.error('failed to fetch output associated with output', {
        error: err,
        node: 'l1',
        section: 'findOutputForIndex',
        outputIndex: this.state.currentOutputIndex,
      })
      this.metrics.nodeConnectionFailures.inc({
        layer: 'l1',
        section: 'findOutputForIndex',
      })
      await sleep(15000)
      return
    }

    let latestBlock: number
    try {
      latestBlock = await this.options.l2RpcProvider.getBlockNumber()
    } catch (err) {
      this.logger.error('failed to query L2 block height', {
        error: err,
        node: 'l2',
        section: 'getBlockNumber',
      })
      this.metrics.nodeConnectionFailures.inc({
        layer: 'l2',
        section: 'getBlockNumber',
      })
      await sleep(15000)
      return
    }

    const outputBlockNumber = outputData.l2BlockNumber
    if (latestBlock < outputBlockNumber) {
      this.logger.info('L2 node is behind, waiting for sync...', {
        l2BlockHeight: latestBlock,
        outputBlock: outputBlockNumber,
      })
      return
    }

    let outputBlock: any
    try {
      outputBlock = await (
        this.options.l2RpcProvider as ethers.providers.JsonRpcProvider
      ).send('eth_getBlockByNumber', [toRpcHexString(outputBlockNumber), false])
    } catch (err) {
      this.logger.error('failed to fetch output block', {
        error: err,
        node: 'l2',
        section: 'getBlock',
        block: outputBlockNumber,
      })
      this.metrics.nodeConnectionFailures.inc({
        layer: 'l2',
        section: 'getBlock',
      })
      await sleep(15000)
      return
    }

    let messagePasserProofResponse: any
    try {
      messagePasserProofResponse = await (
        this.options.l2RpcProvider as ethers.providers.JsonRpcProvider
      ).send('eth_getProof', [
        DEFAULT_L2_CONTRACT_ADDRESSES.BedrockMessagePasser,
        [],
        toRpcHexString(outputBlockNumber),
      ])
    } catch (err) {
      this.logger.error('failed to fetch message passer proof', {
        error: err,
        node: 'l2',
        section: 'getProof',
        block: outputBlockNumber,
      })
      this.metrics.nodeConnectionFailures.inc({
        layer: 'l2',
        section: 'getProof',
      })
      await sleep(15000)
      return
    }

    const outputRoot = ethers.utils.solidityKeccak256(
      ['uint256', 'bytes32', 'bytes32', 'bytes32'],
      [
        0,
        outputBlock.stateRoot,
        messagePasserProofResponse.storageHash,
        outputBlock.hash,
      ]
    )

    if (outputRoot !== outputData.outputRoot) {
      this.state.diverged = true
      this.metrics.isCurrentlyMismatched.set(1)
      this.logger.error('state root mismatch', {
        blockNumber: outputBlock.number,
        expectedStateRoot: outputData.outputRoot,
        actualStateRoot: outputRoot,
        finalizationTime: dateformat(
          new Date(
            (ethers.BigNumber.from(outputBlock.timestamp).toNumber() +
              this.state.faultProofWindow) *
              1000
          ),
          'mmmm dS, yyyy, h:MM:ss TT'
        ),
      })
      return
    }

    const elapsedMs = Date.now() - startMs

    // Mark the current output index as checked
    this.logger.info('checked output ok', {
      outputIndex: this.state.currentOutputIndex,
      timeMs: elapsedMs,
    })
    this.metrics.highestOutputIndex.set(
      { type: 'checked' },
      this.state.currentOutputIndex
    )

    // If we got through the above without throwing an error, we should be
    // fine to reset and move onto the next output
    this.state.diverged = false
    this.state.currentOutputIndex++
    this.metrics.isCurrentlyMismatched.set(0)
  }
}

if (require.main === module) {
  config()
  const service = new FaultDetector()
  service.run()
}
