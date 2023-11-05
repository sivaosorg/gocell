package main

import (
	"os"

	"github.com/sivaosorg/gocell/internal/core"
	"github.com/sivaosorg/govm/cmd"
	"github.com/sivaosorg/govm/logger"
)

func main() {
	command := cmd.NewCommandManager()
	command.AddCommand(&core.CoreCommand{})
	err := command.Execute(os.Args)
	if err != nil {
		logger.Errorf("Internal Server", err, err.Error())
		panic(err)
	}
}
