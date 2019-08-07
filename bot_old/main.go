package bot_old

import (
    "context"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "time"

    "github.com/mitinarseny/telego/bot"
    "github.com/mitinarseny/telego/bot_old/ch_log"
    log "github.com/sirupsen/logrus"
    "github.com/spf13/viper"

    "github.com/go-telegram-bot-api/telegram-bot-api"
    "github.com/mitinarseny/telego/bot_old/handlers"
)

const (
    logUpdatesBufferSize = 10000
    logUpdatesTimeDelta  = 10 * time.Second
)

func startHandlingUpdates(ctx context.Context, botAPI *tgbotapi.BotAPI) (<-chan error, error) {
    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates, err := botAPI.GetUpdatesChan(u)
    if err != nil {
        return nil, err
    }




    errCh := make(chan error)
    defer close(errCh)

    go func() {
        bot := handlers.Bot{
            BotAPI: botAPI,
        }
        ul, err := ch_log.NewUpdatesLogger(
            viper.GetString("log.db.host"),
            viper.GetUint("log.db.port"),
            viper.GetString("log.db.user"),
            viper.GetString("log.db.password"),
            viper.GetString("log.db.name"),
            viper.GetString("log.db.table"))
        if err != nil {
            log.WithFields(log.Fields{
                "context": "LOG",
            }).Error(err)
        }
        if err := bot.HandleUpdates(updates, ul, errCh); err != nil {
            select {
            case errCh <- err:
            case <-ctx.Done():
                return
            }
        }
    }()
    return errCh, nil
}

func Run(token, notifierToken string, notifyChatID int64, debug bool) error {
    log.Info("Starting...")
    botAPI, err := tgbotapi.NewBotAPI(token)
    if err != nil {
        return err
    }
    botAPI.Debug = debug
    log.WithField("bot", botAPI.Self.UserName).Info("Authorized")

    ctx, cancelFunc := context.WithCancel(context.Background())
    defer cancelFunc()

    if notifierToken != "" {
        notifier, err := tgbotapi.NewBotAPI(notifierToken)
        if err != nil {
            return err
        }
        log.WithField("notifier", notifier.Self.UserName).Info()
        _ = bot.notifyUp(notifier, notifyChatID, botAPI.Self.UserName)
        defer bot.notifyDown(notifier, notifyChatID, botAPI.Self.UserName)
    }

    sigErrCh := getSignalErrorCh(ctx)

    updErrCh, err := startHandlingUpdates(ctx, botAPI)
    if err != nil {
        return err
    }

    return waitForSigOrError(ctx, updErrCh, sigErrCh)
}

type SignalError struct{}

func (S SignalError) Error() string {
    return "signal error"
}

func waitForSigOrError(ctx context.Context, errChs ...<-chan error) error {
    ctx, cancelFunc := context.WithCancel(ctx)
    defer cancelFunc()

    errCh := mergeErrorChs(ctx, errChs...)
    for err := range errCh {
        if err != nil {
            switch err.(type) {
            case SignalError:
                return nil
            default:
                return err
            }
        }
    }
    return nil
}

func getSignalErrorCh(ctx context.Context) <-chan error {
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

    sigErrCh := make(chan error, 1)
    go func() {
        for range sigCh {
            sigErrCh <- SignalError{}
        }
    }()
    return sigErrCh
}

func mergeErrorChs(ctx context.Context, cs ...<-chan error) <-chan error {
    var wg sync.WaitGroup
    out := make(chan error)

    output := func(c <-chan error) {
        defer wg.Done()
        for n := range c {
            select {
            case out <- n:
            case <-ctx.Done():
                return
            }
        }
    }

    wg.Add(len(cs))
    for _, c := range cs {
        go output(c)
    }

    go func() {
        defer close(out)
        wg.Wait()
    }()
    return out
}
