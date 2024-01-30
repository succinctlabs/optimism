#!/bin/bash
set -euo pipefail

if [ -z "${DEVNET_ADMIN}" ]; then
    echo "Error: Set DEVNET_ADMIN to secret key of 0x8c20c40180751d93E939DDDee3517AE0d1EBeAd2."
    exit 1
fi
if [ -z "${ETH_RPC_URL}" ]; then
    echo "Error: Set ETH_RPC_URL to Sepolia L1 rpc endpoint."
    exit 1
fi

forge script scripts/SystemConfig-SetGasConfig-devnet.s.sol \
  --private-key $DEVNET_ADMIN \
  --rpc-url $ETH_RPC_URL \
  -vvvvv \
  # --broadcast # Uncomment to send the tx. Without it, you can simulate and check the expected results.
