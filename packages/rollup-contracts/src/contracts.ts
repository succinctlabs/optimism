/* External Imports */
import { ethers } from 'ethers'

/* Contract Imports */
import * as ExecutionManager from '../build/contracts/ExecutionManager.json'
import * as FullStateManager from '../build/contracts/FullStateManager.json'
import * as L2ExecutionManager from '../build/contracts/L2ExecutionManager.json'
import * as ContractAddressGenerator from '../build/contracts/ContractAddressGenerator.json'
import * as L2ToL1MessageReceiver from '../build/contracts/L2ToL1MessageReceiver.json'
import * as L2ToL1MessagePasser from '../build/contracts/L2ToL1MessagePasser.json'
import * as L1ToL2TransactionPasser from '../build/contracts/L1ToL2TransactionPasser.json'
import * as RLPWriter from '../build/contracts/RLPWriter.json'
import * as SafetyChecker from '../build/contracts/SafetyChecker.json'

/* Contract Exports */
export const ExecutionManagerContractDefinition = ExecutionManager
export const L2ExecutionManagerContractDefinition = L2ExecutionManager
export const FullStateManagerContractDefinition = FullStateManager
export const ContractAddressGeneratorContractDefinition = ContractAddressGenerator
export const L2ToL1MessageReceiverContractDefinition = L2ToL1MessageReceiver
export const L2ToL1MessagePasserContractDefinition = L2ToL1MessagePasser
export const L1ToL2TransactionPasserContractDefinition = L1ToL2TransactionPasser
export const RLPWriterContractDefinition = RLPWriter
export const SafetyCheckerContractDefinition = SafetyChecker

export const executionManagerInterface = new ethers.utils.Interface(
  ExecutionManager.interface
)

export const l2ExecutionManagerInterface = new ethers.utils.Interface(
  L2ExecutionManager.interface
)
export const l2ToL1MessagePasserInterface = new ethers.utils.Interface(
  L2ToL1MessagePasser.interface
)
