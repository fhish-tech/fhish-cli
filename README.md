# Fhish CLI

The `fhish-cli` is the command-line orchestrator for the Fhish Private FHE Rollup Stack. It is a Go-based tool (originally forked from Initia's `weave-cli`) designed to streamline the provisioning, deployment, and management of privacy-preserving smart contract rollups on Initia.

## Architecture

The CLI wraps the entire stack and provides native management capabilities:

1. **MiniEVM Node Management**: Native commands to initialize, build, and run the Cosmos/EVM hybrid node locally without Docker if needed (`fhish node start`).
2. **Contract Deployer**: Uses `os/exec` to securely invoke Hardhat deployment scripts (`deploy.js`), parsing the outputs to dynamically capture the deployed Gateway and PrivateVoting contract addresses.
3. **FHE Key Generation**: The `fhish keys generate-fhe` command dynamically invokes the underlying `fhish-wasm` Rust binaries to generate the required Shortint Server/Client keypairs on the fly.
4. **Docker Orchestration**: The `fhish docker` subcommands provide wrappers around `docker-compose`## Installation

The Fhish CLI is distributed as a native binary for maximum performance.

### macOS (Silicon/M1/M2/M3)
```bash
curl -L -o fhish https://github.com/fhish-tech/fhish-cli/releases/download/v0.1.7/fhish-darwin-arm64
chmod +x fhish
sudo mv fhish /usr/local/bin/
```

### macOS (Intel)
```bash
curl -L -o fhish https://github.com/fhish-tech/fhish-cli/releases/download/v0.1.7/fhish-darwin-amd64
chmod +x fhish
sudo mv fhish /usr/local/bin/
```

### Linux (AMD64)
```bash
curl -L -o fhish https://github.com/fhish-tech/fhish-cli/releases/download/v0.1.7/fhish-linux-amd64
chmod +x fhish
sudo mv fhish /usr/local/bin/
```

## Quick Start

1. **Check Installation**:
   ```bash
   fhish help
   ```

2. **Initialize a Node**:
   ```bash
   fhish node init --chain-id fhish-1
   ```

3. **Start the Stack**:
   ```bash
   # Launch MiniEVM, Gateway, and Relayer via Docker
   fhish docker up
   ```

## Commands

- `fhish create`: Wizard to create new rollups or services.
- `fhish node`: Start/Stop/Status of the local MiniEVM node.
- `fhish gateway`: Manage the FHE decryption gateway.
- `fhish relayer`: Manage the FHE relayer service.
- `fhish keys`: Generate FHE evaluation and client keys.
- `fhish version`: Show current version (v0.1.7).

## Documentation

For full documentation, visit [docs.fhish.tech](https://docs.fhish.tech).
