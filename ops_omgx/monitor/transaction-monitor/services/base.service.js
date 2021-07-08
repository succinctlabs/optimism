#!/usr/bin/env node

const ethers = require('ethers');
const util = require('util');
const core_utils_1 = require('@eth-optimism/core-utils');

const Watcher = require('./utilities/watcher');

const addressManagerJSON = require('../artifacts/contracts/optimistic-ethereum/libraries/resolver/Lib_AddressManager.sol/Lib_AddressManager.json');

require('dotenv').config();
const env = process.env;
const L1_NODE_WEB3_URL = env.L1_NODE_WEB3_URL || "http://localhost:8545";
const L2_NODE_WEB3_URL = env.L2_NODE_WEB3_URL || "http://localhost:9545";

const MYSQL_HOST_URL = env.MYSQL_HOST_URL || "127.0.0.1";
const MYSQL_PORT = env.MYSQL_PORT || 3306;
const MYSQL_USERNAME = env.MYSQL_USERNAME;
const MYSQL_PASSWORD = env.MYSQL_PASSWORD;
const MYSQL_DATABASE_NAME = env.MYSQL_DATABASE_NAME || "OMGXV1";

const ADDRESS_MANAGER_ADDRESS = env.ADDRESS_MANAGER_ADDRESS;
const L2_MESSENGER_ADDRESS = env.L2_MESSENGER_ADDRESS || "0x4200000000000000000000000000000000000007";

const DEPLOYER_PRIVATE_KEY = env.DEPLOYER_PRIVATE_KEY;

const CHAIN_SCAN_INTERVAL = env.CHAIN_SCAN_INTERVAL || 60000;
const MESSAGE_SCAN_INTERVAL = env.MESSAGE_SCAN_INTERVAL || 60 * 60 * 1000;

class BaseService extends Watcher {
  constructor() {
    super(...arguments);
    this.L1Provider = new ethers.providers.JsonRpcProvider(L1_NODE_WEB3_URL);
    this.L2Provider = new ethers.providers.JsonRpcProvider(L2_NODE_WEB3_URL);

    this.wallet = new ethers.Wallet(DEPLOYER_PRIVATE_KEY).connect(this.L1Provider);

    this.MySQLHostURL = MYSQL_HOST_URL;
    this.MySQLPort = MYSQL_PORT;
    this.MySQLUsername = MYSQL_USERNAME;
    this.MySQLPassword = MYSQL_PASSWORD;
    this.MySQLDatabaseName = MYSQL_DATABASE_NAME;

    this.addressManagerAddress = ADDRESS_MANAGER_ADDRESS;
    this.L1MessengerAddress = null;
    this.L2MessengerAddress = L2_MESSENGER_ADDRESS;

    this.numberBlockToFetch = 10000;
    this.chainScanInterval = CHAIN_SCAN_INTERVAL;
    this.messageScanInterval = MESSAGE_SCAN_INTERVAL;

    this.logger = new core_utils_1.Logger({ name: this.name });
    this.sleep = util.promisify(setTimeout);
  }

  async initBaseService() {
    const addressManager = new ethers.Contract(
      this.addressManagerAddress,
      addressManagerJSON.abi,
      this.wallet
    )
    this.L1MessengerAddress = await addressManager.getAddress('Proxy__OVM_L1CrossDomainMessenger');
  }
}

module.exports = BaseService;