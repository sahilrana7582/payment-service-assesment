package config

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	once   sync.Once
	config *Config
)

func loadViper() {

	fmt.Println("Loading Config...")

	fs := pflag.NewFlagSet("payment-system", pflag.ContinueOnError)

	fs.String("name", "payment-svc", "Application Name")
	fs.String("port", "8080", "Application Port")
	fs.String("host", "localhost", "Database Host")
	fs.String("db-port", "5432", "Database Port")
	fs.String("user", "postgres", "Database User")
	fs.String("password", "", "Database Password")
	fs.String("db-name", "", "Database Name")
	fs.String("ssl-mode", "disabke", "SSL Mode")

	if err := fs.Parse(os.Args[1:]); err != nil {
		fmt.Printf("Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	if err := viper.BindPFlags(fs); err != nil {
		fmt.Printf("Error binding flags: %v\n", err)
		os.Exit(1)
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
}

func GetConfig() *Config {
	once.Do(func() {
		loadViper()

		config = &Config{
			App: Application{
				Name: viper.GetString("name"),
				Port: viper.GetString("port"),
			},
			Db: Database{
				Host:     viper.GetString("host"),
				Port:     viper.GetString("db-port"),
				User:     viper.GetString("user"),
				Password: viper.GetString("password"),
				Name:     viper.GetString("db-name"),
				SSLMode:  viper.GetString("ssl-mode"),
			},
		}
	})

	return config
}
