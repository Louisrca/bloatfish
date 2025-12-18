package cmd

import (
	"fmt"

	"github.com/Louisrca/bloatfish/internal/analyzer"
	utils "github.com/Louisrca/bloatfish/internal/utils"
	"github.com/spf13/cobra"
)

var depsCmd = &cobra.Command{
	Use:   "deps",
	Short: "Audit dependencies (direct, indirect, unused)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔍 Analyzing dependencies...")

		report, err := analyzer.AnalyzeDependencies()
		if err != nil {
			fmt.Println("❌ Error:", err)
			return
		}

		utils.WriteJSONReport(report, "unused_packages_report.json")
	},
}

func init() {
	auditCmd.AddCommand(depsCmd)
}
