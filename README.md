# Fhish CLI

The `fhish-cli` is the command-line orchestrator for the Fhish Private FHE Rollup Stack. It is a Go-based tool (originally forked from Initia's `weave-cli`) designed to streamline the provisioning, deployment, and management of privacy-preserving smart contract rollups on Initia.

## Installation (v0.1.8)

The Fhish CLI is distributed as a native binary for maximum performance.

### macOS (Silicon/M1/M2/M3)
```bash
curl -L -o fhish https://github.com/fhish-tech/fhish-cli/releases/download/v0.1.8/fhish-darwin-arm64
chmod +x fhish
sudo mv fhish /usr/local/bin/
```

### macOS (Intel)
```bash
curl -L -o fhish https://github.com/fhish-tech/fhish-cli/releases/download/v0.1.8/fhish-darwin-amd64
chmod +x fhish
sudo mv fhish /usr/local/bin/
```

### Linux (AMD64)
```bash
curl -L -o fhish https://github.com/fhish-tech/fhish-cli/releases/download/v0.1.8/fhish-linux-amd64
chmod +x fhish
sudo mv fhish /usr/local/bin/
```

## Quick Start

1. **Check Installation**:
   ```bash
   fhish version
   ```

2. **Launch the Setup Wizard**:
   The CLI provides an interactive TUI to set up your entire FHE stack. It handles EVM Chain ID calculation and FHE key generation automatically.
   ```bash
   fhish create all
   ```

3. **Start the Fhish Stack**:
   If you have an existing configuration, you can start the services via Docker:
   ```bash
   fhish docker up
   ```

4. **Verify the Environment**:
   Run the FHE end-to-end smoke test:
   ```bash
   fhish docker verify
   ```

## Key Features

- **Interactive TUI**: Bubbletea-powered wizard for easy setup.
- **Auto EVM Derivation**: Calculates deterministic EVM Chain IDs based on rollup names.
- **FHE Key Generation**: Integrated Rust-based key generator for SHORTINT FHE parameters.
- **Smart Docker Orchestration**: Automatically handles local source code detection and fallbacks.

## Commands

- `fhish create all`: Interactive wizard for full stack setup.
- `fhish docker`: Orchestrate the full stack (up, down, logs, verify).
- `fhish node`: Manage the local MiniEVM node (init, start, stop).
- `fhish keys`: Generate FHE evaluation and client keys.
- `fhish version`: Show current version (v0.1.8).

## Documentation

For full documentation, visit [docs.fhish.tech](https://docs.fhish.tech).
