package main

import (
	"github.com/milligan22963/cmra/cmd/subcmd"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "AFMCamera",
		Short: "A camera application",
		Long:  `An application to interface with the pi camera and give a web view for it`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
			print("starting up\n")
		},
	}

	rootCmd.AddCommand(subcmd.ServerCmd)
	subcmd.ServerCmd.Flags().String("config", "settings.yaml", "configuration file")

	rootCmd.AddCommand(subcmd.VersionCmd)

	if err := rootCmd.Execute(); err != nil {
		println("failed to initialize: ", err.Error())
	}
}
