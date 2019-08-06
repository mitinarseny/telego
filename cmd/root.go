package cmd

import (
	"strings"

	"github.com/mitinarseny/telego/bot_old"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
		checkMandatoryParams()
		start()
	},
}

func start() {
	if err := bot_old.Run(viper.GetString("bot.token"),
		viper.GetString("notifier.bot.token"),
		viper.GetInt64("notifier.chat.id"),
		true); err != nil {
		log.Fatal(err)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().String("bot.token", "", "Telegram Bot API token")
	_ = viper.BindPFlag("bot.token", rootCmd.PersistentFlags().Lookup("bot.token"))

	rootCmd.PersistentFlags().String("notifier.bot.token", "", "Notifier Telegram Bot API token")
	_ = viper.BindPFlag("notifier.bot.token", rootCmd.PersistentFlags().Lookup("notifier.bot.token"))

	rootCmd.PersistentFlags().Int64("notifier.chat.id", -1, "Notifier Chat ID")
	_ = viper.BindPFlag("notifier.chat.id", rootCmd.PersistentFlags().Lookup("notifier.chat.id"))
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
		log.Fatalf("missing: %s", strings.Join(missing, ", "))
	}
}
