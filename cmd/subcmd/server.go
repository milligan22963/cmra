// Package subcmd is for all subcmds in the cmd tree
package subcmd

import (
	"github.com/milligan22963/cmra/pkg/server"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// ServerCmd is the main server command
var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Server hosts the web site",
	Long:  `A base server interface running on a device`,
	Run: func(cmd *cobra.Command, args []string) {
		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			logrus.Errorf("bad port: %v", err)
		}
		logrus.Infof("listening on port: %v", port)
		serverInstance := server.ServerInstance{ServerPort: port}

		serverInstance.Run()
	},
}
