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
    botToken string
)

var rootCmd = &cobra.Command{
    Run: func(cmd *cobra.Command, args []string) {
        checkMandatoryParams()
        start()
    },
}

func start() {
    if err := bot.Run(viper.GetString("bot.token"),
        viper.GetString("notifier.bot.token"),
        viper.GetInt64("notifier.chat.id"),
        true); err != nil {
        fmt.Println(err)
        os.Exit(1)
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
        fmt.Printf("missing: %s", strings.Join(missing, ", "))
        os.Exit(1)
    }
}
