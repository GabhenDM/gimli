package cmd

import (
	"fmt"

	"github.com/gabhendm/gimli/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(InitCmd)
	InitCmd.Flags().StringP("mongo-host-port", "p", viper.GetString("MongoHostPort"), "Set Host Port for mongoDB connection")
	viper.BindPFlag("MongoHostPort", InitCmd.PersistentFlags().Lookup("mongo-host-port"))
}

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Init Gimli Infrastucture Containers",
	Long:  `This command initializes all service containers required by Gimli (MongoDB, Redis, etc..)`,
	Run: func(cmd *cobra.Command, args []string) {
		debug, _ := cmd.Flags().GetBool("debug")
		mongoHostPort, _ := cmd.Flags().GetString("MongoHostPort")
		fmt.Println(fmt.Sprintf("[!] Starting MongoDB on Port %s...", mongoHostPort))
		mongodb, err := service.StartContainerDetached("mongo", []string{}, mongoHostPort, "27017", debug)
		if err != nil {
			panic(err)
		}

		fmt.Println("[!] MongoDB Started! - ID:", mongodb.ID)

	},
}
