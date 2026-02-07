package cmd

import (
	"errors"
	"fmt"
	"os/exec"
	"path"
	"strings"

	"github.com/htsee/score-processor/internal/util"
	"github.com/spf13/cobra"
)

var ConvertCmd = &cobra.Command{
	Use:   "convert [inputs] [destination]",
	Short: "Convert PDFs to images. Requires mupdf",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputs := args[0 : len(args)-1]
		destination := args[len(args)-1]
		pages, err := cmd.Flags().GetString("pages")
		if err != nil {
			return err
		}
		if err := util.CheckValidIO(inputs, "pdf", destination); err != nil {
			return err
		}
		return util.Batch(inputs, func(input string) error {
			return Convert(input, destination, pages)
		})
	},
}

func Convert(pdf, destination, pages string) error {
	pdf_name, _ := strings.CutSuffix(path.Base(pdf), ".pdf")
	output_path := fmt.Sprintf("%s/%s_%%03d.png", destination, pdf_name)

	convert := exec.Command("mutool", "convert", "-o", output_path, pdf, pages)
	if err := convert.Run(); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return errors.New("\"mutool\" not found. Install \"muPDF\" to use this command")
		}
		return fmt.Errorf("failed to convert PDF: %w", err)
	}
	return nil
}
