package config

import (
	"log/slog"
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

type commandConfigService interface {
	GetCommandFlagArgs() []string
	GetCommandDir() string
	GetCommandEnvPath() *string
}

type loggerService interface {
	NewSlogLoggerEntry(fields ...any) *slog.Logger
	NewSlogNamedLoggerEntry(named string, fields ...any) *slog.Logger
	NewSlogLoggerEntryWithFields(fields ...slog.Attr) *slog.Logger
}

//nolint:interfacebloat //it's ok here, we need it we must use it as one big interface
type errorFormatterService interface {
	ErrorWithCode(err error, code int) error
	ErrWithCode(err error, code int) error
	ErrorGetCode(err error) int
	ErrGetCode(err error) int
	// ErrorNoWrap function for pseudo-wrap error, must be used in case of linter warnings...
	ErrorNoWrap(err error) error
	// ErrNoWrap same with ErrorNoWrap function, just alias for ErrorNoWrap, just short function name...
	ErrNoWrap(err error) error
	ErrorOnly(err error, details ...string) error
	Error(err error, details ...string) error
	Errorf(err error, format string, args ...interface{}) error
	NewError(details ...string) error
	NewErrorf(format string, args ...interface{}) error
}
