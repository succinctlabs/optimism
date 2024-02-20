#!/bin/bash
set -euo pipefail

if [ -z "${ADMIN_PK}" ]; then
    echo "Error: Set ADMIN_PK to secret key of 0xfd1D2e729aE8eEe2E146c033bf4400fE75284301."
    exit 1
fi
if [ -z "${ETH_RPC_URL}" ]; then
    echo "Error: Set ETH_RPC_URL to Sepolia L1 rpc endpoint."
    exit 1
fi

forge script scripts/SystemConfig-SetGasConfig-devnet.s.sol \
  --private-key $ADMIN_PK\
  --rpc-url $ETH_RPC_URL \
  -vvvvv \
  # --broadcast # Uncomment to send the tx. Without it, you can simulate and check the expected results.
