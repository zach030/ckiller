package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"os"
)

const (
	openaiApiKey  = "openaiApiKey"
	openaiBaseURL = "openaiBaseURL"
)

var (
	cfgFile  string
	language string
	cfg      *CkConfig
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ckiller",
	Short: "kill command `-h`",
}

type CkConfig struct {
	OpenaiApiKey  string `mapstructure:"openaiApiKey"`
	OpenaiBaseURL string `mapstructure:"openaiBaseURL"`
	Language      string `mapstructure:"language"`
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&language, "language", "English", "language")
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err != nil {
			log.Panic(err)
		}
	}
	viper.AutomaticEnv() // read in environment variables that match
	cfg = &CkConfig{
		OpenaiApiKey:  viper.GetString(openaiApiKey),
		OpenaiBaseURL: viper.GetString(openaiBaseURL),
		Language:      language,
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
