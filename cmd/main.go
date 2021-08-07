package main

import (
	"github.com/milligan22963/cmra/cmd/subcmd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "AFMCamera",
		Short: "A camera application",
		Long:  `An application to interface with the pi camera and give a web view for it`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
			logrus.Info("starting up")
		},
	}

	rootCmd.AddCommand(subcmd.ServerCmd)
	subcmd.ServerCmd.Flags().IntP("port", "p", 9001, "port on which the server will listen")

	rootCmd.AddCommand(subcmd.VersionCmd)

	if err := rootCmd.Execute(); err != nil {
		logrus.Errorf("error executing cmd: %v", err)
	}
}
