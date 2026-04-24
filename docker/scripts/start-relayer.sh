#!/bin/sh
set -e

# Wait for the deployer to write addresses
ADDRESSES_FILE=${ADDRESSES_FILE:-"/shared/addresses.json"}
echo "==> Waiting for contract addresses at $ADDRESSES_FILE..."
until [ -f "$ADDRESSES_FILE" ]; do
  sleep 2
done

# Export contract addresses as env vars
export VOTING_ADDRESS=$(jq -r '.PrivateVotingV2 // .PrivateVoting // .privateVoting' "$ADDRESSES_FILE")
echo "==> Voting Contract: $VOTING_ADDRESS"

# Start the relayer
if [ -f "/app/dist/index.js" ]; then
  exec node /app/dist/index.js
elif [ -f "/app/src/index.ts" ]; then
  exec npx tsx src/index.ts
else
  # Find main entry from package.json
  MAIN=$(node -e "const pkg = require('./package.json'); console.log(pkg.main || 'index.js')")
  exec node "$MAIN"
fi
