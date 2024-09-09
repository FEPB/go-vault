package auth

import (
	"context"
	"time"

	"github.com/hashicorp/vault-client-go"
	"go.fepb.org.br/logger/pkg/logger"
	"go.uber.org/zap"
)

type VaultClient struct {
	client *vault.Client
}

func NewVaultClient(address, token string) VaultClient {
	client, err := vault.New(vault.WithAddress(address), vault.WithRequestTimeout(30*time.Second))
	if err != nil {
		logger.Fatal("Unable to initialize Vault client", zap.Error(err))
	}

	if err = client.SetToken(token); err != nil {
		logger.Fatal("Error setting Vault client token", zap.Error(err))
	}

	return VaultClient{client: client}
}

func (v VaultClient) GetSecret(ctx context.Context, secretName, tenantID string) (map[string]interface{}, error) {
	secret, err := v.client.Secrets.KvV2Read(ctx, secretName, vault.WithMountPath(tenantID))
	if err != nil {
		logger.Error("Unable to read secret", zap.Error(err))
		return nil, err
	}

	return secret.Data.Data, nil
}
