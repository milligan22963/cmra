package main

import (
	"github.com/milligan22963/cmra/cmd/subcmd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	logrus.Info("starting up")

	rootCmd := &cobra.Command{
		Use:   "AFMCamera",
		Short: "A camera application",
		Long:  `An application to interface with the pi camera and give a web view for it`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}

	rootCmd.AddCommand(subcmd.ServerCmd)

	if err := rootCmd.Execute(); err != nil {
		logrus.Errorf("error executing cmd: %v", err)
	}
}
