package main

import (
	"context"
	"log"

	"gitlab.heronodes.io/bc-platform/bc-wallet-common-migrator/internal/config"

	commonLogger "gitlab.heronodes.io/bc-platform/bc-wallet-common-lib-logger/pkg/logger"
	commonPostgres "gitlab.heronodes.io/bc-platform/bc-wallet-common-lib-postgres/pkg/postgres"

	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

// DO NOT EDIT THESE VARIABLES DIRECTLY. These are build-time constants
// DO NOT USE THESE VARIABLES IN APPLICATION CODE. USE commonConfig.NewLdFlagsManager SERVICE-COMPONENT INSTEAD OF IT
var (
	// Version - version time.RFC3339.
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE
	Version = "DEVELOPMENT.VERSION"

	// ReleaseTag - release tag in TAG.%Y-%m-%dT%H-%M-%SZ.
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE
	ReleaseTag = "DEVELOPMENT.RELEASE_TAG"

	// CommitID - latest commit id.
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE
	CommitID = "DEVELOPMENT.COMMIT_HASH"

	// ShortCommitID - first 12 characters from CommitID.
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE
	ShortCommitID = "DEVELOPMENT.SHORT_COMMIT_HASH"

	// BuildNumber - ci/cd build number for BuildNumber
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE
	BuildNumber uint64 = 0

	// BuildDateTS - ci/cd build date in time stamp
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE
	BuildDateTS uint64 = 0
)

const ApplicationName = "bc-wallet-common-migrator"

func main() {
	var err error
	ctx, cancelCtxFunc := context.WithCancel(context.Background())

	appCfg, _, err := config.Prepare(ctx, Version, ReleaseTag,
		CommitID, ShortCommitID,
		BuildNumber, BuildDateTS, ApplicationName)
	if err != nil {
		log.Fatal("unable prepare application config", err)
	}

	loggerSrv, err := commonLogger.NewService(appCfg)
	if err != nil {
		log.Fatal("unable create logger service", err)
	}
	loggerEntry := loggerSrv.NewLoggerEntry("main")

	pgConn := commonPostgres.NewConnection(context.Background(), appCfg, loggerEntry)
	_, err = pgConn.Connect()
	if err != nil {
		loggerEntry.Fatal("unable to connect to to database", zap.Error(err))
	}

	goose.SetLogger(zap.NewStdLog(loggerEntry.Named("goose.service")))

	commandArgs := appCfg.GetCommandFlagArgs()
	err = goose.RunWithOptions(commandArgs[0], pgConn.Dbx.DB, appCfg.GetCommandDir(), commandArgs[1:])
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
