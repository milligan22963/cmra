// Package subcmd is for all subcmds in the cmd tree
package subcmd

import (
	"github.com/milligan22963/cmra/pkg/server"

	"github.com/spf13/cobra"
)

// ServerCmd is the main server command
var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Server hosts the web site",
	Long:  `A base server interface running on a device`,
	Run: func(cmd *cobra.Command, args []string) {
		serverInstance := server.ServerInstance{}

		serverInstance.Run()
	},
}
