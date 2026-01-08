package main

import (
	"github.com/htsee/score-processor/internal/cmd"
	"github.com/spf13/cobra"
)

func main() {
	cobra.CheckErr(cmd.RootCmd.Execute())
}
