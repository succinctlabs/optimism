/* Externals Import */
import {
  add0x,
  getDeployedContractAddress,
  getLogger,
} from '@eth-optimism/core-utils'
import {
  GAS_LIMIT,
  L2ExecutionManagerContractDefinition,
  L2ToL1MessagePasserContractDefinition,
  DEFAULT_OPCODE_WHITELIST_MASK,
} from '@eth-optimism/ovm'
import { Address } from '@eth-optimism/rollup-core'

import { Contract, Wallet } from 'ethers'
import { createMockProvider, getWallets } from 'ethereum-waffle'
import { readFile } from 'fs'

/* Internal Imports */
import { DEFAULT_ETHNODE_GAS_LIMIT, deployContract } from '../index'
import { JsonRpcProvider } from 'ethers/providers'
import { promisify } from 'util'

const readFileAsync = promisify(readFile)
const log = getLogger('l2-node')

/* Configuration */
const opcodeWhitelistMask: string =
  process.env.OPCODE_WHITELIST_MASK || DEFAULT_OPCODE_WHITELIST_MASK
const volumePath: string = process.env.VOLUME_PATH || '/'
const privateKeyFilePath: string =
  process.env.PRIVATE_KEY_FILE_PATH || volumePath + '/private_key.txt'

export interface L2NodeContext {
  provider: JsonRpcProvider
  wallet: Wallet
  executionManager: Contract
  l2ToL1MessagePasser: Contract
}

export async function initializeL2Node(
  web3Provider?: JsonRpcProvider
): Promise<L2NodeContext> {
  let provider: JsonRpcProvider = web3Provider
  if (!web3Provider) {
    provider = createMockProvider({
      gasLimit: DEFAULT_ETHNODE_GAS_LIMIT,
      allowUnlimitedContractSize: true,
    })
  }

  // Initialize a fullnode for us to interact with
  let wallet

  // If we're given a provider, our wallet must be configured from a private key file
  if (web3Provider) {
    const privateKey: string = await readFileAsync(privateKeyFilePath, 'utf8')
    wallet = new Wallet(add0x(privateKey), provider)
  } else {
    ;[wallet] = getWallets(provider)
  }

  let nonce: number = 0
  const executionManagerAddress: Address = await getDeployedContractAddress(
    nonce++,
    provider,
    wallet.address
  )

  let executionManager: Contract
  let l2ToL1MessagePasser: Contract
  if (executionManagerAddress) {
    log.info(
      `Using existing ExecutionManager deployed at ${executionManagerAddress}`
    )
    executionManager = new Contract(
      executionManagerAddress,
      L2ExecutionManagerContractDefinition.abi,
      wallet
    )
  } else {
    executionManager = await deployExecutionManager(wallet)
  }

  const messagePasserAddress: Address = await getDeployedContractAddress(
    nonce++,
    provider,
    wallet.address
  )
  if (!messagePasserAddress) {
    l2ToL1MessagePasser = await deployL2ToL1MessagePasser(
      wallet,
      executionManager.address
    )
  } else {
    log.info(
      `Using existing L2ToL1MessagePasser deployed at ${messagePasserAddress}`
    )
    l2ToL1MessagePasser = new Contract(
      messagePasserAddress,
      L2ToL1MessagePasserContractDefinition.abi,
      wallet
    )
  }

  return {
    wallet,
    provider,
    executionManager,
    l2ToL1MessagePasser,
  }
}

/**
 * Deploys the ExecutionManager contract with the provided wallet and whitelist,
 * returning the resulting Contract.
 *
 * @param wallet The wallet to be used, containing all connection info.
 * @returns The deployed Contract.
 */
export async function deployExecutionManager(
  wallet: Wallet
): Promise<Contract> {
  log.debug('Deploying execution manager...')

  const executionManager: Contract = await deployContract(
    wallet,
    L2ExecutionManagerContractDefinition,
    [opcodeWhitelistMask, wallet.address, GAS_LIMIT, true],
    { gasLimit: DEFAULT_ETHNODE_GAS_LIMIT }
  )

  log.info('Deployed execution manager to address:', executionManager.address)

  return executionManager
}

/**
 * Deploys the L2ToL1MessagePasser contract with the provided wallet and EM address,
 * returning the resulting Contract.
 *
 * @param wallet The wallet to be used, containing all connection info.
 * @param executionManagerAddress The EM address param to the L2ToL1MessagePasser.
 * @returns The deployed Contract.
 */
export async function deployL2ToL1MessagePasser(
  wallet: Wallet,
  executionManagerAddress: Address
): Promise<Contract> {
  log.debug('Deploying L2ToL1MessagePasser contract...')

  const l2ToL1MessagePasser: Contract = await deployContract(
    wallet,
    L2ToL1MessagePasserContractDefinition,
    [executionManagerAddress],
    { gasLimit: DEFAULT_ETHNODE_GAS_LIMIT }
  )

  log.info(
    'Deployed L2ToL1MessagePasser to address:',
    l2ToL1MessagePasser.address
  )

  return l2ToL1MessagePasser
}
