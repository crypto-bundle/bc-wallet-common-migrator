package config

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	commonConfig "github.com/crypto-bundle/bc-wallet-common-lib-config/pkg/config"
	commonLogger "github.com/crypto-bundle/bc-wallet-common-lib-logger/pkg/logger"
	commonVault "github.com/crypto-bundle/bc-wallet-common-lib-vault/pkg/vault"
	commonVaultTokenClient "github.com/crypto-bundle/bc-wallet-common-lib-vault/pkg/vault/client/token"
)

const (
	MigrationDirParameterName = "dir"
	EnvFileParameterName      = "envPath"
)

func PrepareLogger(ctx context.Context,
	baseCfgSvc baseConfigService,
	errFmtSvc errorFormatterService,
) (*commonLogger.LoggerConfig, error) {
	cfgPreparerSrv := commonConfig.NewConfigManager(errFmtSvc)
	loggerCfg := &commonLogger.LoggerConfig{}

	err := cfgPreparerSrv.PrepareTo(loggerCfg).With(baseCfgSvc).Do(ctx)
	if err != nil {
		return nil, errFmtSvc.ErrorNoWrap(err)
	}

	return loggerCfg, nil
}

func PrepareVault(ctx context.Context,
	errFmtSvc errorFormatterService,
	baseCfgSvc baseConfigService,
	loggerBuilderSvc loggerService,
) (*commonVault.Service, error) {
	cfgPreparerSrv := commonConfig.NewConfigManager(errFmtSvc)
	vaultCfg := &VaultWrappedConfig{
		BaseConfig: &commonVault.BaseConfig{},
		AuthConfig: &commonVaultTokenClient.AuthConfig{},
	}

	err := cfgPreparerSrv.PrepareTo(vaultCfg).With(baseCfgSvc, errFmtSvc).Do(ctx)
	if err != nil {
		return nil, errFmtSvc.ErrorNoWrap(err)
	}

	vaultClientSvc, err := commonVaultTokenClient.NewClient(ctx, errFmtSvc, vaultCfg)
	if err != nil {
		return nil, errFmtSvc.ErrorNoWrap(err)
	}

	vaultSvc, err := commonVault.NewService(loggerBuilderSvc, errFmtSvc, vaultCfg, vaultClientSvc)
	if err != nil {
		return nil, errFmtSvc.ErrorNoWrap(err)
	}

	_, err = vaultSvc.Login(ctx)
	if err != nil {
		return nil, errFmtSvc.ErrorNoWrap(err)
	}

	return vaultSvc, nil
}

func PrepareCommand() (*CommandConfig, error) {
	err := checkForbiddenFlags()
	if err != nil {
		return nil, err
	}

	cmd := &CommandConfig{
		Flags: flag.NewFlagSet("bc-wallet-common-migrator", flag.ContinueOnError),
	}

	cmd.Flags.SetOutput(io.Discard)

	dirPath := ""
	cmd.Flags.StringVar(&dirPath, MigrationDirParameterName, "./migrations", "directory with migration files")

	envFilePath := ""
	cmd.Flags.StringVar(&envFilePath, EnvFileParameterName, "", "environment file with migration settings")

	args := os.Args[1:]

	err = cmd.Flags.Parse(args)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

func checkForbiddenFlags() error {
	for _, arg := range os.Args[1:] {
		if len(arg) == 0 || arg[0] != '-' {
			continue
		}

		flagName := strings.TrimLeft(arg[1:], "-")
		if _, ok := commandForbiddenFlags[flagName]; ok {
			return fmt.Errorf("command contains forbidden flag: %s", flagName)
		}
	}

	return nil
}

func PrepareAppCfg(ctx context.Context,
	loggerBuilderSvc loggerService,
	errFmtSvc errorFormatterService,
	wrappedBaseCfgSvc *BaseConfigWrapper,
) (*Config, *commonVault.Service, error) {
	baseCfg, loggerCfg, commandCfg := wrappedBaseCfgSvc.BaseConfig,
		wrappedBaseCfgSvc.LoggerConfig, wrappedBaseCfgSvc.CommandConfig

	vaultSecretSvc, err := PrepareVault(ctx, errFmtSvc,
		baseCfg, loggerBuilderSvc)
	if err != nil {
		return nil, nil, err
	}

	err = vaultSecretSvc.LoadSecrets(ctx)
	if err != nil {
		return nil, nil, errFmtSvc.ErrorNoWrap(err)
	}

	appCfgPreparerSrv := commonConfig.NewConfigManager(errFmtSvc)
	appConfig := &Config{}

	err = appCfgPreparerSrv.PrepareTo(appConfig).With(baseCfg,
		loggerCfg, vaultSecretSvc, commandCfg).Do(ctx)
	if err != nil {
		return nil, nil, errFmtSvc.ErrorNoWrap(err)
	}

	return appConfig, vaultSecretSvc, nil
}

func PrepareBaseConfig(ctx context.Context,
	errFmtSvc errorFormatterService,
	applicationName string,
	releaseTag,
	commitID,
	shortCommitID,
	buildNumber,
	buildDateTS string,
) (*BaseConfigWrapper, error) {
	commandCfg, err := PrepareCommand()
	if err != nil {
		return nil, errFmtSvc.ErrorNoWrap(err)
	}

	flagManagerSvc, err := commonConfig.NewLdFlagsManager(errFmtSvc, releaseTag,
		commitID, shortCommitID,
		buildNumber, buildDateTS)
	if err != nil {
		return nil, errFmtSvc.ErrorNoWrap(err)
	}

	envPath := commandCfg.GetCommandEnvPath()
	if envPath != nil && *envPath != "" {
		loadErr := commonConfig.LoadEnvFromFile(*envPath)
		if loadErr != nil {
			return nil, errFmtSvc.ErrorNoWrap(loadErr)
		}
	}

	baseCfgPreparerSvc := commonConfig.NewConfigManager(errFmtSvc)
	baseCfg := commonConfig.NewBaseConfig(applicationName)

	err = baseCfgPreparerSvc.PrepareTo(baseCfg).With(flagManagerSvc, errFmtSvc).Do(ctx)
	if err != nil {
		return nil, errFmtSvc.ErrorNoWrap(err)
	}

	loggerConfig, err := PrepareLogger(ctx, baseCfg, errFmtSvc)
	if err != nil {
		return nil, errFmtSvc.ErrorNoWrap(err)
	}

	return &BaseConfigWrapper{
		BaseConfig:    baseCfg,
		LoggerConfig:  loggerConfig,
		CommandConfig: commandCfg,
	}, nil
}
