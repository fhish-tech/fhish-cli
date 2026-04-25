import { ethers } from "ethers"

const RPC_URL = process.env.NEXT_PUBLIC_RPC_URL || "http://localhost:8545"

export const getProvider = () => {
  return new ethers.JsonRpcProvider(RPC_URL)
}

export const VOTING_ABI = [
  "event VoteCast(address indexed voter, bytes32 handleA, bytes32 handleB, uint256 voteId)",
  "function getVoteCount() external view returns (uint256, uint256)",
]
