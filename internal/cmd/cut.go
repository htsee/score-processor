package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var CutCmd = &cobra.Command{
	Use:   "cut [input] [destination]",
	Short: "Cut the score into staves",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := args[0]
		destination := args[1]
		fmt.Printf("Cut %q and output to %q.", input, destination)
		return nil
	},
}
