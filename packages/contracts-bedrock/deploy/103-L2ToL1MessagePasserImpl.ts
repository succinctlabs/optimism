import { DeployFunction } from 'hardhat-deploy/dist/types'
import '@eth-optimism/hardhat-deploy-config'
import '@nomiclabs/hardhat-ethers'

import { assertContractVariable, deploy } from '../src/deploy-utils'

const deployFn: DeployFunction = async (hre) => {
  return // added by kz for multisig rehearsal.
  await deploy({
    hre,
    name: 'L2ToL1MessagePasser',
    args: [],
    postDeployAction: async (contract) => {
      await assertContractVariable(contract, 'MESSAGE_VERSION', 1)
    },
  })
}

deployFn.tags = ['L2ToL1MessagePasserImpl', 'l2']

export default deployFn
