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
