package config

import (
	"context"
	"flag"
	"io"
	"os"

	commonConfig "github.com/crypto-bundle/bc-wallet-common-lib-config/pkg/config"
	commonVault "github.com/crypto-bundle/bc-wallet-common-lib-vault/pkg/vault"
	commonVaultTokenClient "github.com/crypto-bundle/bc-wallet-common-lib-vault/pkg/vault/client/token"

	"github.com/joho/godotenv"
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

func PrepareCommand(baseCfgSrv baseConfigService) (*Command, error) {
	cmd := &Command{
		Flags: flag.NewFlagSet("bc-wallet-common-migrator", flag.ContinueOnError),
	}

	cmd.Flags.SetOutput(io.Discard)

	cmd.Dir = *cmd.Flags.String("dir", "./migrations", "directory with migration files")
	cmd.EnvPath = *cmd.Flags.String("envPath", ".env", "environment file with migration settings")

	err := cmd.Flags.Parse(os.Args[1:])
	if err != nil {
		return nil, err
	}

	if baseCfgSrv.IsDev() {
		loadErr := godotenv.Load(cmd.EnvPath)
		if loadErr != nil {
			return nil, err
		}
	}

	return cmd, nil
}

func Prepare(ctx context.Context,
	version,
	releaseTag,
	commitID,
	shortCommitID string,
	buildNumber,
	buildDateTS uint64,
	applicationName string,
) (*Config, *commonVault.Service, error) {
	baseCfgSrv, err := PrepareBaseConfig(ctx, version, releaseTag,
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
	wrappedConfig.Command = command

	return wrappedConfig, vaultSecretSrv, nil
}

func PrepareBaseConfig(ctx context.Context,
	version,
	releaseTag,
	commitID,
	shortCommitID string,
	buildNumber,
	buildDateTS uint64,
	applicationName string,
) (*commonConfig.BaseConfig, error) {
	flagManagerSrv := commonConfig.NewLdFlagsManager(version, releaseTag,
		commitID, shortCommitID,
		buildNumber, buildDateTS)

	baseCfgPreparerSrv := commonConfig.NewConfigManager()
	baseCfg := commonConfig.NewBaseConfig(applicationName)
	err := baseCfgPreparerSrv.PrepareTo(baseCfg).With(flagManagerSrv).Do(ctx)
	if err != nil {
		return nil, err
	}

	return baseCfg, nil
}
