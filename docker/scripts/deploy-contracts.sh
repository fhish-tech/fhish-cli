#!/bin/bash
set -e

RPC_URL=${L2_RPC_URL:-"http://minievm:8545"}
PRIVATE_KEY=${DEPLOYER_PRIVATE_KEY:-"0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}
OUT_FILE=${ADDRESSES_FILE:-"/shared/addresses.json"}

echo "==> Waiting for MiniEVM JSON-RPC at $RPC_URL..."
until curl -s -X POST -H "Content-Type: application/json" --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' "$RPC_URL" | grep -q "result"; do
  sleep 2
done

echo "==> MiniEVM JSON-RPC is live, waiting 10s for full readiness..."
sleep 10

# Create a robust JS config
cat <<EOF > hardhat.config.docker.js
require('@nomicfoundation/hardhat-toolbox');
module.exports = {
  solidity: {
    version: '0.8.24',
    settings: {
      optimizer: { enabled: true, runs: 200 },
      evmVersion: 'cancun'
    }
  },
  networks: {
    minievm: {
      url: '$RPC_URL',
      accounts: ['$PRIVATE_KEY'],
    },
  },
};
EOF

# Patch the deploy script
cp scripts/deploy.js scripts/deploy_docker.js

# Patch FhishGateway (needs 2 args: admin, kmsVerifier)
sed -i 's/const gateway = await FhishGateway.deploy();/const gateway = await FhishGateway.deploy(deployer.address, "0x000000000000000000000000000000000000005e");/g' scripts/deploy_docker.js

# Patch PrivateVotingV2 (needs 1 arg: gateway)
sed -i 's/const voting = await PrivateVotingV2.deploy(gatewayAddress, aclAddress, executorAddress);/const voting = await PrivateVotingV2.deploy(gatewayAddress);/g' scripts/deploy_docker.js

echo "==> Running patched deployment script..."
npx hardhat run scripts/deploy_docker.js --network minievm --config hardhat.config.docker.js > deploy_output.txt 2>&1
cat deploy_output.txt

# Extract addresses
GATEWAY=$(grep "GATEWAY_ADDRESS=" deploy_output.txt | cut -d'=' -f2)
VOTING=$(grep "VOTING_ADDRESS=" deploy_output.txt | cut -d'=' -f2)

if [ -n "$GATEWAY" ] && [ -n "$VOTING" ]; then
  echo "{
    \"FhishGateway\": \"$GATEWAY\",
    \"PrivateVotingV2\": \"$VOTING\"
  }" > "$OUT_FILE"
  echo "==> Successfully deployed and saved addresses."
else
  echo "==> Deployment failed or addresses not found."
  exit 1
fi
