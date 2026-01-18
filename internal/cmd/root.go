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
	ConvertCmd.Flags().StringP("pages", "p", "1-N", "Comma separated list of page ranges. \"N\" is the last page.")
	ConvertCmd.SilenceErrors = true
	RootCmd.AddCommand(ConvertCmd)
	RootCmd.AddCommand(CutCmd)
	RootCmd.AddCommand(RotateCmd)
	RootCmd.AddCommand(DeskewCmd)
	RootCmd.AddCommand(PaddingCmd)
	RootCmd.AddCommand(DenoiseCmd)
}
