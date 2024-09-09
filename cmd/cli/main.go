package main

import (
	"os"

	"github.com/lucax88x/wentsketchy/cmd/cli/console"
	"github.com/lucax88x/wentsketchy/cmd/cli/executor"
	"github.com/lucax88x/wentsketchy/internal/setup"
	"github.com/spf13/viper"
)

func cli(viper *viper.Viper, console *console.Console) setup.ProgramExecutor {
	return executor.NewCliExecutor(viper, console)
}

func main() {
	result := setup.Run(cli)

	if result == setup.NotOk {
		os.Exit(1)
	}
}
