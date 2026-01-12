package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

var ConvertCmd = &cobra.Command{
	Use:   "convert [inputs] [destination]",
	Short: "Convert PDF to images",
	Long:  "Convert PDF to images (PNG). Automatically create destination folder if it does not exist. Require \"muPDF\".",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputs := args[0 : len(args)-1]
		destination := args[len(args)-1]
		pages, err := cmd.Flags().GetString("pages")
		if err != nil {
			return err
		}
		for _, input := range inputs {
			if err = Convert(input, destination, pages); err != nil {
				return err
			}
		}
		return nil
	},
}

func Convert(pdf, destination, pages string) error {
	if path.Ext(pdf) != ".pdf" {
		return fmt.Errorf("File %q is not a PDF", pdf)
	}

	if _, err := os.Stat(pdf); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("File %q does not exist", pdf)
		}
		return fmt.Errorf("Cannot access file %q: %w", pdf, err)
	}

	if err := os.MkdirAll(destination, 0755); err != nil {
		return fmt.Errorf("Cannot create folder %q: %w", destination, err)
	}

	pdf_name, _ := strings.CutSuffix(path.Base(pdf), ".pdf")
	output_path := fmt.Sprintf("%s/%s_%%03d.png", destination, pdf_name)

	convert := exec.Command("mutool", "convert", "-o", output_path, pdf, pages)
	if err := convert.Run(); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return errors.New("\"mutool\" not found. Install \"muPDF\" to use this command")
		}
		return fmt.Errorf("Failed to convert PDF: %w", err)
	}
	return nil
}
