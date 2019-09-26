package cmd

import (
    "context"
    "os"
    "os/signal"
    "strings"
    "syscall"
    "time"

    "github.com/mitinarseny/telego/bot"
    "github.com/mitinarseny/telego/bot/logging/errors"
    mongolog "github.com/mitinarseny/telego/bot/logging/errors/mongo"
    "github.com/mitinarseny/telego/bot/logging/errors/stderr"
    repolog "github.com/mitinarseny/telego/bot/logging/updates/repo"
    mongoadmin "github.com/mitinarseny/telego/bot/repo/administration/mongo"
    mongotg "github.com/mitinarseny/telego/bot/repo/tg/mongo"
    log "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    mongoDriver "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
    debugKey = "debug"

    botTokenKey = "bot.token"

    dbHostKey     = "db.host"
    dbPortKey     = "db.port"
    dbUserKey     = "db.user"
    dbPasswordKey = "db.password"
    dbNameKey     = "db.name"

    superuserIDKey = "superuser.id"
)

var rootCmd = &cobra.Command{
    Run: func(cmd *cobra.Command, args []string) {
        if viper.GetBool(debugKey) {
            log.SetLevel(log.DebugLevel)
            log.SetReportCaller(true)
        }
        checkMandatoryParams()
        if err := start(); err != nil {
            log.Fatal(err)
        }
    },
}

func start() error {
    logger := log.New()
    mongoOpts := options.Client().SetAppName("bot").SetAuth(options.Credential{
        Username: viper.GetString(dbUserKey),
        Password: viper.GetString(dbPasswordKey),
    }).SetHosts([]string{
        viper.GetString(dbHostKey),
    })

    mongoConnectCtx, dropMongoConnect := context.WithTimeout(context.Background(), 10*time.Second)
    defer dropMongoConnect()

    mongoClient, err := mongoDriver.Connect(mongoConnectCtx, mongoOpts)
    if err != nil {
        logger.WithFields(log.Fields{
            "context": "MongoDB",
            "action":  "CONNECT",
        }).Error(err)
        return err
    }

    defer func() {
        if err := mongoClient.Disconnect(context.Background()); err != nil {
            logger.WithFields(log.Fields{
                "context": "MongoDB",
                "action":  "DISCONNECT",
            }).Error(err)
        }
    }()

    mongoPingCtx, dropMongoPing := context.WithTimeout(context.Background(), 5*time.Second)
    defer dropMongoPing()

    if err := mongoClient.Ping(mongoPingCtx, readpref.Primary()); err != nil {
        logger.WithFields(log.Fields{
            "context": "MongoDB",
            "action":  "PING",
        }).Error(err)
        return err
    }
    logger.WithFields(log.Fields{
        "context": "MongoDB",
        "status":  "CONNECTED",
    }).Info()

    botMongoDB := mongoClient.Database(viper.GetString(dbNameKey))

    botRoles := mongoadmin.NewRolesRepo(botMongoDB)
    botAdmins := mongoadmin.NewAdminsRepo(botMongoDB, botRoles)

    botTgUsers := mongotg.NewUsersRepo(botMongoDB)
    botTgChats := mongotg.NewChatsRepo(botMongoDB)
    botTgUpdates := mongotg.NewUpdatesRepo(botMongoDB, botTgUsers, botTgChats)

    botStdErrorLogger := stderr.NewErrorLogger(logger)
    botDBErrorLogger := mongolog.NewErrorLogger(botMongoDB, botStdErrorLogger)

    botPrefs := bot.Settings{
        Token:        viper.GetString(botTokenKey),
        LastUpdateID: 0, // TODO: set from env
        Storage: &bot.Storage{
            Admins: botAdmins,
            Roles:  botRoles,
        },
        UpdateLogger: repolog.NewUpdatesLogger(botTgUpdates),
        ErrorLogger: errors.NewCompositeErrorLogger(
            botStdErrorLogger,
            botDBErrorLogger,
        ),
        SuperuserID: viper.GetInt64(superuserIDKey),
    }

    b, err := bot.New(&botPrefs)
    if err != nil {
        logger.WithFields(log.Fields{
            "context": "BOT",
            "action":  "CREATE",
        }).Error(err)
        return err
    }

    go b.Start()
    defer b.Stop()

    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

    gotSig := <-sigCh
    log.WithFields(log.Fields{
        "signal": gotSig.String(),
        "status": "STOPPING",
    }).Info()

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
        dbHostKey,
        dbPortKey,
        dbUserKey,
        dbPasswordKey,
        dbNameKey,
        superuserIDKey,
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
