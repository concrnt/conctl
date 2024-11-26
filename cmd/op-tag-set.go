package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/totegamma/concurrent/core"
)

var setCmd = &cobra.Command{
	Use:        "set",
	Short:      "Set the tag of the entity",
	Args:       cobra.ExactArgs(2),
	ArgAliases: []string{"ccid", "tag"},
	Run: func(cmd *cobra.Command, args []string) {
		ccid := args[0]
		tag := args[1]

		err := db.Model(&core.Entity{}).Where("id = ?", ccid).Update("tag", tag)
		if err.Error != nil {
			fmt.Println(err.Error)
		}
	},
}

func init() {
	opTagCmd.AddCommand(setCmd)
}
