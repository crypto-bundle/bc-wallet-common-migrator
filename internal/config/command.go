package config

import (
	"flag"
)

var (
	_ commandConfigService = (*CommandConfig)(nil)

	commandForbiddenFlags = map[string]struct{}{
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

func (c *CommandConfig) GetCommandEnvPath() *string {
	lookUp := c.Flags.Lookup(EnvFileParameterName)
	if lookUp == nil {
		return nil
	}

	value := lookUp.Value.String()

	return &value
}
