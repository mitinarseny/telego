package cmd

import (
    "strings"

    "github.com/mitinarseny/telego/bot"
    "github.com/mitinarseny/telego/notifier"
    log "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    tb "gopkg.in/tucnak/telebot.v2"
)

var rootCmd = &cobra.Command{
    Run: func(cmd *cobra.Command, args []string) {
        checkMandatoryParams()
        checkDependentParams()
        start()
    },
}

func start() {
    var (
        notifierBot notifier.StatusNotifier
    )

    tgEndpoint := viper.GetString("telegram.endpoint")

    botToken := viper.GetString("bot.token")

    b, err := tb.NewBot(tb.Settings{
        URL:    tgEndpoint,
        Token:  botToken,
        Poller: nil,
        Reporter: func(err error) {
            log.WithFields(log.Fields{
                "context": "NOTIFIER",
            }).Error(err)
        },
    })
    if err != nil {
        log.WithFields(log.Fields{
            "context": "BOT",
            "action":  "AUTHENTICATE",
        }).Fatal(err)
    }

    if _, err := bot.Configure(b); err != nil {
        log.WithFields(log.Fields{
            "context": "BOT",
            "action":  "CONFIGURE",
        }).Fatal(err)
    }

    notifierToken := viper.GetString("notifier.bot.token")
    if notifierToken != "" {
        nb, err := tb.NewBot(tb.Settings{
            URL:   tgEndpoint,
            Token: notifierToken,
        })
        if err != nil {
            log.WithFields(log.Fields{
                "context": "NOTIFIER",
                "action":  "AUTHENTICATE",
            }).Fatal(err)
        }
        notifierBot = &notifier.Bot{
            Bot: nb,
            Chat: &tb.Chat{
                ID: viper.GetInt64("notifier.chat.id"),
            },
        }
    }

    log.WithFields(log.Fields{
        "context": "BOT",
        "action":  "STARTING",
    }).Info()

    if notifierBot != nil {
        if err := notifierBot.NotifyUp(b.Me.Username); err != nil {
            log.WithFields(log.Fields{
                "context": "NOTIFIER",
                "action":  "NOTIFY",
            }).Error(err)
        }
        defer func() {
            if err := notifierBot.NotifyDown(b.Me.Username); err != nil {
                log.WithFields(log.Fields{
                    "context": "NOTIFIER",
                    "action":  "NOTIFY",
                }).Error(err)
            }
        }()
    }

    b.Start()
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

func checkDependentParams() {
    notifierBotToken := viper.GetString("notifier.bot.token")
    notifierBotChatID := viper.GetInt64("notifier.chat.id")
    if (notifierBotToken == "") != (notifierBotChatID == 0) {
        log.Fatalf("%s must be provided simultaneously", strings.Join([]string{
            "notifier.bot.token",
            "notifier.chat.id",
        }, ", "))
    }
}
