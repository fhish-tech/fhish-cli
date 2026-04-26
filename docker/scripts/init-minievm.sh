#!/bin/bash
set -e

HOME_DIR=${HOME_DIR:-"/data/minievm"}
GAS_DENOM=${GAS_DENOM:-"uinit"}

# The Cosmos address corresponding to Hardhat's first account (0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266)
HARDHAT_COSMOS_ADDR="init17w0adeg64ky0daxwd2ugyuneellmjgnxdtmpqz"

# Force re-init
rm -rf "$HOME_DIR"/*
echo "==> Initializing MiniEVM node..."
minievm init "node" --chain-id "fhish-1" --home "$HOME_DIR"

# Force empty blocks and faster production
sed -i 's/create_empty_blocks = false/create_empty_blocks = true/g' "$HOME_DIR/config/config.toml"
sed -i 's/timeout_propose = "3s"/timeout_propose = "1s"/g' "$HOME_DIR/config/config.toml"
sed -i 's/timeout_commit = "5s"/timeout_commit = "1s"/g' "$HOME_DIR/config/config.toml"

# Create key for validator
echo "==> Creating validator key..."
minievm keys add validator --keyring-backend test --home "$HOME_DIR"
VAL_ADDR=$(minievm keys show validator -a --keyring-backend test --home "$HOME_DIR")

# Add Hardhat account and validator to genesis
echo "==> Adding accounts to genesis..."
minievm genesis add-genesis-account "$HARDHAT_COSMOS_ADDR" "1000000000000$GAS_DENOM" --home "$HOME_DIR"
minievm genesis add-genesis-account "$VAL_ADDR" "1000000000000$GAS_DENOM" --home "$HOME_DIR"

# Patch genesis (Python)
echo "==> Patching genesis..."
python3 -c "
import json
path = '$HOME_DIR/config/genesis.json'
with open(path, 'r') as f:
    data = json.load(f)
data['app_state']['opchild']['params']['bridge_executors'] = ['$HARDHAT_COSMOS_ADDR']
data['app_state']['opchild']['params']['admin'] = '$HARDHAT_COSMOS_ADDR'
data['app_state']['evm']['params']['fee_denom'] = '$GAS_DENOM'
if 'consensus_params' in data:
    data['consensus_params']['block']['max_gas'] = '100000000000'
with open(path, 'w') as f:
    json.dump(data, f)
"

# Add genesis validator
echo "==> Adding genesis validator..."
minievm genesis add-genesis-validator validator --keyring-backend test --home "$HOME_DIR"

# Enable CORS for RPC and API
sed -i 's/cors_allowed_origins = \[\]/cors_allowed_origins = ["*"]/g' "$HOME_DIR/config/config.toml"
sed -i 's/enabled-unsafe-cors = false/enabled-unsafe-cors = true/g' "$HOME_DIR/config/app.toml"

echo "==> Starting MiniEVM node..."
minievm start --home "$HOME_DIR" \
    --rpc.laddr tcp://0.0.0.0:26657 \
    --api.enable true \
    --api.enabled-unsafe-cors true \
    --api.address "tcp://0.0.0.0:1317" \
    --indexer.enable true \
    --json-rpc.address 0.0.0.0:8545 \
    --json-rpc.enable true \
    --json-rpc.enable-unsafe-cors true &

# Wait for node to be live
until curl -s http://localhost:26657/status | grep -q "latest_block_height"; do
  sleep 2
done

echo "==> Node is live, initializing Hardhat account..."
minievm tx bank send "$VAL_ADDR" "$HARDHAT_COSMOS_ADDR" "1$GAS_DENOM" --chain-id "fhish-1" --keyring-backend test --home "$HOME_DIR" --yes --broadcast-mode sync || echo "WARNING: account initialization failed"
sleep 5

echo "==> MiniEVM ready."
wait
