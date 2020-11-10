package lib

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	vault "github.com/hashicorp/vault/api"
)

var vaultClient *vault.Client

// GetVaultClient gets a Vault Client authenticated by Kubernetes
func GetVaultClient() *vault.Client {
	if vaultClient != nil {
		return vaultClient
	}

	vaultClient, _ = vault.NewClient(nil)

	go func() {
		for {
			jwt, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
			if err != nil {
				log.Fatal(err, "Unable to open ServiceAccount Token")
			}

			resp, err := vaultClient.Logical().Write("auth/kubernetes/login", map[string]interface{}{
				"role": os.Getenv("VAULT_ROLE"),
				"jwt":  string(jwt),
			})
			if err != nil {
				log.Fatal(err, "Unable to refresh vault token")
			}

			vaultClient.SetToken(resp.Auth.ClientToken)
			log.Println("Refreshed vault token")

			ttl, _ := resp.TokenTTL()
			time.Sleep(ttl - (5 * time.Minute))
		}
	}()

	return vaultClient
}
