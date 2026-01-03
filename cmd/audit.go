package cmd

import (
	"fmt"

	"github.com/Louisrca/bloatfish/internal/analyzer"
	"github.com/Louisrca/bloatfish/internal/audit"
	"github.com/Louisrca/bloatfish/internal/server"
	utils "github.com/Louisrca/bloatfish/internal/utils"
	"github.com/spf13/cobra"
)

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Audit your web project for environmental impact",
	Long:  `Run an audit on your web project to analyze his dependencies and minimal web performance.`,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("🔍 Analyzing dependencies...")
		depReport, _ := analyzer.AnalyzeDependencies()

		fmt.Println("🔍 Running full audit...")
		deepAuditreport := audit.DeepAudit([]string{"http://localhost:3000/fr", "http://localhost:3000/fr/playlist", "http://localhost:3000/en/playlist", "http://localhost:3000/en"})

		fmt.Println("✅ Audit completed! Reports saved to 'full_audit_report.json'.")
		utils.WriteJSONReport(utils.JSONReport{DependencyReport: depReport, ChromeDPReport: deepAuditreport}, "full_audit_report.json")

		server.StartServer()

	},
}

func init() {
	RootCmd.AddCommand(auditCmd)
}
