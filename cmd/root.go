package cmd

import (
	"fmt"
	"os"

	"github.com/gabhendm/gimli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// The name of our config file, without the file extension because viper supports many different config file languages.
	defaultConfigFilename = ".gimli"

	// The environment variable prefix of all environment variables bound to our command line flags.
	// For example, --number is bound to STING_NUMBER.
	envPrefix = "GIMLI"
)

var (
	cfgFile string
	debug   bool
	rootCmd = &cobra.Command{
		Use:   "gimli",
		Short: "Gimli is a recon orchestration tool ",
		Long:  ` A fast and easy to use recon tool orchestator Complete documentation is available at https://github.com/gabhendm/gimli`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Run gimli init to initialize all required services before scanning")
		},
	}
)

//Execute is used to bootstrap commands tree
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gimli.yaml)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Debug Output Flag (True or False)")
}

func initializeConfig(cmd *cobra.Command) error {
	v := viper.New()
	v.SetConfigType("yaml")

	if cfgFile == "" {
		v.AddConfigPath("$HOME/")
		v.SetConfigName(".gimli")
	} else {
		v.SetConfigFile(cfgFile)
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	v.SetEnvPrefix(envPrefix)

	v.AutomaticEnv()

	utils.BindFlags(cmd, v, envPrefix)

	fmt.Println("Using config file:", v.ConfigFileUsed())

	return nil
}
