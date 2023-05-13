package cmd

import (
	"github.com/krobus00/storage-service/internal/bootstrap"
	"github.com/spf13/cobra"
)

// initPermissionCmd represents the initPermission command.
var initPermissionCmd = &cobra.Command{
	Use:   "init-permission",
	Short: "init auth permission",
	Long:  `init auth permission`,
	Run: func(cmd *cobra.Command, args []string) {
		bootstrap.StartInitPermission()
	},
}

func init() {
	rootCmd.AddCommand(initPermissionCmd)
}
