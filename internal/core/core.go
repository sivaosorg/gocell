package core

import (
	"fmt"

	"github.com/sivaosorg/govm/configx"
)

// Conf global
// Get & Set key-value
// Based on conf, to create new cluster / instance
var Conf *configx.KeysConfig

type CoreCommand struct {
}

func (c *CoreCommand) Name() string {
	return "start_server"
}

func (c *CoreCommand) Description() string {
	return "Start Server HTTP/1.1"
}

func (c *CoreCommand) Execute(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("CoreCommand Args is required")
	}
	keys, err := configx.ReadConfig[configx.KeysConfig](args[0])
	if err != nil {
		return err
	}
	Conf = keys
	return nil
}
