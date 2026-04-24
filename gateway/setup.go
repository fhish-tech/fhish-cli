package gateway

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fhish/fhish-cli/contracts"
	"github.com/fhish/fhish-cli/utils"
)

func WriteGatewayConfig(homeDir string, addrs *contracts.ContractAddresses, relayerSecret string) error {
	envContent := fmt.Sprintf(`PORT=3000
BASE_URL=http://localhost:3000
FHISH_RELAYER_SECRET=%s
KEYS_DIR=%s
`, relayerSecret, filepath.Join(homeDir, "keys"))

	return os.WriteFile(filepath.Join(homeDir, "gateway.env"), []byte(envContent), 0644)
}

func GenerateKeys(gatewaySrc string, targetKeysDir string) error {
	utils.PrintInfo("Generating FHE keys...")
	_ = os.MkdirAll(targetKeysDir, 0755)
	
	// Use the keygen script in fhish-gateway
	err := utils.RunCommandWithOutput(gatewaySrc, "npm", "run", "keygen")
	if err != nil {
		return err
	}

	// Copy keys to target directory
	// (Implementation depends on where keygen outputs them)
	return nil
}
