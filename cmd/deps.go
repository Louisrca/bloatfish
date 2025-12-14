package cmd

import (
	"fmt"

	"github.com/Louisrca/bloatfish/internal/analyzer"
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

		analyzer.WriteJSONReport(report)
	},
}

func init() {
	auditCmd.AddCommand(depsCmd)
}
