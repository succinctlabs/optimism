const config = {
  numDeployConfirmations: 4,
  gasPrice: 45_000_000_000,
  l1BlockTimeSeconds: 15,
  l2BlockGasLimit: 15_000_000,
  l2ChainId: 10,
  ctcL2GasDiscountDivisor: 32,
  ctcEnqueueGasCost: 60_000,
  sccFaultProofWindowSeconds: 604800,
  sccSequencerPublishWindowSeconds: 12592000,
  ovmSequencerAddress: '0x78339d822c23D943E4a2d4c3DD5408F66e6D662D',
  ovmProposerAddress: '0x78339d822c23D943E4a2d4c3DD5408F66e6D662D',
  ovmBlockSignerAddress: '0x78339d822c23D943E4a2d4c3DD5408F66e6D662D',
  ovmFeeWalletAddress: '0x78339d822c23D943E4a2d4c3DD5408F66e6D662D',
  // the one that matters
  ovmAddressManagerOwner: '0x78339d822c23D943E4a2d4c3DD5408F66e6D662D',
  ovmGasPriceOracleOwner: '0x78339d822c23D943E4a2d4c3DD5408F66e6D662D',
  ovmWhitelistOwner: '0x78339d822c23D943E4a2d4c3DD5408F66e6D662D',
}

export default config
