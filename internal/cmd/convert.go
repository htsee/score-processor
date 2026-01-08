package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/spf13/cobra"
)

var ConvertCmd = &cobra.Command{
	Use:   "convert [input] [output]",
	Short: "Convert PDF to images",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		pdf := args[0]
		output := args[1]

		if path.Ext(pdf) != ".pdf" {
			return fmt.Errorf("File %q is not a PDF", pdf)
		}

		if _, err := os.Stat(pdf); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("File %q does not exist", pdf)
			} else {
				return fmt.Errorf("Cannot access file %q: %w", pdf, err)
			}
		}

		info, err := os.Stat(output)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("Folder %q does not exist", output)
			} else {
				return fmt.Errorf("Cannot access folder %q: %w", output, err)
			}
		}

		if !info.IsDir() {
			return fmt.Errorf("%q is not a folder", output)
		}

		output_path := fmt.Sprintf("%s/img_%%03d.png", output)
		convert := exec.Command("mutool", "convert", "-o", output_path, pdf)
		if err := convert.Run(); err != nil {
			return err
		}
		return nil
	},
}
