# Fhish CLI Integration

This document outlines the integration between the `fhish` CLI and the rest of the FHE stack.

## Architecture

The `fhish` CLI acts as a central orchestrator for the private FHE rollup stack:

1.  **MiniEVM Rollup**: Based on `initia-labs/minievm`. The CLI handles cloning, building, and genesis patching (adding FHE precompiles).
2.  **FHE Contracts**: Solidity contracts from `packages/fhish-contracts-v2`. Deployed using `forge` (Foundry).
3.  **FHE Gateway**: Decryption service from `fhish-gateway`. The CLI manages key generation and service lifecycle.
4.  **FHE Relayer**: Event monitor and decryption coordinator from `packages/fhish-relayer-v2`.

## Configuration Mapping

### FHE Gateway Env Vars
| Env Var | Description | Source in CLI |
|---------|-------------|---------------|
| `PORT` | Service port | Hardcoded to 3000 |
| `FHISH_RELAYER_SECRET` | Secret for relayer auth | Generated/Configured in `create` wizard |
| `KEYS_DIR` | Path to FHE keys | `~/.fhish/minievm/<chain-id>/keys` |

### FHE Relayer Env Vars
| Env Var | Description | Source in CLI |
|---------|-------------|---------------|
| `PRIVATE_KEY` | Relayer operator key | Configured in `create` wizard |
| `RPC_URL` | MiniEVM RPC URL | `http://localhost:8545` |
| `GATEWAY_URL` | Gateway service URL | `http://localhost:3000` |
| `VOTING_ADDRESS` | Main contract address | From `contracts.json` after deployment |
| `FHISH_RELAYER_SECRET` | Auth secret for gateway | Same as Gateway secret |

## Prerequisites

- **Go**: v1.23+
- **Node.js & npm**: For Gateway and Relayer.
- **Foundry (forge/cast)**: For contract deployment.
- **Git**: For cloning repositories.

## Service Management

All services are managed using PID files and logs:
- PID files: `~/.fhish/minievm/<chain-id>/run/*.pid`
- Logs: `~/.fhish/minievm/<chain-id>/logs/*.log`

Commands like `fhish node log` or `fhish gateway status` use these files to monitor the processes.
