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

	RotateCmd.Flags().Float64P("angle", "a", 0, "Angle (in degrees)")
	RotateCmd.SilenceErrors = true
	RootCmd.AddCommand(RotateCmd)

	DeskewCmd.SilenceErrors = true
	RootCmd.AddCommand(DeskewCmd)

	PadCmd.Flags().IntP("vpad", "V", 10, "Vertical padding (in mm)")
	PadCmd.Flags().IntP("hpad", "H", 10, "Horizontal padding (in mm)")
	PadCmd.SilenceErrors = true
	RootCmd.AddCommand(PadCmd)

	DenoiseCmd.Flags().IntP("size", "s", 2, "Size (radius) of noise removed (in mm)")
	DenoiseCmd.SilenceErrors = true
	RootCmd.AddCommand(DenoiseCmd)

	TrimCmd.Flags().IntP("top", "t", 0, "Trim top (in mm)")
	TrimCmd.Flags().IntP("bottom", "b", 0, "Trim bottom (in mm)")
	TrimCmd.Flags().IntP("left", "l", 0, "Trim left (in mm)")
	TrimCmd.Flags().IntP("right", "r", 0, "Trim right (in mm)")
	TrimCmd.SilenceErrors = true
	RootCmd.AddCommand(TrimCmd)

	SpliceCmd.SilenceErrors = true
	RootCmd.AddCommand(SpliceCmd)
}
