import { expect } from '../../../setup'

/* External Imports */
import { ethers, waffle } from 'hardhat'
import { ContractFactory, Contract, Wallet, Signer } from 'ethers'
import { MockContract, smockit } from '@eth-optimism/smock'

/* Internal Imports */
import { getContractInterface } from '../../../../src'
import { toPlainObject } from 'lodash'

describe('OVM_ProxyEOA', () => {
  const eoaDefaultAddr = '0x4200000000000000000000000000000000000003'

  let wallet: Wallet
  before(async () => {
    const provider = waffle.provider
    ;[wallet] = provider.getWallets()
  })

  let signer: Signer
  before(async () => {
    ;[signer] = await ethers.getSigners()
  })

  let Mock__OVM_ECDSAContractAccount: MockContract
  before(async () => {
    Mock__OVM_ECDSAContractAccount = await smockit(
      getContractInterface('OVM_ECDSAContractAccount'),
      {
        address: eoaDefaultAddr,
      }
    )
  })

  let Factory__OVM_ProxyEOA: ContractFactory
  before(async () => {
    Factory__OVM_ProxyEOA = await ethers.getContractFactory('OVM_ProxyEOA')
  })

  let OVM_ProxyEOA: Contract
  beforeEach(async () => {
    OVM_ProxyEOA = await Factory__OVM_ProxyEOA.deploy()
  })

  describe('getImplementation()', () => {
    it(`should be created with implementation at predeploy address`, async () => {
      expect(await OVM_ProxyEOA.getImplementation()).to.equal(eoaDefaultAddr)
    })
  })

  describe('upgrade()', () => {
    // TODO: How do we test this?
    it.skip(`should upgrade the proxy implementation`, async () => {
      const newImpl = `0x${'81'.repeat(20)}`

      await expect(OVM_ProxyEOA.upgrade(newImpl)).to.not.be.reverted

      expect(await OVM_ProxyEOA.getImplementation()).to.equal(newImpl)
    })

    it(`should not allow upgrade of the proxy implementation by another account`, async () => {
      const newImpl = `0x${'81'.repeat(20)}`

      await expect(OVM_ProxyEOA.upgrade(newImpl)).to.be.revertedWith(
        'EOAs can only upgrade their own EOA implementation'
      )
    })
  })

  describe('fallback()', () => {
    it(`should call delegateCall with right calldata`, async () => {
      const data = Mock__OVM_ECDSAContractAccount.interface.encodeFunctionData(
        'execute',
        ['0x12341234']
      )

      await signer.sendTransaction({
        to: OVM_ProxyEOA.address,
        data,
      })

      expect(
        toPlainObject(Mock__OVM_ECDSAContractAccount.smocked.execute.calls[0])
      ).to.deep.include({
        _encodedTransaction: '0x12341234',
      })
    })

    it.skip(`should return data from fallback`, async () => {
      // TODO: test return data from fallback
    })

    it.skip(`should revert in fallback`, async () => {
      // TODO: test reversion from fallback
    })
  })
})
