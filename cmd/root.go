package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "bloatfish",
	Short: "Modern eco-design analyzer for web projects — Reduce your environmental footprint with actionable insights",
	Long:  `Bloatfish is a command-line tool that audits your web projects against GreenIT and RGESN (French Eco-Design Reference Framework) principles. It combines a Go-based CLI with an ESLint plugin to provide comprehensive environmental impact analysis.`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatalf("Erreur lors de l'exécution : %v", err)
	}
}

func init() {
}
