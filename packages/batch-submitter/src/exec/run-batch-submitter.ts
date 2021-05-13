/* External Imports */
import { injectL2Context, Bcfg } from '@eth-optimism/core-utils'
import { Logger, Metrics, createMetricsServer } from '@eth-optimism/common-ts'
import { exit } from 'process'
import { Signer, Wallet } from 'ethers'
import { JsonRpcProvider, TransactionReceipt } from '@ethersproject/providers'
import * as dotenv from 'dotenv'
import Config from 'bcfg'

/* Internal Imports */
import {
  TransactionBatchSubmitter,
  AutoFixBatchOptions,
  StateBatchSubmitter,
  STATE_BATCH_SUBMITTER_LOG_TAG,
  TX_BATCH_SUBMITTER_LOG_TAG,
} from '..'

interface RequiredEnvVars {
  // The HTTP provider URL for L1.
  L1_NODE_WEB3_URL: string
  // The HTTP provider URL for L2.
  L2_NODE_WEB3_URL: string
  // The layer one address manager address
  ADDRESS_MANAGER_ADDRESS: string
  // The minimum size in bytes of any L1 transactions generated by the batch submitter.
  MIN_L1_TX_SIZE: number
  // The maximum size in bytes of any L1 transactions generated by the batch submitter.
  MAX_L1_TX_SIZE: number
  // The maximum number of L2 transactions that can ever be in a batch.
  MAX_TX_BATCH_COUNT: number
  // The maximum number of L2 state roots that can ever be in a batch.
  MAX_STATE_BATCH_COUNT: number
  // The maximum amount of time (seconds) that we will wait before submitting an under-sized batch.
  MAX_BATCH_SUBMISSION_TIME: number
  // The delay in milliseconds between querying L2 for more transactions / to create a new batch.
  POLL_INTERVAL: number
  // The number of confirmations which we will wait after appending new batches.
  NUM_CONFIRMATIONS: number
  // The number of seconds to wait before resubmitting a transaction.
  RESUBMISSION_TIMEOUT: number
  // The number of confirmations that we should wait before submitting state roots for CTC elements.
  FINALITY_CONFIRMATIONS: number
  // Whether or not to run the tx batch submitter.
  RUN_TX_BATCH_SUBMITTER: boolean
  // Whether or not to run the state batch submitter.
  RUN_STATE_BATCH_SUBMITTER: boolean
  // The safe minimum amount of ether the batch submitter key should
  // hold before it starts to log errors.
  SAFE_MINIMUM_ETHER_BALANCE: number
  // A boolean to clear the pending transactions in the mempool
  // on start up.
  CLEAR_PENDING_TXS: boolean
}

/* Optional Env Vars
 * FRAUD_SUBMISSION_ADDRESS
 * DISABLE_QUEUE_BATCH_APPEND
 * SEQUENCER_PRIVATE_KEY
 * PROPOSER_PRIVATE_KEY
 * MNEMONIC
 * SEQUENCER_MNEMONIC
 * PROPOSER_MNEMONIC
 * SEQUENCER_HD_PATH
 * PROPOSER_HD_PATH
 * BLOCK_OFFSET
 * USE_HARDHAT
 * DEBUG_IMPERSONATE_SEQUENCER_ADDRESS
 * DEBUG_IMPERSONATE_PROPOSER_ADDRESS
 * RUN_PROMETHEUS_SERVER
 */

