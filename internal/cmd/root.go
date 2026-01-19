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

	CutCmd.SilenceErrors = true
	RootCmd.AddCommand(CutCmd)

	RotateCmd.SilenceErrors = true
	RootCmd.AddCommand(RotateCmd)

	DeskewCmd.SilenceErrors = true
	RootCmd.AddCommand(DeskewCmd)

	PaddingCmd.SilenceErrors = true
	RootCmd.AddCommand(PaddingCmd)

	DenoiseCmd.SilenceErrors = true
	RootCmd.AddCommand(DenoiseCmd)

	TrimCmd.SilenceErrors = true
	RootCmd.AddCommand(TrimCmd)
}
