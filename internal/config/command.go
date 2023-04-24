package config

import (
	"flag"
)

var (
	CommandForbiddenFlags = map[string]struct{}{
		"port":     {},
		"host":     {},
		"user":     {},
		"password": {},
		"dbname":   {},
	}
)

type CommandConfig struct {
	Flags            *flag.FlagSet
	MigrationDirPath string
	EnvFilePath      string
}

func (c *CommandConfig) GetCommandFlagArgs() []string {
	return c.Flags.Args()
}

func (c *CommandConfig) GetCommandDir() string {
	return c.MigrationDirPath
}

func (c *CommandConfig) GetCommandEnvPath() string {
	return c.EnvFilePath
}
