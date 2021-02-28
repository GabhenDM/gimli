package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {

	rootCmd.AddCommand(VersionCmd)

}

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Gimli",
	Long:  `Gotta version all the things`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Gimli Recon Orchestrator v0.1 -- HEAD")
	},
}
