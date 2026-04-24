# Fhish CLI

The `fhish-cli` is the command-line orchestrator for the Fhish Private FHE Rollup Stack. It is a Go-based tool (originally forked from Initia's `weave-cli`) designed to streamline the provisioning, deployment, and management of privacy-preserving smart contract rollups on Initia.

## Architecture

The CLI wraps the entire stack and provides native management capabilities:

1. **MiniEVM Node Management**: Native commands to initialize, build, and run the Cosmos/EVM hybrid node locally without Docker if needed (`fhish node start`).
2. **Contract Deployer**: Uses `os/exec` to securely invoke Hardhat deployment scripts (`deploy.js`), parsing the outputs to dynamically capture the deployed Gateway and PrivateVoting contract addresses.
3. **FHE Key Generation**: The `fhish keys generate-fhe` command dynamically invokes the underlying `fhish-wasm` Rust binaries to generate the required Shortint Server/Client keypairs on the fly.
4. **Docker Orchestration**: The `fhish docker` subcommands provide wrappers around `docker-compose` to manage the fully containerized end-to-end verification stack.

## Subcommands
- `fhish create all`: Deploy the entire node and contract stack locally.
- `fhish docker up`: Launch the containerized verification stack.
- `fhish keys generate-fhe`: Generate FHE evaluation and client keys.
- `fhish node status`: Check the block height and RPC connectivity natively.
