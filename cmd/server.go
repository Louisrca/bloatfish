package cmd

import (
	"github.com/Louisrca/bloatfish/internal/server"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the Bloatfish web server",
	Long:  `Start the Bloatfish web server to view and edit pages.`,
	Run: func(cmd *cobra.Command, args []string) {
		server.StartServer()
	},
}

func init() {
	RootCmd.AddCommand(serverCmd)
}
