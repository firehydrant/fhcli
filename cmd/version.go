package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cliRevision string
var releaseVersion string

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of FireHydrant CLI",
	Run: func(cmd *cobra.Command, args []string) {
		versionString := fmt.Sprintf("FireHydrant API CLI v%s - %s", releaseVersion, cliRevision)
		fmt.Println(versionString)
	},
}
