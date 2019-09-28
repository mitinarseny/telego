package cmd

import (
    "context"
    "os"
    "os/signal"
    "strings"
    "syscall"
    "time"

    "github.com/mitinarseny/telego/admins"
    mongoadmin "github.com/mitinarseny/telego/admins/mongo"
    "github.com/mitinarseny/telego/bot"
    "github.com/mitinarseny/telego/bot/tglog"
    mongotglog "github.com/mitinarseny/telego/bot/tglog/mongo"
    "github.com/mitinarseny/telego/log"
    mongolog "github.com/mitinarseny/telego/log/mongo"
    "github.com/mitinarseny/telego/log/stderr"
    "github.com/pkg/errors"
    "github.com/sirupsen/logrus"
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

var logger = logrus.New()

var rootCmd = &cobra.Command{
    Run: func(cmd *cobra.Command, args []string) {
        if err := checkMandatoryParams(); err != nil {
            logger.Fatal(err)
        }
        if err := start(); err != nil {
            logger.Fatal(err)
        }
    },
}

func start() error {
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
        return errors.Wrap(err, "unable to connect to mongo")
    }

    defer func() {
        if err := mongoClient.Disconnect(context.Background()); err != nil {
            logger.WithFields(logrus.Fields{
                "context": "MongoDB",
                "action":  "DISCONNECT",
            }).Error(err)
        }
    }()

    mongoPingCtx, dropMongoPing := context.WithTimeout(context.Background(), 5*time.Second)
    defer dropMongoPing()

    if err := mongoClient.Ping(mongoPingCtx, readpref.Primary()); err != nil {
        return errors.Wrap(err, "error while pining mongo")
    }
    logger.WithFields(logrus.Fields{
        "context": "MongoDB",
        "status":  "CONNECTED",
    }).Info()

    botMongoDB := mongoClient.Database(viper.GetString(dbNameKey))
    mongos, err := createMongoRepos(botMongoDB)
    if err != nil {
        return errors.Wrap(err, "unable to create mongo repositories")
    }
    botStdInfoErrorLogger := stderr.NewErrorLogger(logger)
    botPrefs := bot.Settings{
        Token:        viper.GetString(botTokenKey),
        LastUpdateID: 0, // TODO: set from env
        Storage: &bot.Storage{
            Admins: mongos.Admins,
            Roles:  mongos.Roles,
        },
        UpdateLogger: tglog.NewUpdatesLogger(mongos.TgUpdates),
        Logger: log.Unsafe(log.Multi(
            botStdInfoErrorLogger,
            log.NewPropagateInfoError(
                mongos.InfoErrorLogs,
                botStdInfoErrorLogger,
            ),
        )),
        SuperuserID: viper.GetInt64(superuserIDKey),
    }

    b, err := bot.New(&botPrefs)
    if err != nil {
        return errors.Wrap(err, "unable to create bot")
    }

    go b.Start()
    defer b.Stop()

    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

    gotSig := <-sigCh
    logger.WithFields(logrus.Fields{
        "signal": gotSig.String(),
        "status": "STOPPING",
    }).Info()

    return nil
}

type Repositories struct {
    Admins        admins.AdminsRepo
    Roles         admins.RolesRepo
    TgUsers       tglog.UsersRepo
    TgChats       tglog.ChatsRepo
    TgUpdates     tglog.UpdatesRepo
    InfoErrorLogs log.InfoErrorLogger
}

func createMongoRepos(db *mongoDriver.Database) (*Repositories, error) {
    botRoles, err := mongoadmin.NewRolesRepo(db)
    if err != nil {
        return nil, errors.Wrap(err, "unable to create mongoadmin.RolesRepo")
    }

    botAdmins, err := mongoadmin.NewAdminsRepo(db, botRoles)
    if err != nil {
        return nil, errors.Wrap(err, "unable to create mongoadmin.AdminsRepo")
    }

    botTgUsers, err := mongotglog.NewUsersRepo(db)
    if err != nil {
        return nil, errors.Wrap(err, "unable to create mongotglog.UsersRepo")
    }

    botTgChats, err := mongotglog.NewChatsRepo(db)
    if err != nil {
        return nil, errors.Wrap(err, "unable to create mongotglog.ChatsRepo")
    }

    botTgUpdates := mongotglog.NewUpdatesRepo(db, botTgUsers, botTgChats)

    botDBInfoErrorLogger := mongolog.NewErrorLogger(db)
    return &Repositories{
        Admins:        botAdmins,
        Roles:         botRoles,
        TgUsers:       botTgUsers,
        TgChats:       botTgChats,
        TgUpdates:     botTgUpdates,
        InfoErrorLogs: botDBInfoErrorLogger,
    }, nil
}

func checkMandatoryParams() error {
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
        return errors.Errorf("missing: %s", strings.Join(missing, ", "))
    }
    return nil
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        logger.Fatal(err)
    }
}

func initConfig() {
    viper.SetEnvPrefix("TELEGO")
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
    viper.AutomaticEnv()

    if viper.GetBool(debugKey) {
        logger.SetLevel(logrus.DebugLevel)
        logger.SetReportCaller(true)
    }
}

func init() {
    cobra.OnInitialize(initConfig)

    rootCmd.PersistentFlags().Bool(debugKey, false, "Debug mode")
    _ = viper.BindPFlag(debugKey, rootCmd.PersistentFlags().Lookup(debugKey))
}
