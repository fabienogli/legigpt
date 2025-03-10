package main

import (
	"os"

	"github.com/fabienogli/legigpt/cmd/findsimilitude"
	"github.com/fabienogli/legigpt/cmd/retrieve"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "legi-cli",
	Short: "Legi allows to search legi API",
	// Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		return nil
	},
}

func main() {
	rootCmd.AddCommand(
		retrieve.InitCMD(),
		findsimilitude.InitCMD(),
	)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
