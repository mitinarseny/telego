package cmd

import (
    "context"
    "os"
    "os/signal"
    "strings"
    "syscall"
    "time"

    "github.com/mitinarseny/telego/bot"
    "github.com/mitinarseny/telego/notifier"
    log "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    tb "gopkg.in/tucnak/telebot.v2"
)

const (
    debugKey             = "debug"
    botTokenKey          = "bot.token"
    notifierBotTokenKey  = "notifier.bot.token"
    notifierBotChatIDKey = "notifier.chat.id"
)

var rootCmd = &cobra.Command{
    Run: func(cmd *cobra.Command, args []string) {
        checkMandatoryParams()
        checkDependentParams()
        start()
    },
}

func start() {
    botLogEntry := log.WithField("context", "BOT")
    notifierLogEntry := log.WithField("context", "NOTIFIER")

    tgEndpoint := viper.GetString("telegram.endpoint")

    botToken := viper.GetString("bot.token")
    poller := &tb.LongPoller{
        Timeout: 60 * time.Second,
    }
    logPoller := tb.NewMiddlewarePoller(poller, func(update *tb.Update) bool {
        log.Debug(*update)
        return true
    })
    b, err := tb.NewBot(tb.Settings{
        URL:    tgEndpoint,
        Token:  botToken,
        Poller: logPoller,
        Reporter: func(err error) {
            notifierLogEntry.Error(err)
        },
    })
    if err != nil {
        botLogEntry.WithField("action", "AUTHENTICATE").Fatal(err)
    }
    botLogEntry.WithFields(log.Fields{
        "action":  "AUTHENTICATE",
        "account": b.Me.Username,
    }).Info()

    if _, err := bot.Configure(b); err != nil {
        botLogEntry.WithField("action", "CONFIGURE").Fatal(err)
    }
    statusNotifier := notifier.NewBaseNotifier()

    notifierToken := viper.GetString("notifier.bot.token")
    if notifierToken != "" {
        nb, err := tb.NewBot(tb.Settings{
            URL:   tgEndpoint,
            Token: notifierToken,
        })
        if err != nil {
            notifierLogEntry.WithField("action", "AUTHENTICATE").Fatal(err)
        }
        notifierLogEntry.WithFields(log.Fields{
            "action":  "AUTHENTICATE",
            "account": nb.Me.Username,
        }).Info()
        statusNotifier.Register(&notifier.Bot{
            Bot: nb,
            Chat: &tb.Chat{
                ID: viper.GetInt64("notifier.chat.id"),
            },
        })
    }
    ctx, cancelFunc := context.WithCancel(context.Background())
    defer cancelFunc()

    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

    statusNotifier.Start(ctx)
    defer statusNotifier.Stop()

    go func() {
        botLogEntry.WithField("action", "START").Info()
        b.Start()
    }()
    defer func() {
        b.Stop()
        botLogEntry.WithFields(log.Fields{
            "action":       "STOP",
            "lastUpdateID": poller.LastUpdateID,
        }).Info()
    }()

    statusNotifier.NotifyUp(b.Me.Username)
    defer statusNotifier.NotifyDown(b.Me.Username)

    gotSig := <-sigCh
    log.WithFields(log.Fields{
        "signal": gotSig.String(),
    }).Info("Got signal, stopping...")
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        log.Fatal(err)
    }
}

func init() {
    cobra.OnInitialize(initConfig)

    rootCmd.PersistentFlags().String(debugKey, "", "Debug mode")
    _ = viper.BindPFlag(debugKey, rootCmd.PersistentFlags().Lookup(debugKey))

    rootCmd.PersistentFlags().String(botTokenKey, "", "Telegram Bot API token")
    _ = viper.BindPFlag(botTokenKey, rootCmd.PersistentFlags().Lookup(botTokenKey))

    rootCmd.PersistentFlags().String(notifierBotTokenKey, "", "Notifier Telegram Bot API token")
    _ = viper.BindPFlag(notifierBotTokenKey, rootCmd.PersistentFlags().Lookup(notifierBotTokenKey))

    rootCmd.PersistentFlags().Int64(notifierBotChatIDKey, -1, "Notifier Chat ID")
    _ = viper.BindPFlag(notifierBotChatIDKey, rootCmd.PersistentFlags().Lookup(notifierBotChatIDKey))

    if viper.GetString(debugKey) != "" {
        log.SetLevel(log.DebugLevel)
        log.SetReportCaller(true)
    }
}

func initConfig() {
    viper.SetEnvPrefix("TELEGO")
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
    viper.AutomaticEnv()
}

func checkMandatoryParams() {
    var missing []string

    if v := viper.GetString(botTokenKey); v == "" {
        missing = append(missing, botTokenKey)
    }

    if len(missing) > 0 {
        log.Fatalf("missing: %s", strings.Join(missing, ", "))
    }
}

func checkDependentParams() {
    notifierBotToken := viper.GetString(notifierBotTokenKey)
    notifierBotChatID := viper.GetInt64(notifierBotChatIDKey)
    if (notifierBotToken == "") != (notifierBotChatID == 0) {
        log.Fatalf("%s must be provided simultaneously", strings.Join([]string{
            notifierBotTokenKey,
            notifierBotChatIDKey,
        }, ", "))
    }
}
