package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/krobus00/storage-service/internal/config"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "storage-service",
	Short: "storage-service",
	Long:  "storage-service",
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func Init() {
	if err := config.LoadConfig(); err != nil {
		log.Fatalln(err.Error())
	}

	log.Info(fmt.Sprintf("starting %s:%s...", config.ServiceName(), config.ServiceVersion()))
}
