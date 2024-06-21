package config

import (
	"time"

	commonConfig "github.com/crypto-bundle/bc-wallet-common-lib-config/pkg/config"
	commonLogger "github.com/crypto-bundle/bc-wallet-common-lib-logger/pkg/logger"
	commonVault "github.com/crypto-bundle/bc-wallet-common-lib-vault/pkg/vault"
	commonVaultTokenClient "github.com/crypto-bundle/bc-wallet-common-lib-vault/pkg/vault/client/token"
)

type baseConfigService interface {
	GetHostName() string
	GetEnvironmentName() string
	IsProd() bool
	IsStage() bool
	IsTest() bool
	IsDev() bool
	IsDebug() bool
	IsLocal() bool
	GetStageName() string
	GetApplicationPID() int
	GetApplicationName() string
	SetApplicationName(appName string)
	GetReleaseTag() string
	GetCommitID() string
	GetShortCommitID() string
	GetBuildNumber() uint64
	GetBuildDateTS() int64
	GetBuildDate() time.Time
}

type VaultWrappedConfig struct {
	*commonVault.BaseConfig
	*commonVaultTokenClient.AuthConfig
}

type BaseConfigWrapper struct {
	*commonConfig.BaseConfig
	*commonLogger.LoggerConfig
	*CommandConfig
}

type loggerCfgService interface {
	GetMinimalLogLevel() string
	IsStacktraceEnabled() bool
	GetSkipBuildInfo() bool
}
