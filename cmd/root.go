package cmd

import (
    "context"
    "database/sql"
    "net/url"
    "os"
    "os/signal"
    "strconv"
    "strings"
    "syscall"
    "time"

    "github.com/mitinarseny/telego/bot"
    "github.com/mitinarseny/telego/notifier"
    "github.com/mitinarseny/telego/tg_log/clickhouse"
    "github.com/mitinarseny/telego/tg_log/repository"
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
    logDBHostKey         = "log.db.host"
    logDBPortKey         = "log.db.port"
    logDBUserKey         = "log.db.user"
    logDBPasswordKey     = "log.db.password"
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

    botToken := viper.GetString(botTokenKey)
    clickhouseDB, err := getClickHouseDB(
        viper.GetString(logDBHostKey),
        viper.GetInt(logDBPortKey),
        viper.GetString(logDBUserKey),
        viper.GetString(logDBPasswordKey),
        "log",
    )
    if err != nil {
        log.WithFields(log.Fields{
            "context": "CLICKHOUSE",
            "action":  "CONNECT",
        }).Fatal(err)
    }
    defer func() {
        if err := clickhouseDB.Close(); err != nil {
            log.WithFields(log.Fields{
                "context": "CLICKHOUSE",
                "action":  "CLOSE",
            }).Error(err)
        }
    }()
    updateLogger, err := clickhouse.NewBufferedUpdateLogger(clickhouseDB)
    if err != nil {
        log.WithFields(log.Fields{
            "context": "UpdatesLogger",
            "action":  "CREATE",
        }).Fatal(err)
    }
    defer updateLogger.Close()

    poller := &tb.LongPoller{
        Timeout: 60 * time.Second,
    }
    logPoller := tb.NewMiddlewarePoller(poller, func(update *tb.Update) bool {
        go func() {
            if err := updateLogger.LogUpdate(repository.FromTelebotUpdate(update)); err != nil {
                log.WithFields(log.Fields{
                    "context": "UpdatesLogger",
                    "action":  "LOG",
                }).Error(err)
            }
        }()
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

    notifierToken := viper.GetString(notifierBotTokenKey)
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
                ID: viper.GetInt64(notifierBotChatIDKey),
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

    rootCmd.PersistentFlags().String(logDBHostKey, "", "Host of DB for logging telegram updates")
    _ = viper.BindPFlag(logDBHostKey, rootCmd.PersistentFlags().Lookup(logDBHostKey))

    rootCmd.PersistentFlags().String(logDBPortKey, "", "Port of DB for logging telegram updates")
    _ = viper.BindPFlag(logDBPortKey, rootCmd.PersistentFlags().Lookup(logDBPortKey))

    rootCmd.PersistentFlags().String(logDBUserKey, "", "User of DB for logging telegram updates")
    _ = viper.BindPFlag(logDBUserKey, rootCmd.PersistentFlags().Lookup(logDBUserKey))

    rootCmd.PersistentFlags().String(logDBPasswordKey, "", "Password for user of DB for logging telegram updates")
    _ = viper.BindPFlag(logDBPasswordKey, rootCmd.PersistentFlags().Lookup(logDBPasswordKey))

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

func getClickHouseDB(host string, port int, username, password, dbName string) (*sql.DB, error) {
    connURL := url.URL{
        Scheme: "tcp",
        Host:   host,
        Path:   dbName,
    }
    if port != 0 {
        connURL.Host += ":" + strconv.Itoa(port)
    }
    if username != "" {
        connURL.RawQuery += "&username=" + username
    }
    if password != "" {
        connURL.RawQuery += "&password=" + password
    }

    db, err := sql.Open("clickhouse", connURL.String())
    if err != nil {
        return nil, err
    }
    if err := db.Ping(); err != nil {
        return nil, err
    }
    return db, nil
}
