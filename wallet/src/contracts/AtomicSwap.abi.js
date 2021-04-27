module.exports = [
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": false,
        "internalType": "bytes32",
        "name": "_swapID",
        "type": "bytes32"
      }
    ],
    "name": "Close",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": false,
        "internalType": "bytes32",
        "name": "_swapID",
        "type": "bytes32"
      }
    ],
    "name": "Expire",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": false,
        "internalType": "bytes32",
        "name": "_swapID",
        "type": "bytes32"
      },
      {
        "indexed": false,
        "internalType": "address",
        "name": "_closeTrader",
        "type": "address"
      }
    ],
    "name": "Open",
    "type": "event"
  },
  {
    "inputs": [
      {
        "internalType": "bytes32",
        "name": "_swapID",
        "type": "bytes32"
      }
    ],
    "name": "check",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "openValue",
        "type": "uint256"
      },
      {
        "internalType": "address",
        "name": "openContractAddress",
        "type": "address"
      },
      {
        "internalType": "uint256",
        "name": "closeValue",
        "type": "uint256"
      },
      {
        "internalType": "address",
        "name": "closeTrader",
        "type": "address"
      },
      {
        "internalType": "address",
        "name": "closeContractAddress",
        "type": "address"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "bytes32",
        "name": "_swapID",
        "type": "bytes32"
      }
    ],
    "name": "close",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "bytes32",
        "name": "_swapID",
        "type": "bytes32"
      }
    ],
    "name": "expire",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "bytes32",
        "name": "_swapID",
        "type": "bytes32"
      },
      {
        "internalType": "uint256",
        "name": "_openValue",
        "type": "uint256"
      },
      {
        "internalType": "address",
        "name": "_openContractAddress",
        "type": "address"
      },
      {
        "internalType": "uint256",
        "name": "_closeValue",
        "type": "uint256"
      },
      {
        "internalType": "address",
        "name": "_closeTrader",
        "type": "address"
      },
      {
        "internalType": "address",
        "name": "_closeContractAddress",
        "type": "address"
      }
    ],
    "name": "open",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  }
];