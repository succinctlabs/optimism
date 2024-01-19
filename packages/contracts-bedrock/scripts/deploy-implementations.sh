#!/usr/bin/env bash

# Used to deploy only the implementations of L1 contracts
set -euo pipefail

verify_flag=""
if [ -n "${DEPLOY_VERIFY:-}" ]; then
  verify_flag="--verify"
fi

SENDER=$(cast wallet address --keystore "$DEPLOY_PRIVATE_KEY")
echo "Sender: $SENDER"

echo "> Deploying contracts"
forge script -vvv scripts/Deploy.s.sol:DeployExtendedPauseUpgrade \
--rpc-url "$DEPLOY_ETH_RPC_URL"  \
--keystore "$DEPLOY_PRIVATE_KEY" \
--sender "$SENDER" \
--broadcast \
  $verify_flag

# if [ -n "${DEPLOY_GENERATE_HARDHAT_ARTIFACTS:-}" ]; then
#   echo "> Generating hardhat artifacts"
#   forge script -vvv scripts/Deploy.s.sol:Deploy --sig 'sync()' --rpc-url "$DEPLOY_ETH_RPC_URL"  \
#   --broadcast \
#   --keystore "$DEPLOY_PRIVATE_KEY"
# fi
