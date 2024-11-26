/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/spf13/cobra"
)

// generateVapidKeysCmd represents the generateVapidKeys command
var generateVapidKeysCmd = &cobra.Command{
	Use:   "generateVapidKeys",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		vapidPrivateKey, vapidPublicKey, err := webpush.GenerateVAPIDKeys()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("VAPID private key:", vapidPrivateKey)
		fmt.Println("VAPID public key:", vapidPublicKey)
	},
}

func init() {
	generateCmd.AddCommand(generateVapidKeysCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateVapidKeysCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateVapidKeysCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
