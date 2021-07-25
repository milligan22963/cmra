package main

import (
	"cmd/main.go/cmd/subcmd"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cli struct {
	Port uint16 `help:"Default serving port."`
}

func main() {
	logrus.Info("starting up")
	var rootCmd = &cobra.Command{Use: "Set port to listen on"}
	rootCmd.AddCommand(subcmd.ServerCmd)
	rootCmd.Execute()
}
