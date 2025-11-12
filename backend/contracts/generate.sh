#!/bin/bash

# Generate Go bindings for InsurancePool contract
# Make sure you have abigen installed: go install github.com/ethereum/go-ethereum/cmd/abigen@latest

CONTRACTS_DIR="../contracts/artifacts/contracts"
OUTPUT_DIR="."

echo "Generating Go bindings for InsurancePool..."

abigen --abi="${CONTRACTS_DIR}/InsurancePool.sol/InsurancePool.json" \
       --pkg=contracts \
       --type=InsurancePool \
       --out="${OUTPUT_DIR}/insurancepool.go"

if [ $? -eq 0 ]; then
    echo "✅ Successfully generated insurancepool.go"
else
    echo "❌ Failed to generate bindings. Make sure abigen is installed:"
    echo "   go install github.com/ethereum/go-ethereum/cmd/abigen@latest"
    exit 1
fi
