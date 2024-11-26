package cmd

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/totegamma/concurrent/core"
	"github.com/totegamma/concurrent/x/jwt"
)

var inviteCmd = &cobra.Command{
	Use:   "invite",
	Short: "generate invite code for the server",
	Run: func(cmd *cobra.Command, args []string) {

		configPath, _ := cmd.Flags().GetString("configpath")
		rootConf := Config{}
		err := rootConf.Load(configPath)
		if err != nil {
			fmt.Println("failed to open config file.", err)
			return
		}

		conf := core.SetupConfig(rootConf.Concrnt)
		config = &conf

		if config.PrivateKey == "" {
			fmt.Println("config.PrivateKey is empty")
			return
		}

		if config.FQDN == "" {
			fmt.Println("config.FQDN is empty")
			return
		}

		if config.CSID == "" {
			fmt.Println("config.CSID is empty")
			return
		}

		token, err := jwt.Create(jwt.Claims{
			Subject:        "CONCRNT_INVITE",
			JWTID:          uuid.New().String(),
			Audience:       config.FQDN,
			Issuer:         config.CSID,
			ExpirationTime: fmt.Sprintf("%d", time.Now().Add(time.Hour*24).Unix()),
		}, config.PrivateKey)

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(token)
	},
}

func init() {
	generateCmd.AddCommand(inviteCmd)
	inviteCmd.Flags().StringP("configpath", "c", "/etc/concrnt/config/config.yaml", "Config file path")
}
