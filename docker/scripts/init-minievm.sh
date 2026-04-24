#!/bin/bash
set -e

CHAIN_ID=${CHAIN_ID:-"fhish-1"}
EVM_CHAIN_ID=${EVM_CHAIN_ID:-"1234"}
MONIKER=${MONIKER:-"fhish-node"}
HOME_DIR=${HOME_DIR:-"/data/minievm"}
GAS_DENOM=${GAS_DENOM:-"uinit"}

# Only init if not already initialized
if [ ! -f "$HOME_DIR/config/genesis.json" ]; then
  echo "==> Initializing MiniEVM node..."
  minievm init "$MONIKER" --chain-id "$CHAIN_ID" --home "$HOME_DIR"

  # Patch genesis: set EVM chain ID and high gas limits for FHE
  jq --arg evm_id "$EVM_CHAIN_ID" --arg denom "$GAS_DENOM" \
     '.app_state.evm.params.evm_denom = $denom |
      .app_state.evm.params.allowed_publishers = [] |
      .consensus_params.block.max_gas = "100000000000"' \
     "$HOME_DIR/config/genesis.json" > /tmp/genesis_patched.json
  mv /tmp/genesis_patched.json "$HOME_DIR/config/genesis.json"

  # Add a funded genesis account for contract deployer
  DEPLOYER_ADDR=${DEPLOYER_ADDRESS:-"0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"}
  if [ -n "$DEPLOYER_ADDR" ]; then
    echo "==> Adding genesis account: $DEPLOYER_ADDR"
    minievm genesis add-genesis-account "$DEPLOYER_ADDR" "1000000000000$GAS_DENOM" --home "$HOME_DIR"
  fi

  # Patch app.toml: enable JSON-RPC
  sed -i 's/enable = false/enable = true/g' "$HOME_DIR/config/app.toml"
  sed -i 's|address = "127.0.0.1:8545"|address = "0.0.0.0:8545"|g' "$HOME_DIR/config/app.toml"
  sed -i 's|ws-address = "127.0.0.1:8546"|ws-address = "0.0.0.0:8546"|g' "$HOME_DIR/config/app.toml"
  sed -i 's/api = "eth,net,web3"/api = "eth,net,web3,personal,txpool,debug"/g' "$HOME_DIR/config/app.toml"
  sed -i 's/minimum-gas-prices = ""/minimum-gas-prices = "0'"$GAS_DENOM"'"/g' "$HOME_DIR/config/app.toml"

  # Patch config.toml: allow all RPC connections
  sed -i 's|laddr = "tcp://127.0.0.1:26657"|laddr = "tcp://0.0.0.0:26657"|g' "$HOME_DIR/config/config.toml"
  sed -i 's|cors_allowed_origins = \[\]|cors_allowed_origins = ["*"]|g' "$HOME_DIR/config/config.toml"

  echo "==> MiniEVM node initialized."
else
  echo "==> MiniEVM node already initialized, skipping init."
fi

echo "==> Starting MiniEVM node..."
exec minievm start --home "$HOME_DIR"
