package config

import "flag"

type CommandConfig struct {
	Flags   *flag.FlagSet
	Dir     string
	EnvPath string
}

func (c *CommandConfig) GetCommandFlagArgs() []string {
	return c.Flags.Args()
}

func (c *CommandConfig) GetCommandDir() string {
	return c.Dir
}

func (c *CommandConfig) GetCommandEnvPath() string {
	return c.EnvPath
}
