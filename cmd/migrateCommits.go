/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"github.com/totegamma/concurrent/cdid"
	"github.com/totegamma/concurrent/core"
)

// migrateCommitsCmd represents the migrateCommits command
var migrateCommitsCmd = &cobra.Command{
	Use:   "commits-redisstream-to-db",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		if config == nil {
			fmt.Println("Config must be loaded")
			return
		}

		ctx := context.Background()

		result, err := rdb.XRead(ctx, &redis.XReadArgs{
			Streams: []string{
				"repository-all",
				"0",
			},
			Block: 0,
		}).Result()

		if err != nil {
			panic(err)
		}

		if len(result) != 1 {
			fmt.Println("Invalid result")
			return
		}

		messages := result[0].Messages

		var processed = 0
		// print progress
		go func() {
			for {
				fmt.Printf("Processed %d messages of %d\n", processed, len(messages))
				time.Sleep(1 * time.Second)
			}
		}()

		var created = 0
		var updated = 0
		var skipped = 0
		var errored = 0
		for _, message := range messages {

			owner, ok := message.Values["owner"].(string)
			if !ok {
				errored++
				continue
			}

			content, ok := message.Values["entry"].(string)
			if !ok {
				errored++
				continue
			}

			split := strings.Split(content, " ")

			signature := split[0]
			document := strings.Join(split[1:], " ")

			var base core.DocumentBase[any]
			err := json.Unmarshal([]byte(document), &base)
			if err != nil {
				errored++
				continue
			}

			hash := core.GetHash([]byte(document))
			hash10 := [10]byte{}
			copy(hash10[:], hash[:10])
			signedAt := base.SignedAt
			documentID := cdid.New(hash10, signedAt).String()

			if !core.IsCCID(owner) {
				if owner == config.FQDN {
					owner = config.CSID
				} else {
					skipped++
					continue
				}
			}

			// check if the document already exists
			var existing CommitLog
			db.Where("document_id = ?", documentID).First(&existing)
			if existing.ID != 0 { // already exists
				// check if the owner already exists
				if !slices.Contains(existing.Owners, owner) {
					existing.Owners = append(existing.Owners, owner)
					db.Save(&existing)
					updated++
				} else {
					skipped++
				}
			} else { // create new
				log := CommitLog{
					DocumentID:  documentID,
					IsEphemeral: false,
					Type:        base.Type,
					Document:    document,
					Signature:   signature,
					SignedAt:    base.SignedAt,
					Owners:      []string{owner},
				}

				err := db.Create(&log).Error
				if err != nil {
					errored++
					continue
				}
				created++
			}

			processed++
		}

		fmt.Printf("Processed %d messages of %d\n", processed, len(messages))
		fmt.Printf("Created: %d\n", created)
		fmt.Printf("Updated: %d\n", updated)
		fmt.Printf("Skipped: %d\n", skipped)
		fmt.Printf("Errored: %d\n", errored)
	},
}

type CommitLog struct {
	ID          uint           `json:"id" gorm:"primaryKey;auto_increment"`
	IP          string         `json:"ip" gorm:"type:text"`
	DocumentID  string         `json:"documentID" gorm:"type:char(26);uniqueIndex:idx_document_id"`
	IsEphemeral bool           `json:"isEphemeral" gorm:"type:boolean;default:false"`
	Type        string         `json:"type" gorm:"type:text"`
	Document    string         `json:"document" gorm:"type:json"`
	Signature   string         `json:"signature" gorm:"type:char(130)"`
	SignedAt    time.Time      `json:"signedAt" gorm:"type:timestamp with time zone;not null;default:clock_timestamp()"`
	Owners      pq.StringArray `json:"owners" gorm:"type:char(42)[]"`
	CDate       time.Time      `json:"cdate" gorm:"type:timestamp with time zone;not null;default:clock_timestamp()"`
}

func init() {
	migrationCmd.AddCommand(migrateCommitsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrateCommitsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// migrateCommitsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
