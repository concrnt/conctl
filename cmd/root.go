/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-yaml/yaml"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/totegamma/concurrent/core"
)

var (
	db     *gorm.DB
	rdb    *redis.Client
	config *core.Config
)

type Config struct {
	Concrnt core.ConfigInput `yaml:"concrnt"`
}

func (c *Config) Load(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	err = yaml.NewDecoder(f).Decode(&c)
	if err != nil {
		return err
	}

	return nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "conctl",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		var err error

		logger := logger.New(
			log.New(os.Stderr, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:        time.Second,
				LogLevel:             logger.Silent,
				Colorful:             true,
				ParameterizedQueries: true,
			},
		)

		dbhost, err := cmd.Flags().GetString("dbhost")
		if err != nil {
			return err
		}
		dbuser, err := cmd.Flags().GetString("dbuser")
		if err != nil {
			return err
		}
		dbpass, err := cmd.Flags().GetString("dbpass")
		if err != nil {
			return err
		}
		dbname, err := cmd.Flags().GetString("dbname")
		if err != nil {
			return err
		}
		dbport, err := cmd.Flags().GetString("dbport")
		if err != nil {
			return err
		}

		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			dbhost, dbuser, dbpass, dbname, dbport,
		)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger,
		})
		if err != nil {
			return err
		}

		redisAddr, _ := cmd.Flags().GetString("redisaddr")

		rdb = redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: "",
			DB:       0,
		})

		if rdb == nil {
			return fmt.Errorf("Failed to connect to redis")
		}

		configPath, _ := cmd.Flags().GetString("configpath")
		rootConf := Config{}
		err = rootConf.Load(configPath)
		if err == nil {
			conf := core.SetupConfig(rootConf.Concrnt)
			config = &conf
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("dbname", "d", "concrnt", "Database name")
	rootCmd.PersistentFlags().StringP("dbhost", "H", "localhost", "Database host")
	rootCmd.PersistentFlags().StringP("dbuser", "u", "postgres", "Database user")
	rootCmd.PersistentFlags().StringP("dbpass", "p", "postgres", "Database password")
	rootCmd.PersistentFlags().StringP("dbport", "P", "5432", "Database port")
	rootCmd.PersistentFlags().StringP("redisaddr", "r", "localhost:6379", "Redis address")
	rootCmd.PersistentFlags().StringP("configpath", "c", "/etc/concrnt/config/config.yaml", "Config file path")
}
