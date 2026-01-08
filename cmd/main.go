package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "sp",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
}

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert PDF to images",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			cmd.Help()
			return nil
		}

		pdf := args[0]
		isPdf, _ := filepath.Match("*.pdf", pdf)

		if !isPdf {
			return fmt.Errorf("File %q is not a PDF", pdf)
		}

		if _, err := os.Stat(pdf); err != nil {
			return fmt.Errorf("File %q does not exist", pdf)
		}

		convert := exec.Command("mutool", "convert", "-o", "output/img_%03d.png", pdf)
		if err := convert.Run(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	convertCmd.SilenceErrors = true
	rootCmd.AddCommand(convertCmd)
}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}
