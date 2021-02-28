package cmd

import (
	"github.com/gabhendm/gimli/service"
	"github.com/spf13/cobra"
)

var ListContainersCmd = &cobra.Command{
	Use:   "list-containers",
	Short: "List all running containers",
	Long:  `This command lists all running containers `,
	Run: func(cmd *cobra.Command, args []string) {
		service.ListRunningContainers()
	},
}
