/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/totegamma/concurrent/core"
)

// statsCmd represents the stats command
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
	rootCmd.AddCommand(statsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
