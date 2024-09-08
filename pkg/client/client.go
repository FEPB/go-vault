package auth

import (
	"context"
	"time"

	"github.com/hashicorp/vault-client-go"
	"go.fepb.org.br/logger/pkg/logger"
)

type VaultClient struct {
	client *vault.Client
	path   string
}

func NewVaultClient(address, token, path string) VaultClient {
	client, err := vault.New(vault.WithAddress(address), vault.WithRequestTimeout(30*time.Second))
	if err != nil {
		logger.Fatal("Unable to initialize Vault client", err)
	}

	if err = client.SetToken(token); err != nil {
		logger.Fatal("Error setting Vault client token", err)
	}

	return VaultClient{client: client, path: path}
}

func (v VaultClient) GetSecret(ctx context.Context, tenantID string) (map[string]interface{}, error) {
	secret, err := v.client.Secrets.KvV2Read(ctx, v.path, vault.WithMountPath(tenantID))
	if err != nil {
		logger.Error("Unable to read secret", err)
		return nil, err
	}

	return secret.Data.Data, nil
}
