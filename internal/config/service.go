/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2025 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

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
		return nil, nil, errFmtSvc.ErrorNoWrap(err)
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
