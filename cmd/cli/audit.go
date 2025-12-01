package cli

import (
	"github.com/Louisrca/bloatfish/cmd"
	"github.com/Louisrca/bloatfish/internal/analyser"

	"github.com/spf13/cobra"
)

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Audit your web project for environmental impact",
	Long:  `Run an audit on your web project to analyze its environmental impact based on GreenIT and RGESN principles.`,
	Run: func(cmd *cobra.Command, args []string) {
		analyser := analyser.UnusedPackageAnalyser{}
		analyser.Analyze("./node_modules")
	},
}

func init() {
	cmd.RootCmd.AddCommand(auditCmd)
}
