#!/bin/bash
set -euo pipefail

if [ -z "${DEVNET_ADMIN}" ]; then
    echo "Error: Set DEVNET_ADMIN to secret key of 0x858F0751ef8B4067f0d2668C076BDB50a8549fbF."
    exit 1
fi
if [ -z "${ETH_RPC_URL}" ]; then
    echo "Error: Set ETH_RPC_URL to Goerli L1 rpc endpoint."
    exit 1
fi

forge script scripts/SystemConfig-SetGasConfig-devnet.s.sol \
  --private-key $DEVNET_ADMIN \
  --rpc-url $ETH_RPC_URL \
  -vvvvv \
  # --broadcast # Uncomment to send the tx. Without it, you can simulate and check the expected results.
