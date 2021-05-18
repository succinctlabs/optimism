import { expect } from '../../../setup'

/* External Imports */
import { ethers } from 'hardhat'
import { ContractFactory, Contract, Signer } from 'ethers'

describe('OVM_SequencerEntrypoint', () => {
  let signer1: Signer
  let signer2: Signer
  before(async () => {
    ;[signer1, signer2] = await ethers.getSigners()
  })

  let Factory__OVM_CongestionPriceOracle: ContractFactory
  before(async () => {
    Factory__OVM_CongestionPriceOracle = await ethers.getContractFactory(
      'OVM_CongestionPriceOracle'
    )
  })

  let OVM_CongestionPriceOracle: Contract
  beforeEach(async () => {
    OVM_CongestionPriceOracle = await Factory__OVM_CongestionPriceOracle.deploy(
      await signer1.getAddress()
    )
  })

  describe('owner', () => {
    it('should have an owner', async () => {
      expect(await OVM_CongestionPriceOracle.owner()).to.equal(
        await signer1.getAddress()
      )
    })
  })

  describe('setCongestionPrice', () => {
    it('should revert if called by someone other than the owner', async () => {
      await expect(
        OVM_CongestionPriceOracle.connect(signer2).setCongestionPrice(1234)
      ).to.be.reverted
    })

    it('should succeed if called by the owner', async () => {
      await expect(
        OVM_CongestionPriceOracle.connect(signer1).setCongestionPrice(1234)
      ).to.not.be.reverted
    })
  })

  describe('getCongestionPrice', () => {
    it('should return zero at first', async () => {
      expect(await OVM_CongestionPriceOracle.getCongestionPrice()).to.equal(0)
    })

    it('should change when setCongestionPrice is called', async () => {
      const congestionPrice = 1234

      await OVM_CongestionPriceOracle.connect(signer1).setCongestionPrice(
        congestionPrice
      )

      expect(await OVM_CongestionPriceOracle.getCongestionPrice()).to.equal(
        congestionPrice
      )
    })
  })
})
