package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/totegamma/concurrent/core"
)

var roleGetCmd = &cobra.Command{
	Use:        "get",
	Short:      "Get the tag of the entity",
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ccid"},
	Run: func(cmd *cobra.Command, args []string) {

		ccid := args[0]

		var entity core.Entity
		err := db.Where("id = ?", ccid).First(&entity)
		if err.Error != nil {
			fmt.Println(err.Error)
		} else {
			fmt.Printf(entity.Tag)
		}
	},
}

func init() {
	opTagCmd.AddCommand(roleGetCmd)
}
