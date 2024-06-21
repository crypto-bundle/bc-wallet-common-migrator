package config

import (
	commonConfig "github.com/crypto-bundle/bc-wallet-common-lib-config/pkg/config"
	commonLogger "github.com/crypto-bundle/bc-wallet-common-lib-logger/pkg/logger"
	commonPostgres "github.com/crypto-bundle/bc-wallet-common-lib-postgres/pkg/postgres"
)

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
	// ----------------------------
	// Dependencies
	baseAppCfgSvc baseConfigService
	loggerCfgSvc  loggerCfgService
}
