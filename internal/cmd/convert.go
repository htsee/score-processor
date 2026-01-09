package cmd

import (
	"errors"
	"fmt"
	"io/fs"
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
		pages, err := cmd.Flags().GetString("pages")
		if err != nil {
			return err
		}

		if path.Ext(pdf) != ".pdf" {
			return fmt.Errorf("File %q is not a PDF", pdf)
		}

		if _, err := os.Stat(pdf); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return fmt.Errorf("File %q does not exist", pdf)
			}
			return fmt.Errorf("Cannot access file %q: %w", pdf, err)
		}

		info, err := os.Stat(output)
		if err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				return fmt.Errorf("Cannot access folder %q: %w", output, err)
			}
			err = os.Mkdir(output, 0755)
			if err != nil {
				return fmt.Errorf("Cannot create folder %q: %w", output, err)
			}
		}

		if !info.IsDir() {
			return fmt.Errorf("%q is not a folder", output)
		}

		output_path := fmt.Sprintf("%s/img_%%03d.png", output)
		convert := exec.Command("mutool", "convert", "-o", output_path, pdf, pages)
		if err := convert.Run(); err != nil {
			if errors.Is(err, exec.ErrNotFound) {
				return errors.New("\"mutool\" not found. Install \"muPDF\" to use this command")
			}
			return fmt.Errorf("Failed to convert PDF: %w", err)
		}
		return nil
	},
}
