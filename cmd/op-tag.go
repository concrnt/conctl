package cmd

import (
	"github.com/spf13/cobra"
)

var opTagCmd = &cobra.Command{
	Use: "tag",
}

func init() {
	opCmd.AddCommand(opTagCmd)
}
