#!/bin/bash
set -e

RPC_URL=${L2_RPC_URL:-"http://minievm:8545"}
PRIVATE_KEY=${DEPLOYER_PRIVATE_KEY:-"0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}
OUT_FILE=${ADDRESSES_FILE:-"/shared/addresses.json"}
CHAIN_ID=${EVM_CHAIN_ID:-"1234"}

echo "==> Waiting for MiniEVM RPC at $RPC_URL..."
until cast block-number --rpc-url "$RPC_URL" 2>/dev/null | grep -q "[0-9]"; do
  echo "   ...waiting"
  sleep 3
done

BLOCK=$(cast block-number --rpc-url "$RPC_URL")
echo "==> MiniEVM is live at block $BLOCK"

echo "==> Deploying FHE contracts..."

# Find the deploy script
DEPLOY_SCRIPT=$(find ./scripts -name "deploy.js" | head -1)

if [ -n "$DEPLOY_SCRIPT" ]; then
  echo "==> Using hardhat deploy script: $DEPLOY_SCRIPT"
  # Set up env for hardhat
  export RPC_URL=$RPC_URL
  export PRIVATE_KEY=$PRIVATE_KEY
  
  # Run hardhat deployment
  # Note: assuming hardhat is in node_modules
  npx hardhat run "$DEPLOY_SCRIPT" --network localhost > deploy_output.txt 2>&1
  cat deploy_output.txt
  
  # Extract addresses from output (naive grep for demo)
  GATEWAY=$(grep "FhishGateway deployed to:" deploy_output.txt | awk '{print $NF}')
  VOTING=$(grep "PrivateVotingV2 deployed to:" deploy_output.txt | awk '{print $NF}')
  
  # If hardhat output is clean, we can do better. For now, we'll write what we found.
  echo "{
    \"FhishGateway\": \"$GATEWAY\",
    \"PrivateVotingV2\": \"$VOTING\"
  }" > "$OUT_FILE"
else
  # Fallback to individual deployment if no script
  echo "==> No deploy script found, skipping automated deployment."
  echo "{}" > "$OUT_FILE"
fi

echo "==> Deployed contracts:"
cat "$OUT_FILE"
echo "==> Addresses written to $OUT_FILE"
