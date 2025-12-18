package cmd

import (
	"fmt"

	"github.com/Louisrca/bloatfish/internal/audit"
	"github.com/spf13/cobra"
)

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Audit your web project for environmental impact",
	Long:  `Run an audit on your web project to analyze its environmental impact based on GreenIT and RGESN principles.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running full audit...")
		audit.DeepAudit("https://vite.dev/")
	},
}

func init() {
	RootCmd.AddCommand(auditCmd)
}
