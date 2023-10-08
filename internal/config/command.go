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
	Flags *flag.FlagSet
}

func (c *CommandConfig) GetCommandFlagArgs() []string {
	return c.Flags.Args()
}

func (c *CommandConfig) GetCommandDir() string {
	return c.Flags.Lookup(MigrationDirParameterName).Value.String()
}

func (c *CommandConfig) GetCommandEnvPath() string {
	return c.Flags.Lookup(EnvFileParameterName).Value.String()
}
