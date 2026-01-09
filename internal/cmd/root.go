package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use: "sp",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	ConvertCmd.SilenceErrors = true
	RootCmd.AddCommand(ConvertCmd)
}
