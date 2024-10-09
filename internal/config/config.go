package config

import (
	commonConfig "github.com/crypto-bundle/bc-wallet-common-lib-config/pkg/config"
	commonLogger "github.com/crypto-bundle/bc-wallet-common-lib-logger/pkg/logger"
	commonPostgres "github.com/crypto-bundle/bc-wallet-common-lib-postgres/pkg/postgres"
)

// Config for application...
type Config struct {
	// -------------------
	// External common configs
	// -------------------
	*commonConfig.BaseConfig
	*commonLogger.LoggerConfig
	*commonPostgres.PostgresConfig
	*VaultWrappedConfig
	// ----------------------------
	// Dependencies
	baseAppCfgSvc baseConfigService
	loggerCfgSvc  loggerCfgService
	commandCfgSvc commandConfigService
}

func (c *Config) GetCommandFlagArgs() []string {
	return c.commandCfgSvc.GetCommandFlagArgs()
}

func (c *Config) GetCommandDir() string {
	return c.commandCfgSvc.GetCommandDir()
}

func (c *Config) GetCommandEnvPath() *string {
	return c.commandCfgSvc.GetCommandEnvPath()
}

func (c *Config) PrepareWith(dependentCfgList ...interface{}) error {
	for _, cfgSrv := range dependentCfgList {
		switch castedDep := cfgSrv.(type) {
		case baseConfigService:
			c.baseAppCfgSvc = castedDep
		case loggerCfgService:
			c.loggerCfgSvc = castedDep
		case commandConfigService:
			c.commandCfgSvc = castedDep
		default:
			continue
		}
	}

	return nil
}
