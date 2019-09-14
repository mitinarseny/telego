package cmd

import (
    "context"
    "database/sql"
    "net/http"
    "net/url"
    "os"
    "os/signal"
    "strconv"
    "strings"
    "syscall"
    "time"

    "github.com/mitinarseny/telego/bot"
    "github.com/mitinarseny/telego/notifier"
    "github.com/mitinarseny/telego/tglog"
    "github.com/mitinarseny/telego/tglog/dblog"
    "github.com/mitinarseny/telego/tglog/repo"
    "github.com/mitinarseny/telego/tglog/repo/mongo"
    log "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    mongoDriver "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
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
    logDBNameKey         = "log.db.name"
)

var rootCmd = &cobra.Command{
    Run: func(cmd *cobra.Command, args []string) {
        if viper.GetBool(debugKey) {
            log.SetLevel(log.DebugLevel)
            log.SetReportCaller(true)
        }
        checkMandatoryParams()
        checkDependentParams()
        if err := start(); err != nil {
            log.Fatal(err)
        }
    },
}

func start() error {
    botLogEntry := log.WithField("context", "BOT")
    notifierLogEntry := log.WithField("context", "NOTIFIER")

    tgEndpoint := viper.GetString("telegram.endpoint")
    botToken := viper.GetString(botTokenKey)

    mongoOpts := options.Client().SetAppName("bot").SetAuth(options.Credential{
        Username: viper.GetString(logDBUserKey),
        Password: viper.GetString(logDBUserKey),
    }).SetHosts([]string{
        viper.GetString(logDBHostKey),
    })

    mongoConnectCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    mongoClient, err := mongoDriver.Connect(mongoConnectCtx, mongoOpts)
    if err != nil {
        log.WithFields(log.Fields{
            "context": "MongoDB",
            "action":  "CONNECT",
        }).Error(err)
        return err
    }
    mongoPingCtx, _ := context.WithTimeout(context.Background(), 5*time.Second)
    if err := mongoClient.Ping(mongoPingCtx, readpref.Primary()); err != nil {
        log.WithFields(log.Fields{
            "context": "MongoDB",
            "action":  "PING",
        }).Error(err)
        return err
    }
    log.WithFields(log.Fields{
        "context": "MongoDB",
        "status":  "CONNECTED",
    }).Info()
    defer func() {
        if err := mongoClient.Disconnect(context.Background()); err != nil {
            log.WithFields(log.Fields{
                "context": "MongoDB",
                "action":  "DISCONNECT",
            }).Error(err)
        }
    }()
    var updatesLogger tglog.UpdatesLogger

    updatesLogger = &dblog.DBLogger{
        UpdatesRepo: mongo.NewUpdatesRepository(mongoClient.Database(viper.GetString(logDBNameKey))),
    }

    poller := &tb.LongPoller{
        Timeout: 60 * time.Second,
    }
    logPoller := tb.NewMiddlewarePoller(poller, func(update *tb.Update) bool {
        go func() {
            upd := repo.FromTelebotUpdate(update)
            if err := updatesLogger.LogUpdates([]repo.Update{*upd}); err != nil {
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
        Client: &http.Client{
            Timeout: 30 * time.Second,
        },
        Reporter: func(err error) {
            notifierLogEntry.Error(err)
        },
    })
    if err != nil {
        botLogEntry.WithField("action", "AUTHENTICATE").Error(err)
        return err
    }
    botLogEntry.WithFields(log.Fields{
        "status":  "AUTHENTICATED",
        "account": b.Me.Username,
    }).Info()

    if _, err := bot.Configure(b); err != nil {
        botLogEntry.WithField("action", "CONFIGURE").Error(err)
        return err
    }
    statusNotifier := notifier.NewBaseNotifier()

    notifierToken := viper.GetString(notifierBotTokenKey)
    if notifierToken != "" {
        nb, err := tb.NewBot(tb.Settings{
            URL:   tgEndpoint,
            Token: notifierToken,
        })
        if err != nil {
            notifierLogEntry.WithField("action", "AUTHENTICATE").Error(err)
            return err
        }
        notifierLogEntry.WithFields(log.Fields{
            "status":  "AUTHENTICATED",
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
        botLogEntry.WithField("status", "STARTED").Info()
        b.Start()
    }()
    defer func() {
        b.Stop()
        botLogEntry.WithFields(log.Fields{
            "status":       "STOPPED",
            "lastUpdateID": poller.LastUpdateID,
        }).Info()
    }()

    statusNotifier.NotifyUp(b.Me.Username)
    defer statusNotifier.NotifyDown(b.Me.Username)

    gotSig := <-sigCh
    log.WithFields(log.Fields{
        "signal": gotSig.String(),
    }).Info("Got signal, stopping...")

    return nil
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        log.Fatal(err)
    }
}

func init() {
    cobra.OnInitialize(initConfig)

    rootCmd.PersistentFlags().Bool(debugKey, false, "Debug mode")
    _ = viper.BindPFlag(debugKey, rootCmd.PersistentFlags().Lookup(debugKey))
}

func initConfig() {
    viper.SetEnvPrefix("TELEGO")
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
    viper.AutomaticEnv()
}

func checkMandatoryParams() {
    mandatoryParams := [...]string{
        botTokenKey,
        logDBHostKey,
        logDBPortKey,
        logDBUserKey,
        logDBPasswordKey,
        logDBNameKey,
    }
    var missing []string

    for _, k := range mandatoryParams {
        if !viper.IsSet(k) {
            missing = append(missing, k)
        }
    }

    if len(missing) > 0 {
        log.Fatalf("missing: %s", strings.Join(missing, ", "))
    }
}

func checkDependentParams() {
    if viper.IsSet(notifierBotTokenKey) != viper.IsSet(notifierBotChatIDKey) {
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
