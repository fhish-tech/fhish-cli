const fs = require('fs');
const { ethers } = require('ethers');

// Setup crypto for WASM
if (typeof globalThis.crypto === 'undefined') {
  globalThis.crypto = require('crypto').webcrypto;
}

// Load WASM module
const wasm = require('/packages/fhish-wasm/pkg-node/fhish_wasm.js');

const RPC_URL = process.env.RPC_URL || 'http://minievm:8545';
const GATEWAY_URL = process.env.GATEWAY_URL || 'http://gateway:8080';
const ADDRESSES_FILE = process.env.ADDRESSES_FILE || '/shared/addresses.json';
const PRIVATE_KEY = process.env.PRIVATE_KEY || '0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80';

const PrivateVotingV2_ABI = [
  "function vote(bytes32 handleA, bytes32 handleB, bytes memory, bytes memory) external",
  "function requestDecryptResults() external",
  "function getVoteCount() external view returns (uint32, uint32)",
  "function isDecrypted() external view returns (bool)",
  "function finalTallyA() external view returns (uint32)",
  "function finalTallyB() external view returns (uint32)"
];

const FhishGateway_ABI = [
  "function submitCiphertext(bytes calldata ciphertext) external returns (bytes32 handle)"
];

async function main() {
  console.log("🚀 Starting E2E FHE Verification...");

  // 1. Wait for addresses
  console.log(`Waiting for addresses at ${ADDRESSES_FILE}...`);
  while (!fs.existsSync(ADDRESSES_FILE)) {
    await new Promise(r => setTimeout(r, 2000));
  }
  const addrs = JSON.parse(fs.readFileSync(ADDRESSES_FILE, 'utf8'));
  const gatewayAddr = addrs.FhishGateway;
  const votingAddr = addrs.PrivateVotingV2;
  console.log(`✅ Loaded Contracts: Gateway=${gatewayAddr}, Voting=${votingAddr}`);

  // 2. Connect to RPC
  const provider = new ethers.JsonRpcProvider(RPC_URL);
  const wallet = new ethers.Wallet(PRIVATE_KEY, provider);
  const gateway = new ethers.Contract(gatewayAddr, FhishGateway_ABI, wallet);
  const voting = new ethers.Contract(votingAddr, PrivateVotingV2_ABI, wallet);

  // 3. Fetch FHE Public Key
  console.log(`Fetching FHE Public Key from ${GATEWAY_URL}/get-public-key...`);
  let pkRes;
  while (true) {
    try {
      const res = await fetch(`${GATEWAY_URL}/get-public-key`);
      if (res.ok) {
        pkRes = await res.json();
        break;
      }
    } catch (e) {
      // ignore
    }
    await new Promise(r => setTimeout(r, 2000));
  }
  
  if (!pkRes.publicKey) {
    throw new Error("Failed to get public key from gateway");
  }

  // Convert hex to bytes
  const pkBytes = Uint8Array.from(Buffer.from(pkRes.publicKey.replace('0x', ''), 'hex'));
  console.log(`✅ Fetched FHE Public Key (${pkBytes.length} bytes, type: ${pkRes.type})`);

  // 4. Encrypt Vote (YES on A, NO on B) using Pre-encrypted ciphertexts
  console.log("🔐 Loading pre-encrypted votes...");
  
  const ctABytes = fs.readFileSync(__dirname + '/ct_yes.bin');
  const ctBBytes = fs.readFileSync(__dirname + '/ct_no.bin');

  console.log(`✅ Loaded: A=1 (${ctABytes.length} bytes), B=0 (${ctBBytes.length} bytes)`);

  // 5. Submit Ciphertexts to Gateway HTTP API
  console.log("📤 Submitting ciphertexts to Gateway HTTP API...");
  
  async function submitCtHttp(ctBytes) {
    const hex = "0x" + Buffer.from(ctBytes).toString('hex');
    const res = await fetch(`${GATEWAY_URL}/ciphertext`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ ciphertext: hex }),
    });
    if (!res.ok) throw new Error("Failed to submit to gateway: " + await res.text());
    const data = await res.json();
    if (!data.handle.startsWith("0x")) data.handle = "0x" + data.handle;
    return data.handle;
  }

  const handleA = await submitCtHttp(ctABytes);
  console.log(`✅ Submitted Ciphertext A -> Handle: ${handleA}`);

  const handleB = await submitCtHttp(ctBBytes);
  console.log(`✅ Submitted Ciphertext B -> Handle: ${handleB}`);

  // 6. Cast Vote
  console.log("🗳️ Casting encrypted vote...");

  const txVote = await voting.vote(handleA, "0x0000000000000000000000000000000000000000000000000000000000000000", "0x", "0x");
  await txVote.wait();
  console.log("✅ Vote cast successfully");

  // 7. Request Decryption
  console.log("🔓 Requesting decryption of results...");
  const txDecrypt = await voting.requestDecryptResults();
  await txDecrypt.wait();
  console.log("✅ Decryption requested. Waiting for relayer and gateway...");

  // 8. Wait for decryption
  process.stdout.write("⏳ Polling for decrypted results.");
  let isDecrypted = false;
  for (let i = 0; i < 30; i++) {
    isDecrypted = await voting.isDecrypted();
    if (isDecrypted) break;
    process.stdout.write(".");
    await new Promise(r => setTimeout(r, 2000));
  }
  console.log("");

  if (!isDecrypted) {
    throw new Error("❌ Decryption timeout! Relayer/Gateway did not fulfill.");
  }

  // 9. Assert Results
  const tallyA = await voting.finalTallyA();
  const tallyB = await voting.finalTallyB();

  console.log(`\n🎉 Decryption fulfilled!`);
  console.log(`📊 Tally A: ${tallyA.toString()}`);
  console.log(`📊 Tally B: ${tallyB.toString()}`);

  if (tallyA.toString() === "1" && tallyB.toString() === "0") {
    console.log("\n✅ E2E VERIFICATION SUCCESSFUL!");
    process.exit(0);
  } else {
    console.error("\n❌ E2E VERIFICATION FAILED: Unexpected tally.");
    process.exit(1);
  }
}

main().catch(console.error);
