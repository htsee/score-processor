package cmd

import "github.com/spf13/cobra"

var VSpliceCmd = &cobra.Command{
	Use:   "vsplice [inputs] [destination]",
	Short: "Splice 2 vertical images",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputs := args[0 : len(args)-1]
		destination := args[len(args)-1]
		for _, input := range inputs {
			if err := VSplice(input, destination); err != nil {
				return err
			}
		}
		return nil
	},
}

func VSplice(input, destination string) error {
	return nil
}
