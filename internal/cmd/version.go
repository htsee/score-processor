package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "dev"

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("sp %s\n", Version)
	},
}
