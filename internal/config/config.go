package config

import (
	commonConfig "gitlab.heronodes.io/bc-platform/bc-wallet-common-lib-config/pkg/config"
	commonLogger "gitlab.heronodes.io/bc-platform/bc-wallet-common-lib-logger/pkg/logger"
	commonPostgres "gitlab.heronodes.io/bc-platform/bc-wallet-common-lib-postgres/pkg/postgres"
	commonVault "gitlab.heronodes.io/bc-platform/bc-wallet-common-lib-vault/pkg/vault"
	commonVaultTokenClient "gitlab.heronodes.io/bc-platform/bc-wallet-common-lib-vault/pkg/vault/client/token"
)

type VaultWrappedConfig struct {
	*commonVault.BaseConfig
	*commonVaultTokenClient.AuthConfig
}

// Config for application
type Config struct {
	// -------------------
	// External common configs
	// -------------------
	*commonConfig.BaseConfig
	*commonLogger.LoggerConfig
	*commonPostgres.PostgresConfig
	*VaultWrappedConfig
	// -------------------
	// Internal configs
	// -------------------
	*CommandConfig
}
