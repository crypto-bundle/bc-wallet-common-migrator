package config

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	commonConfig "github.com/crypto-bundle/bc-wallet-common-lib-config/pkg/config"
	commonVault "github.com/crypto-bundle/bc-wallet-common-lib-vault/pkg/vault"
	commonVaultTokenClient "github.com/crypto-bundle/bc-wallet-common-lib-vault/pkg/vault/client/token"

	"github.com/joho/godotenv"
)

const (
	MigrationDirParameterName = "dir"
	EnvFileParameterName      = "envPath"
)

func PrepareVault(ctx context.Context, baseCfgSrv baseConfigService) (*commonVault.Service, error) {
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

	vaultSrv, err := commonVault.NewService(ctx, vaultCfg, vaultClientSrv)
	if err != nil {
		return nil, err
	}

	_, err = vaultSrv.Login(ctx)
	if err != nil {
		return nil, err
	}

	return vaultSrv, nil
}

func PrepareCommand(baseCfgSrv baseConfigService) (*CommandConfig, error) {
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

func Prepare(ctx context.Context,
	releaseTag,
	commitID,
	shortCommitID,
	buildNumber,
	buildDateTS,
	applicationName string,
) (*Config, *commonVault.Service, error) {
	baseCfgSrv, err := PrepareBaseConfig(ctx, releaseTag,
		commitID, shortCommitID,
		buildNumber, buildDateTS, applicationName)
	if err != nil {
		return nil, nil, err
	}

	command, err := PrepareCommand(baseCfgSrv)
	if err != nil {
		return nil, nil, err
	}

	vaultSecretSrv, err := PrepareVault(ctx, baseCfgSrv)
	if err != nil {
		return nil, nil, err
	}

	err = vaultSecretSrv.LoadSecrets(ctx)
	if err != nil {
		return nil, nil, err
	}

	appCfgPreparerSrv := commonConfig.NewConfigManager()
	wrappedConfig := &Config{}
	err = appCfgPreparerSrv.PrepareTo(wrappedConfig).With(baseCfgSrv, vaultSecretSrv).Do(ctx)
	if err != nil {
		return nil, nil, err
	}

	wrappedConfig.BaseConfig = baseCfgSrv
	wrappedConfig.CommandConfig = command

	return wrappedConfig, vaultSecretSrv, nil
}

func PrepareBaseConfig(ctx context.Context,
	releaseTag,
	commitID,
	shortCommitID,
	buildNumber,
	buildDateTS,
	applicationName string,
) (*commonConfig.BaseConfig, error) {
	flagManagerSrv, err := commonConfig.NewLdFlagsManager(releaseTag,
		commitID, shortCommitID,
		buildNumber, buildDateTS)
	if err != nil {
		return nil, err
	}

	baseCfgPreparerSrv := commonConfig.NewConfigManager()
	baseCfg := commonConfig.NewBaseConfig(applicationName)
	err = baseCfgPreparerSrv.PrepareTo(baseCfg).With(flagManagerSrv).Do(ctx)
	if err != nil {
		return nil, err
	}

	return baseCfg, nil
}
