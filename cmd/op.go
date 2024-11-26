package cmd

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-yaml/yaml"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"gorm.io/gorm"

	"github.com/totegamma/concurrent/core"
)

var (
	db     *gorm.DB
	rdb    *redis.Client
	config *core.Config
)

type Config struct {
	Concrnt core.ConfigInput `yaml:"concrnt"`
	Server  Server           `yaml:"server"`
}

type Server struct {
	Dsn           string `yaml:"dsn"`
	RedisAddr     string `yaml:"redisAddr"`
	RedisDB       int    `yaml:"redisDB"`
	MemcachedAddr string `yaml:"memcachedAddr"`

	RepositoryPath string `yaml:"repositoryPath"`
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

var opCmd = &cobra.Command{
	Use:     "operation",
	Aliases: []string{"op"},
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

		configPath, _ := cmd.Flags().GetString("configpath")
		if configPath == "" {
			configPath = "/etc/concrnt/config/config.yaml"
		}
		rootConf := Config{}
		err = rootConf.Load(configPath)
		if err == nil {
			conf := core.SetupConfig(rootConf.Concrnt)
			config = &conf
		}

		if rootConf.Server.Dsn != "" {
			split := strings.Split(rootConf.Server.Dsn, " ")
			for _, s := range split {
				if strings.Contains(s, "host=") {
					if dbhost == "" {
						dbhost = strings.Split(s, "=")[1]
					}
				} else if strings.Contains(s, "user=") {
					if dbuser == "" {
						dbuser = strings.Split(s, "=")[1]
					}
				} else if strings.Contains(s, "password=") {
					if dbpass == "" {
						dbpass = strings.Split(s, "=")[1]
					}
				} else if strings.Contains(s, "dbname=") {
					if dbname == "" {
						dbname = strings.Split(s, "=")[1]
					}
				} else if strings.Contains(s, "port=") {
					if dbport == "" {
						dbport = strings.Split(s, "=")[1]
					}
				}
			}
		}

		if dbhost == "" {
			dbhost = "localhost"
		}
		if dbuser == "" {
			dbuser = "postgres"
		}
		if dbpass == "" {
			dbpass = "postgres"
		}
		if dbname == "" {
			dbname = "concrnt"
		}
		if dbport == "" {
			dbport = "5432"
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

		if redisAddr == "" {
			redisAddr = rootConf.Server.RedisAddr
		}
		if redisAddr == "" {
			redisAddr = "localhost:6379"
		}

		rdb = redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: "",
			DB:       0,
		})
		if rdb == nil {
			return fmt.Errorf("Failed to connect to redis")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(opCmd)

	opCmd.PersistentFlags().StringP("dbname", "d", "", "Database name (default: concrnt)")
	opCmd.PersistentFlags().StringP("dbhost", "H", "", "Database host (default: localhost)")
	opCmd.PersistentFlags().StringP("dbuser", "u", "", "Database user (default: postgres)")
	opCmd.PersistentFlags().StringP("dbpass", "p", "", "Database password (default: postgres)")
	opCmd.PersistentFlags().StringP("dbport", "P", "", "Database port (default: 5432)")
	opCmd.PersistentFlags().StringP("redisaddr", "r", "", "Redis address (default: localhost:6379)")
	opCmd.PersistentFlags().StringP("configpath", "c", "", "Config file path (default: /etc/concrnt/config/config.yaml)")
}
