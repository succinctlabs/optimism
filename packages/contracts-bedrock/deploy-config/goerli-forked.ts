import { DeployConfig } from '../src/deploy-config'

const config: DeployConfig = {
  // Core config
  finalSystemOwner: '0x62790eFcB3a5f3A5D398F95B47930A9Addd83807',
  l1StartingBlockTag:
    '0xcbcbf93fa6ef953a491466311fce85846252a808aa3cc6a8fd432afc5ea4487b',
  l1ChainID: 5,
  l2ChainID: 420,
  l2BlockTime: 2,
  maxSequencerDrift: 1200,
  sequencerWindowSize: 3600,
  channelTimeout: 120,
  p2pSequencerAddress: '0xCBABF46d40982B4530c0EAc9889f6e44e17f0383',
  batchInboxAddress: '0xff00000000000000000000000000000000000420',
  batchSenderAddress: '0x3a2baA0160275024A50C1be1FC677375E7DB4Bd7',
  l2OutputOracleSubmissionInterval: 20,
  l2OutputOracleStartingTimestamp: 1670043996,
  l2OutputOracleProposer: '0x88BCa4Af3d950625752867f826E073E337076581',
  l2OutputOracleChallenger: '0x88BCa4Af3d950625752867f826E073E337076581',
  finalizationPeriodSeconds: 2,

  // L2 network config
  l2GenesisBlockGasLimit: '0x17D7840',
  l2GenesisBlockCoinbase: '0x4200000000000000000000000000000000000011',
  l2GenesisBlockBaseFeePerGas: '0x3b9aca00',
  gasPriceOracleOverhead: 2100,
  gasPriceOracleScalar: 1000000,
}

export default config
