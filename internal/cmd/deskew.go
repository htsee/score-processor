package cmd

import (
	"github.com/htsee/score-processor/internal/util"
	"github.com/spf13/cobra"
)

var DeskewCmd = &cobra.Command{
	Use:   "deskew [inputs]",
	Short: "Deskew images",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, input := range args {
			if err := Deskew(input); err != nil {
				return err
			}
		}
		return nil
	},
}

func Deskew(input string) error {
	if err := util.CheckFileType(input, "png"); err != nil {
		return err
	}
	return nil
}
