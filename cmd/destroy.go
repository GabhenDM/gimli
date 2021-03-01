package cmd

import (
	"fmt"

	"github.com/gabhendm/gimli/service"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(DestroyCmd)
}

// DestroyCmd is used to Destroy all Service Containers
var DestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy Gimli Infrastucture Containers",
	Long:  `This command stops and removes all service containers required by Gimli (MongoDB, Redis, etc..)`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[!] Stopping all Gimli containers")
		err := service.RemoveRunningContainers()
		if err != nil {
			panic(err)
		}
		fmt.Println("[!] Containers Stopped!")
	},
}
