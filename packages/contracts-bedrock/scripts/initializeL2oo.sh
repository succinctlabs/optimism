#!/bin/bash

# Check if --network argument is provided
if [[ "$1" != "--network" || -z "$2" ]]; then
    echo "Usage: $0 --network <network_name>"
    exit 1
fi

network=$2

# Read address from JSON file
address=$(jq -r '.address' "../deployments/${network}/L2OutputOracle.json")

# Read environment variables
rpc_url=$L1_RPC_URL
private_key=$GS_ADMIN_PRIVATE_KEY

# Get output root
# ZTODO: Is the port saved anywhere so we don't have to hardcode this?
output_root=$(cast rpc --rpc-url http://localhost:8547 optimism_outputAtBlock 0x0 | jq -r .outputRoot)

echo "Initial Output Root: $output_root"

# Send transaction
echo "Initializing L2OutputOracle..."
tx_output=$(cast send --rpc-url "$rpc_url" "$address" "setInitialOutputRoot(bytes32)" "$output_root" --private-key "$private_key")
echo "Transaction sent. Waiting for receipt..."

# Extract transaction hash
tx_hash=$(echo "$tx_output" | grep -oP '(?<=tx: )0x[a-fA-F0-9]+')

# Wait for receipt and print it
receipt=$(cast receipt --rpc-url "$rpc_url" "$tx_hash")
echo "L2OutputOracle successfully initialized with Genesis Output Root:"
echo "$receipt"

# Get VK
echo "Getting VK..."
vk=$(./get_vk.bin)

# Send setVKey transaction
echo "Setting VKey..."
vk_tx_output=$(cast send --rpc-url "$rpc_url" "$address" "setVKey(bytes32)" "$vk" --private-key "$private_key")

# Extract transaction hash for setVKey
vk_tx_hash=$(echo "$vk_tx_output" | grep -oP '(?<=tx: )0x[a-fA-F0-9]+')

# Wait for setVKey receipt and print it
echo "Waiting for setVKey transaction receipt..."
vk_receipt=$(cast receipt --rpc-url "$rpc_url" "$vk_tx_hash")
echo "setVKey transaction receipt:"
echo "$vk_receipt"

# Set Verifier Gateway address
verifier_gateway="0x3B6041173B80E77f038f3F2C0f9744f04837185e"

# Send setVerifierGateway transaction
echo "Setting Verifier Gateway..."
vg_tx_output=$(cast send --rpc-url "$rpc_url" "$address" "setVerifierGateway(address)" "$verifier_gateway" --private-key "$private_key")

# Extract transaction hash for setVerifierGateway
vg_tx_hash=$(echo "$vg_tx_output" | grep -oP '(?<=tx: )0x[a-fA-F0-9]+')

# Wait for setVerifierGateway receipt and print it
echo "Waiting for setVerifierGateway transaction receipt..."
vg_receipt=$(cast receipt --rpc-url "$rpc_url" "$vg_tx_hash")
echo "setVerifierGateway transaction receipt:"
echo "$vg_receipt"

# Check Verifier Gateway code size
code_size=$(cast codesize "$verifier_gateway" --rpc-url "$L1_RPC_URL")

if [ "$code_size" -eq 0 ]; then
    echo "WARNING: Verifier Gateway not deployed on this L1. Please request it on the SP1 Telegram: https://t.me/succinct_sp1"
else
