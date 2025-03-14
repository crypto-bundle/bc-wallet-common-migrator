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

package main

import (
	"context"
	"log"

	"github.com/crypto-bundle/bc-wallet-common-migrator/internal/config"

	commonErrors "github.com/crypto-bundle/bc-wallet-common-lib-errors/pkg/errformatter"
	commonLogger "github.com/crypto-bundle/bc-wallet-common-lib-logger/pkg/logger"
	commonPostgres "github.com/crypto-bundle/bc-wallet-common-lib-postgres/pkg/postgres"

	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

// DO NOT EDIT THESE VARIABLES DIRECTLY. These are build-time constants
// DO NOT USE THESE VARIABLES IN APPLICATION CODE. USE commonConfig.NewLdFlagsManager SERVICE-COMPONENT INSTEAD OF IT...
var (
	// ReleaseTag - release tag in TAG.SHORT_COMMIT_ID.BUILD_NUMBER.
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE...
	ReleaseTag = "v0.0.0-00000000-100500"

	// CommitID - latest commit id.
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE...
	CommitID = "0000000000000000000000000000000000000000"

	// ShortCommitID - first 12 characters from CommitID.
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE...
	ShortCommitID = "00000000"

	// BuildNumber - ci/cd build number for BuildNumber
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE...
	BuildNumber string = "100500"

	// BuildDateTS - ci/cd build date in time stamp
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE...
	BuildDateTS string = "1713280105"
)

const ApplicationName = "bc-wallet-common-migrator"

func main() {
	var err error

	ctx, cancelCtxFunc := context.WithCancel(context.Background())

	cfgErrFmtSvc := commonErrors.NewScopedErrorFormatter("config")

	wrappedBaseCfg, err := config.PrepareBaseConfig(ctx, cfgErrFmtSvc,
		ApplicationName, ReleaseTag,
		CommitID, ShortCommitID,
		BuildNumber, BuildDateTS)
	if err != nil {
		log.Fatal(err.Error(), err)
	}

	loggerSvc, err := commonLogger.NewService(wrappedBaseCfg,
		commonErrors.NewScopedErrorFormatter("logger"))
	if err != nil {
		log.Fatal(err.Error(), err)
	}

	loggerEntry := loggerSvc.NewZapNamedLoggerEntry("migrator")

	appCfg, _, err := config.PrepareAppCfg(ctx, loggerSvc, cfgErrFmtSvc, wrappedBaseCfg)
	if err != nil {
		log.Fatal(err.Error(), err)
	}

	pgConn := commonPostgres.NewConnection(loggerSvc,
		commonErrors.NewScopedErrorFormatter("postgres"),
		appCfg)

	_, err = pgConn.Connect()
	if err != nil {
		loggerEntry.Fatal("unable to connect to database", zap.Error(err))
	}

	goose.SetLogger(loggerSvc.NewStdNamedLoggerEntry("goose"))

	commandArgs := appCfg.GetCommandFlagArgs()

	err = goose.RunWithOptionsContext(ctx, commandArgs[0],
		pgConn.Dbx.DB, appCfg.GetCommandDir(), commandArgs[1:])
	if err != nil {
		loggerEntry.Fatal("unable to run goose migration", zap.Error(err))
	}

	cancelCtxFunc()

	syncErr := loggerEntry.Sync()
	if syncErr != nil {
		log.Print(syncErr.Error(), syncErr)
	}

	log.Print("stopped")
}
