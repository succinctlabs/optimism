#!/bin/bash
set -euo pipefail

cast send  \
  --mnemonic "test test test test test test test test test test test junk" \
  --mnemonic-derivation-path "m/44'/60'/0'/0/4" \
  --chain 900  \
  $1 'attack(uint256,bytes32)' $2 0xaa00000000000000000000000000000000000000000000000000000000000000
