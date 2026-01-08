package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
)

var ConvertCmd = &cobra.Command{
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
