package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mitinarseny/telego/bot"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	BotToken string
)

var rootCmd = &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
		checkMandatoryParams()
		start()
	},
}

func start() {
	if err := bot.Start(BotToken); err != nil {

	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&BotToken, "bot.token", "t", "", "Telegram Bot API token. You can get it from https://t.me/BotFather")
}

func initConfig() {
	viper.SetEnvPrefix("TELEGO")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

func checkMandatoryParams() {
	var missing []string

	if v := viper.GetString("bot.token"); v == "" {
		missing = append(missing, "bot.token")
	}

	if len(missing) > 0 {
		fmt.Printf("missing: %s", strings.Join(missing, ", "))
		os.Exit(1)
	}
}
