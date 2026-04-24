package relayer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fhish/fhish-cli/contracts"
)

func WriteRelayerConfig(homeDir string, addrs *contracts.ContractAddresses, privateKey string, gatewaySecret string, rpcURL string) error {
	envContent := fmt.Sprintf(`PRIVATE_KEY=%s
RPC_URL=%s
GATEWAY_URL=http://localhost:3000
VOTING_ADDRESS=%s
FHISH_RELAYER_SECRET=%s
HEALTH_PORT=3001
`, privateKey, rpcURL, addrs.PrivateVoting, gatewaySecret)

	return os.WriteFile(filepath.Join(homeDir, "relayer.env"), []byte(envContent), 0644)
}
