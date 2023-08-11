#!/bin/bash
set -euo pipefail

rm -rf ~/Downloads/data || true

../../op-challenger/bin/op-challenger \
  --l1-eth-rpc http://127.0.0.1:8545 \
  --trace-type cannon \
  --cannon-rollup-config ../../.devnet/rollup.json \
  --cannon-l2-genesis ../../.devnet/genesis-l2.json \
  --cannon-bin ../../cannon/bin/cannon \
  --cannon-server ../../op-program/bin/op-program \
  --cannon-prestate ../../op-program/bin/prestate.json \
  --cannon-datadir ~/Downloads/data \
  --cannon-l2 http://127.0.0.1:9545 \
  --mnemonic "test test test test test test test test test test test junk" \
  --hd-path "m/44'/60'/0'/0/7" \
  --cannon-snapshot-freq 10000000 \
  --num-confirmations 1 \
  --agree-with-proposed-output=true \
  --game-address $1
