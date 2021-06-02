/* Imports: External */
import { DeployFunction } from 'hardhat-deploy/dist/types'

/* Imports: Internal */
import { getDeployedContract } from '../src/hardhat-deploy-ethers'
import { predeploys } from '../src/predeploys'

const deployFn: DeployFunction = async (hre) => {
  const { deploy } = hre.deployments
  const { deployer } = await hre.getNamedAccounts()

  const Lib_AddressManager = await getDeployedContract(
    hre,
    'Lib_AddressManager',
    {
      signerOrProvider: deployer,
    }
  )

  const result = await deploy('Proxy__OVM_L1StandardBridge', {
    contract: 'Lib_ResolvedDelegateProxy',
    from: deployer,
    args: [Lib_AddressManager.address, 'OVM_L1StandardBridge'],
    log: true,
  })

  if (!result.newlyDeployed) {
    return
  }

  const Proxy__OVM_L1StandardBridge = await getDeployedContract(
    hre,
    'Proxy__OVM_L1StandardBridge',
    {
      signerOrProvider: deployer,
      iface: 'OVM_L1StandardBridge',
    }
  )

  await Proxy__OVM_L1StandardBridge.initialize(
    Lib_AddressManager.address,
    predeploys.OVM_ETH
  )

  const libAddressManager = await Proxy__OVM_L1StandardBridge.libAddressManager()
  if (libAddressManager !== Lib_AddressManager.address) {
    throw new Error(
      `\n**FATAL ERROR. THIS SHOULD NEVER HAPPEN. CHECK YOUR DEPLOYMENT.**:\n` +
        `Proxy__OVM_L1StandardBridge could not be succesfully initialized.\n` +
        `Attempted to set Lib_AddressManager to: ${Lib_AddressManager.address}\n` +
        `Actual address after initialization: ${libAddressManager}\n` +
        `This could indicate a compromised deployment.`
    )
  }

  await Lib_AddressManager.setAddress('Proxy__OVM_L1StandardBridge', result.address)
}

deployFn.dependencies = ['Lib_AddressManager', 'OVM_L1StandardBridge']
deployFn.tags = ['Proxy__OVM_L1StandardBridge']

export default deployFn
