package cmd

import "github.com/spf13/cobra"

var SpliceCmd = &cobra.Command{
	Use:   "splice [inputs] [destination]",
	Short: "Splice images",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputs := args[0 : len(args)-1]
		destination := args[len(args)-1]
		for _, input := range inputs {
			if err := Splice(input, destination); err != nil {
				return err
			}
		}
		return nil
	},
}

func Splice(input, destination string) error {
	return nil
}
