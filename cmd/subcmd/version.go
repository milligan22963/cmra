// Package subcmd is for all subcmds in the cmd tree
package subcmd

import (
	"github.com/milligan22963/cmra/pkg/version"

	"github.com/spf13/cobra"
)

// VersionCmd is the version command
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays the version of the application",
	Long:  `A command to display the version of the application`,
	Run: func(cmd *cobra.Command, args []string) {
		versionInstance := version.VersionInstance{}

		versionInstance.Run()
	},
}
