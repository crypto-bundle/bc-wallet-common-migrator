package config

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	commonConfig "github.com/crypto-bundle/bc-wallet-common-lib-config/pkg/config"
	commonLogger "github.com/crypto-bundle/bc-wallet-common-lib-logger/pkg/logger"
	commonVault "github.com/crypto-bundle/bc-wallet-common-lib-vault/pkg/vault"
	commonVaultTokenClient "github.com/crypto-bundle/bc-wallet-common-lib-vault/pkg/vault/client/token"

	"github.com/joho/godotenv"
)

const (
	MigrationDirParameterName = "dir"
	EnvFileParameterName      = "envPath"
)

func PrepareLogger(ctx context.Context,
	baseCfgSrv baseConfigService,
) (*commonLogger.LoggerConfig, error) {
	cfgPreparerSrv := commonConfig.NewConfigManager()
	loggerCfg := &commonLogger.LoggerConfig{}

	err := cfgPreparerSrv.PrepareTo(loggerCfg).With(baseCfgSrv).Do(ctx)
	if err != nil {
		return nil, err
	}

	return loggerCfg, nil
}

func PrepareVault(ctx context.Context,
	baseCfgSrv baseConfigService,
	stdLogger *log.Logger,
) (*commonVault.Service, error) {
	cfgPreparerSrv := commonConfig.NewConfigManager()
	vaultCfg := &VaultWrappedConfig{
		BaseConfig: &commonVault.BaseConfig{},
		AuthConfig: &commonVaultTokenClient.AuthConfig{},
	}
	err := cfgPreparerSrv.PrepareTo(vaultCfg).With(baseCfgSrv).Do(ctx)
	if err != nil {
		return nil, err
	}

	vaultClientSrv, err := commonVaultTokenClient.NewClient(ctx, vaultCfg)
	if err != nil {
		return nil, err
	}

	vaultSrv, err := commonVault.NewService(stdLogger, vaultCfg, vaultClientSrv)
	if err != nil {
		return nil, err
	}

	_, err = vaultSrv.Login(ctx)
	if err != nil {
		return nil, err
	}

	return vaultSrv, nil
}

func PrepareCommand() (*CommandConfig, error) {
	err := CheckForbiddenFlags()
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

	// TODO: Add env path validation
	envPath := cmd.GetCommandEnvPath()
	if envPath != nil && *envPath != "" {
		loadErr := godotenv.Load(*envPath)
		if loadErr != nil {
			return nil, loadErr
		}

		return cmd, nil
	}

	return cmd, nil
}

func CheckForbiddenFlags() error {
	for _, arg := range os.Args[1:] {
		if len(arg) > 1 && arg[0] == '-' {
			flagName := strings.TrimLeft(arg[1:], "-")
			if _, ok := CommandForbiddenFlags[flagName]; ok {
				return fmt.Errorf("command contains forbidden flag: %s", flagName)
			}
		}
	}

	return nil
}

func PrepareAppCfg(ctx context.Context,
	wrappedBaseCfgSvc *BaseConfigWrapper,
	stdLogger *log.Logger,
) (*Config, *commonVault.Service, error) {
	baseCfg, loggerCfg, commandCfg := wrappedBaseCfgSvc.BaseConfig,
		wrappedBaseCfgSvc.LoggerConfig, wrappedBaseCfgSvc.CommandConfig

	vaultSecretSvc, err := PrepareVault(ctx, baseCfg, stdLogger)
	if err != nil {
		return nil, nil, err
	}

	err = vaultSecretSvc.LoadSecrets(ctx)
	if err != nil {
		return nil, nil, err
	}

	appCfgPreparerSrv := commonConfig.NewConfigManager()
	wrappedConfig := &Config{}
	err = appCfgPreparerSrv.PrepareTo(wrappedConfig).With(baseCfg,
		loggerCfg, vaultSecretSvc).Do(ctx)
	if err != nil {
		return nil, nil, err
	}

	wrappedConfig.baseAppCfgSvc = baseCfg
	wrappedConfig.loggerCfgSvc = loggerCfg
	wrappedConfig.CommandConfig = commandCfg

	return wrappedConfig, vaultSecretSvc, nil
}

func PrepareBaseConfig(ctx context.Context,
	applicationName string,
	releaseTag,
	commitID,
	shortCommitID,
	buildNumber,
	buildDateTS string,
) (*BaseConfigWrapper, error) {
	commandCfg, err := PrepareCommand()
	if err != nil {
		return nil, err
	}

	flagManagerSvc, err := commonConfig.NewLdFlagsManager(releaseTag,
		commitID, shortCommitID,
		buildNumber, buildDateTS)
	if err != nil {
		return nil, err
	}

	err = commonConfig.LoadLocalEnvIfDev()
	if err != nil {
		return nil, err
	}

	baseCfgPreparerSvc := commonConfig.NewConfigManager()
	baseCfg := commonConfig.NewBaseConfig(applicationName)
	err = baseCfgPreparerSvc.PrepareTo(baseCfg).With(flagManagerSvc).Do(ctx)
	if err != nil {
		return nil, err
	}

	loggerConfig, err := PrepareLogger(ctx, baseCfg)
	if err != nil {
		return nil, err
	}

	return &BaseConfigWrapper{
		BaseConfig:    baseCfg,
		LoggerConfig:  loggerConfig,
		CommandConfig: commandCfg,
	}, nil
}
