package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/stonebirdjx/go-layout/internal/bootstrap"
)

var cfgFile string

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the backend application server",
	Long:  `Start the backend application server with lifecycle management and dependency injection.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := bootstrap.RunServer(cfgFile); err != nil {
			fmt.Fprintf(os.Stderr, "Server failed: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.
	serverCmd.Flags().StringVarP(&cfgFile, "config", "c", "configs/server.yaml", "config file path")
}
