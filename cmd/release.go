package cmd

import (
	"github.com/spf13/cobra"

	"go.amplifyedge.org/booty-v2/dep"
)

// runs Goreleaser under the hood
func ReleaseCommand(a dep.Agent) *cobra.Command {
	runCmd := &cobra.Command{Use: "release", DisableFlagParsing: true}
	runCmd.DisableFlagParsing = true
	runCmd.Flags().SetInterspersed(true)
	runCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return a.Run("goreleaser", args...)
	}
	return runCmd
}
