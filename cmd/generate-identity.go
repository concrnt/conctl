/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"encoding/hex"

	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/go-bip39"
	"github.com/spf13/cobra"
)

// identityCmd represents the identity command
var identityCmd = &cobra.Command{
	Use:   "identity",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		entropy, err := bip39.NewEntropy(128)

		mnemonic, err := bip39.NewMnemonic(entropy)
		if err != nil {
			panic(err)
		}

		seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
		if err != nil {
			panic(err)
		}
		master, ch := hd.ComputeMastersFromSeed(seed)
		priv, err := hd.DerivePrivateKeyForPath(master, ch, "m/44'/118'/0'/0/0")
		if err != nil {
			panic(err)
		}

		privHex := hex.EncodeToString(priv)

		privKey := &secp256k1.PrivKey{Key: priv}

		// get the public key
		pubKey := privKey.PubKey()
		pubKeyHex := hex.EncodeToString(pubKey.Bytes())

		fa := sdk.AccAddress(pubKey.Address())

		// convert the address to string
		addrCdc := address.NewBech32Codec("con")
		addrStr, err := addrCdc.BytesToString(fa)
		if err != nil {
			panic(err)
		}

		fmt.Println("ccid:\t\t", addrStr)
		fmt.Println("mnemonic:\t", mnemonic)
		fmt.Println("privatekey:\t", privHex)
		fmt.Println("publickey:\t", pubKeyHex)

	},
}

func init() {
	generateCmd.AddCommand(identityCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// identityCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// identityCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
