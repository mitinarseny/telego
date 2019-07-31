package bot

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mitinarseny/telego/bot/handlers"
)

func startHandlingUpdates(ctx context.Context, bot *tgbotapi.BotAPI) (<-chan error, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return nil, err
	}

	errCh := make(chan error)
	go func() {
		defer close(errCh)
		for update := range updates {
			select {
			case <-ctx.Done():
				return
			default:
			}
			go func() {
				switch {
				case update.Message != nil:
					switch {
					case update.Message.Command() == "hello":
						if err := handlers.HandleHello(bot, update); err != nil {
							errCh <- err
						}
					default:
						if err := handlers.HandleUnsupported(bot, update); err != nil {
							errCh <- err
						}
					}
				default:
					if err := handlers.HandleUnsupported(bot, update); err != nil {
						errCh <- err
					}
				}
			}()
		}
	}()
	return errCh, nil
}

func Run(token, notifierToken string, notifyChatID int64, debug bool) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}
	bot.Debug = debug
	log.Printf("Authorized on account: @%s", bot.Self.UserName)

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	if notifierToken != "" {
		notifier, err := tgbotapi.NewBotAPI(notifierToken)
		if err != nil {
			return err
		}
		log.Printf("Notifier: @%s", notifier.Self.UserName)
		_ = notifyUp(notifier, notifyChatID, bot.Self.UserName)
		defer notifyDown(notifier, notifyChatID, bot.Self.UserName)
	}

	sigErrCh := getSignalErrorCh(ctx)

	updErrCh, err := startHandlingUpdates(ctx, bot)
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
	signal.Notify(sigCh, os.Interrupt, os.Kill)

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