export const run = async () => {
  dotenv.config()

  const config: Bcfg = new Config('batch-submitter')
  config.load({
    env: true,
    argv: true,
  })

  // Parse config
  const env = process.env
  const environment = config.str('node-env', env.NODE_ENV)
  const network = config.str('eth-network-name', env.ETH_NETWORK_NAME)
  const release = `batch-submitter@${env.npm_package_version}`
  const sentryDsn = config.str('sentry-dsn', env.SENTRY_DSN)
  const sentryTraceRate = config.ufloat(
    'sentry-trace-rate',
    parseFloat(env.SENTRY_TRACE_RATE) || 0.05
  )

  // Default is 1 because Geth normally has 1 more block than L1
  const BLOCK_OFFSET = config.uint(
    'block-offset',
    parseInt(env.BLOCK_OFFSET, 10) || 1
  )

  /* Logger */
  const name = 'oe:batch_submitter:init'
  let logger

  if (config.bool('use-sentry', env.USE_SENTRY === 'true')) {
    // Initialize Sentry for Batch Submitter deployed to a network
    logger = new Logger({
      name,
      sentryOptions: {
        release,
        dsn: sentryDsn,
        tracesSampleRate: sentryTraceRate,
        environment: network,
      },
    })
  } else {
    // Skip initializing Sentry
    logger = new Logger({ name })
  }

  const useHardhat = config.bool('use-hardhat', !!env.USE_HARDAT)
  const DEBUG_IMPERSONATE_SEQUENCER_ADDRESS = config.str(
    'debug-impersonate-sequencer-address',
    env.DEBUG_IMPERSONATE_SEQUENCER_ADDRESS
  )
  const DEBUG_IMPERSONATE_PROPOSER_ADDRESS = config.str(
    'debug-impersonate-proposer-address',
    env.DEBUG_IMPERSONATE_PROPOSER_ADDRESS
  )

  const getSequencerSigner = async (): Promise<Signer> => {
    const l1Provider = new JsonRpcProvider(requiredEnvVars.L1_NODE_WEB3_URL)

    if (useHardhat) {
      if (!DEBUG_IMPERSONATE_SEQUENCER_ADDRESS) {
        throw new Error('Must pass DEBUG_IMPERSONATE_SEQUENCER_ADDRESS')
      }
      await l1Provider.send('hardhat_impersonateAccount', [
        DEBUG_IMPERSONATE_SEQUENCER_ADDRESS,
      ])
      return l1Provider.getSigner(DEBUG_IMPERSONATE_SEQUENCER_ADDRESS)
    }

    if (SEQUENCER_PRIVATE_KEY) {
      return new Wallet(SEQUENCER_PRIVATE_KEY, l1Provider)
    } else if (SEQUENCER_MNEMONIC) {
      return Wallet.fromMnemonic(SEQUENCER_MNEMONIC, SEQUENCER_HD_PATH).connect(
        l1Provider
      )
    }
    throw new Error(
      'Must pass one of SEQUENCER_PRIVATE_KEY, MNEMONIC, or SEQUENCER_MNEMONIC'
    )
  }

  const getProposerSigner = async (): Promise<Signer> => {
    const l1Provider = new JsonRpcProvider(requiredEnvVars.L1_NODE_WEB3_URL)

    if (useHardhat) {
      if (!DEBUG_IMPERSONATE_PROPOSER_ADDRESS) {
        throw new Error('Must pass DEBUG_IMPERSONATE_PROPOSER_ADDRESS')
      }
      await l1Provider.send('hardhat_impersonateAccount', [
        DEBUG_IMPERSONATE_PROPOSER_ADDRESS,
      ])
      return l1Provider.getSigner(DEBUG_IMPERSONATE_PROPOSER_ADDRESS)
    }

    if (PROPOSER_PRIVATE_KEY) {
      return new Wallet(PROPOSER_PRIVATE_KEY, l1Provider)
    } else if (PROPOSER_MNEMONIC) {
      return Wallet.fromMnemonic(PROPOSER_MNEMONIC, PROPOSER_HD_PATH).connect(
        l1Provider
      )
    }
    throw new Error(
      'Must pass one of PROPOSER_PRIVATE_KEY, MNEMONIC, or PROPOSER_MNEMONIC'
    )
  }

  /* Metrics */
  const metrics = new Metrics({
    prefix: name,
    labels: { environment, release, network },
  })

  const FRAUD_SUBMISSION_ADDRESS = config.str(
    'fraud-submisison-address',
    env.FRAUD_SUBMISSION_ADDRESS || 'no fraud'
  )
  const DISABLE_QUEUE_BATCH_APPEND = config.bool(
    'disable-queue-batch-append',
    !!env.DISABLE_QUEUE_BATCH_APPEND
  )
  const MIN_GAS_PRICE_IN_GWEI = config.uint(
    'min-gas-price-in-gwei',
    parseInt(env.MIN_GAS_PRICE_IN_GWEI, 10) || 0
  )
  const MAX_GAS_PRICE_IN_GWEI = config.uint(
    'max-gas-price-in-gwei',
    parseInt(env.MAX_GAS_PRICE_IN_GWEI, 10) || 70
  )
  const GAS_RETRY_INCREMENT = config.uint(
    'gas-retry-increment',
    parseInt(env.GAS_RETRY_INCREMENT, 10) || 5
  )
  const GAS_THRESHOLD_IN_GWEI = config.uint(
    'gas-threshold-in-gwei',
    parseInt(env.GAS_THRESHOLD_IN_GWEI, 10) || 100
  )

  // Private keys & mnemonics
  const SEQUENCER_PRIVATE_KEY = config.str(
    'sequencer-private-key',
    env.SEQUENCER_PRIVATE_KEY
  )
  // Kept for backwards compatibility
  const PROPOSER_PRIVATE_KEY = config.str(
    'proposer-private-key',
    env.PROPOSER_PRIVATE_KEY || env.SEQUENCER_PRIVATE_KEY
  )
  const SEQUENCER_MNEMONIC = config.str(
    'sequencer-mnemonic',
    env.SEQUENCER_MNEMONIC || env.MNEMONIC
  )
  const PROPOSER_MNEMONIC = config.str(
    'proposer-mnemonic',
    env.PROPOSER_MNEMONIC || env.MNEMONIC
  )
  const SEQUENCER_HD_PATH = config.str(
    'sequencer-hd-path',
    env.SEQUENCER_HD_PATH || env.HD_PATH
  )
  const PROPOSER_HD_PATH = config.str(
    'proposer-hd-path',
    env.PROPOSER_HD_PATH || env.HD_PATH
  )

  // Auto fix batch options -- TODO: Remove this very hacky config
  const AUTO_FIX_BATCH_OPTIONS_CONF = config.str(
    'auto-fix-batch-conf',
    env.AUTO_FIX_BATCH_OPTIONS_CONF || ''
  )
  const autoFixBatchOptions: AutoFixBatchOptions = {
    fixDoublePlayedDeposits: AUTO_FIX_BATCH_OPTIONS_CONF
      ? AUTO_FIX_BATCH_OPTIONS_CONF.includes('fixDoublePlayedDeposits')
      : false,
    fixMonotonicity: AUTO_FIX_BATCH_OPTIONS_CONF
      ? AUTO_FIX_BATCH_OPTIONS_CONF.includes('fixMonotonicity')
      : false,
    fixSkippedDeposits: AUTO_FIX_BATCH_OPTIONS_CONF
      ? AUTO_FIX_BATCH_OPTIONS_CONF.includes('fixSkippedDeposits')
      : false,
  }

  logger.info('Starting batch submitter...')

  const requiredEnvVars: RequiredEnvVars = {
    L1_NODE_WEB3_URL: config.str('l1-node-web3-url', env.L1_NODE_WEB3_URL),
    L2_NODE_WEB3_URL: config.str('l2-node-web3-url', env.L2_NODE_WEB3_URL),
    ADDRESS_MANAGER_ADDRESS: config.str(
      'address-manager-address',
      env.ADDRESS_MANAGER_ADDRESS
    ),
    MIN_L1_TX_SIZE: config.uint(
      'min-l1-tx-size',
      parseInt(env.MIN_L1_TX_SIZE, 10)
    ),
    MAX_L1_TX_SIZE: config.uint(
      'max-l1-tx-size',
      parseInt(env.MAX_L1_TX_SIZE, 10)
    ),
    MAX_TX_BATCH_COUNT: config.uint(
      'max-tx-batch-count',
      parseInt(env.MAX_TX_BATCH_COUNT, 10)
    ),
    MAX_STATE_BATCH_COUNT: config.uint(
      'max-state-batch-count',
      parseInt(env.MAX_STATE_BATCH_COUNT, 10)
    ),
    MAX_BATCH_SUBMISSION_TIME: config.uint(
      'max-batch-submisison-time',
      parseInt(env.MAX_BATCH_SUBMISSION_TIME, 10)
    ),
    POLL_INTERVAL: config.uint(
      'poll-interval',
      parseInt(env.POLL_INTERVAL, 10)
    ),
    NUM_CONFIRMATIONS: config.uint(
      'num-confirmations',
      parseInt(env.NUM_CONFIRMATIONS, 10)
    ),
    RESUBMISSION_TIMEOUT: config.uint(
      'resubmission-timeout',
      parseInt(env.RESUBMISSION_TIMEOUT, 10)
    ),
    FINALITY_CONFIRMATIONS: config.uint(
      'finality-confirmations',
      parseInt(env.FINALITY_CONFIRMATIONS, 10)
    ),
    RUN_TX_BATCH_SUBMITTER: config.bool(
      'run-tx-batch-submitter',
      env.RUN_TX_BATCH_SUBMITTER === 'true'
    ),
    RUN_STATE_BATCH_SUBMITTER: config.bool(
      'run-state-batch-submitter',
      env.RUN_STATE_BATCH_SUBMITTER === 'true'
    ),
    SAFE_MINIMUM_ETHER_BALANCE: config.ufloat(
      'safe-minimum-ether-balance',
      parseFloat(env.SAFE_MINIMUM_ETHER_BALANCE)
    ),
    CLEAR_PENDING_TXS: config.bool(
      'clear-pending-txs',
      env.CLEAR_PENDING_TXS === 'true'
    ),
  }

  for (const [key, val] of Object.entries(requiredEnvVars)) {
    if (val === null || val === undefined) {
      logger.warn('Missing environment variable', {
        key,
        value: val,
      })
      exit(1)
    }
  }

  const clearPendingTxs = requiredEnvVars.CLEAR_PENDING_TXS

  const l2Provider = injectL2Context(
    new JsonRpcProvider(requiredEnvVars.L2_NODE_WEB3_URL)
  )

  const sequencerSigner: Signer = await getSequencerSigner()
  let proposerSigner: Signer = await getProposerSigner()

  const sequencerAddress = await sequencerSigner.getAddress()
  const proposerAddress = await proposerSigner.getAddress()
  // If the sequencer & proposer are the same, use a single wallet
  if (sequencerAddress === proposerAddress) {
    proposerSigner = sequencerSigner
  }

  logger.info('Configured batch submitter addresses', {
    sequencerAddress,
    proposerAddress,
    addressManagerAddress: requiredEnvVars.ADDRESS_MANAGER_ADDRESS,
  })

  const txBatchSubmitter = new TransactionBatchSubmitter(
    sequencerSigner,
    l2Provider,
    requiredEnvVars.MIN_L1_TX_SIZE,
    requiredEnvVars.MAX_L1_TX_SIZE,
    requiredEnvVars.MAX_TX_BATCH_COUNT,
    requiredEnvVars.MAX_BATCH_SUBMISSION_TIME * 1_000,
    requiredEnvVars.NUM_CONFIRMATIONS,
    requiredEnvVars.RESUBMISSION_TIMEOUT * 1_000,
    requiredEnvVars.ADDRESS_MANAGER_ADDRESS,
    requiredEnvVars.SAFE_MINIMUM_ETHER_BALANCE,
    MIN_GAS_PRICE_IN_GWEI,
    MAX_GAS_PRICE_IN_GWEI,
    GAS_RETRY_INCREMENT,
    GAS_THRESHOLD_IN_GWEI,
    BLOCK_OFFSET,
    logger.child({ name: TX_BATCH_SUBMITTER_LOG_TAG }),
    metrics,
    DISABLE_QUEUE_BATCH_APPEND,
    autoFixBatchOptions
  )

  const stateBatchSubmitter = new StateBatchSubmitter(
    proposerSigner,
    l2Provider,
    requiredEnvVars.MIN_L1_TX_SIZE,
    requiredEnvVars.MAX_L1_TX_SIZE,
    requiredEnvVars.MAX_STATE_BATCH_COUNT,
    requiredEnvVars.MAX_BATCH_SUBMISSION_TIME * 1_000,
    requiredEnvVars.NUM_CONFIRMATIONS,
    requiredEnvVars.RESUBMISSION_TIMEOUT * 1_000,
    requiredEnvVars.FINALITY_CONFIRMATIONS,
    requiredEnvVars.ADDRESS_MANAGER_ADDRESS,
    requiredEnvVars.SAFE_MINIMUM_ETHER_BALANCE,
    MIN_GAS_PRICE_IN_GWEI,
    MAX_GAS_PRICE_IN_GWEI,
    GAS_RETRY_INCREMENT,
    GAS_THRESHOLD_IN_GWEI,
    BLOCK_OFFSET,
    logger.child({ name: STATE_BATCH_SUBMITTER_LOG_TAG }),
    metrics,
    FRAUD_SUBMISSION_ADDRESS
  )

  // Loops infinitely!
  const loop = async (
    func: () => Promise<TransactionReceipt>
  ): Promise<void> => {
    // Clear all pending transactions
    if (clearPendingTxs) {
      try {
        const pendingTxs = await sequencerSigner.getTransactionCount('pending')
        const latestTxs = await sequencerSigner.getTransactionCount('latest')
        if (pendingTxs > latestTxs) {
          logger.info(
            'Detected pending transactions. Clearing all transactions!'
          )
          for (let i = latestTxs; i < pendingTxs; i++) {
            const response = await sequencerSigner.sendTransaction({
              to: await sequencerSigner.getAddress(),
              value: 0,
              nonce: i,
            })
            logger.info('Submitted empty transaction', {
              nonce: i,
              txHash: response.hash,
              to: response.to,
              from: response.from,
            })
            logger.debug('empty transaction data', {
              data: response.data,
            })
            await sequencerSigner.provider.waitForTransaction(
              response.hash,
              requiredEnvVars.NUM_CONFIRMATIONS
            )
          }
        }
      } catch (err) {
        logger.error('Cannot clear transactions', {
          message: err.toString(),
          stack: err.stack,
          code: err.code,
        })
        process.exit(1)
      }
    }

    while (true) {
      try {
        await func()
      } catch (err) {
        logger.error('Error submitting batch', {
          message: err.toString(),
          stack: err.stack,
          code: err.code,
        })
        logger.info('Retrying...')
      }
      // Sleep
      await new Promise((r) => setTimeout(r, requiredEnvVars.POLL_INTERVAL))
    }
  }

  // Run batch submitters in two seperate infinite loops!
  if (requiredEnvVars.RUN_TX_BATCH_SUBMITTER) {
    loop(() => txBatchSubmitter.submitNextBatch())
  }
  if (requiredEnvVars.RUN_STATE_BATCH_SUBMITTER) {
    loop(() => stateBatchSubmitter.submitNextBatch())
  }

  if (
    (config.bool('run-prometheus-server'), env.RUN_PROMETHEUS_SERVER === 'true')
  ) {
    // Initialize metrics server
    await createMetricsServer({
      logger,
      registry: metrics.registry,
      port: config.uint('prometheus-port'),
    })
  }
}
