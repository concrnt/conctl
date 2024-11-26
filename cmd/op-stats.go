package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/totegamma/concurrent/core"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show statistics of the server",
	Run: func(cmd *cobra.Command, args []string) {
		var messageCount int64
		db.Model(&core.Message{}).Count(&messageCount)
		fmt.Printf("Message count: %d\n", messageCount)
		var associationCount int64
		db.Model(&core.Association{}).Count(&associationCount)
		fmt.Printf("Association count: %d\n", associationCount)
		var entityCount int64
		db.Model(&core.Entity{}).Count(&entityCount)
		fmt.Printf("Entity count: %d\n", entityCount)
		var entityMetaCount int64
		db.Model(&core.EntityMeta{}).Count(&entityMetaCount)
		fmt.Printf("EntityMeta count: %d\n", entityMetaCount)
	},
}

func init() {
	opCmd.AddCommand(statsCmd)
}
