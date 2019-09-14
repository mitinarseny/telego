package mongo

import (
    "context"
    "errors"
    "fmt"
    "time"

    tb "github.com/charithe/timedbuf"
    "github.com/mitinarseny/telego/tglog/repo"
    log "github.com/sirupsen/logrus"
    "go.mongodb.org/mongo-driver/mongo"
)

const (
    updatesCollectionName = "updates"
)

type UpdatesRepo struct {
    this  *mongo.Collection
    users repo.UsersRepo
    chats repo.ChatsRepo
}

type UpdatesRepoDependentRepos struct {
    Users repo.UsersRepo
    Chats repo.ChatsRepo
}

func (d *UpdatesRepoDependentRepos) Validate() error {
    shouldBeNotNil := [...]interface{}{
        d.Chats,
        d.Users,
    }
    for _, e := range shouldBeNotNil {
        if e == nil {
            return errors.New(fmt.Sprintf("%T should be not nil", e))
        }
    }
    return nil
}

func NewUpdatesRepo(db *mongo.Database, deps *UpdatesRepoDependentRepos) (*UpdatesRepo, error) {
    if err := deps.Validate(); err != nil {
        return nil, err
    }
    return &UpdatesRepo{
        this:  db.Collection(updatesCollectionName),
        users: deps.Users,
        chats: deps.Chats,
    }, nil
}

func (r *UpdatesRepo) Create(ctx context.Context, updates ...*repo.Update) ([]*repo.Update, error) {
    upds := make([]interface{}, 0, len(updates))
    for _, update := range updates {
        if msg := update.Message; msg != nil {
            if from := msg.From; from != nil {
                if _, err := r.users.CreateIfNotExists(ctx, from); err != nil {
                    return nil, err
                }
            }
            if chat := msg.Chat; chat != nil {
                if _, err := r.chats.CreateIfNotExists(ctx, chat); err != nil {
                    return nil, err
                }
            }
        }
        ca := time.Now()
        update.CreatedAt = &ca
        upds = append(upds, update)
    }
    // _, err := r.this.UpdateMany(ctx, )
    _, err := r.this.InsertMany(ctx, upds)
    if err != nil {
        log.WithFields(log.Fields{
            "context": "UpdatesRepo",
            "action":  "CREATE",
        }).Error(err)
        return nil, err
    }
    log.WithFields(log.Fields{
        "context": "UpdatesRepo",
        "status":  "CREATED",
        "count":   len(updates),
    }).Info()
    return updates, nil
}

const (
    buffSize     = 10000
    buffDuration = 5 * time.Second // TODO: minute?
)

type BufferedUpdatesRepo struct {
    *UpdatesRepo
    tb *tb.TimedBuf
}

func NewBufferedUpdatesRepo(ur *UpdatesRepo) *BufferedUpdatesRepo {
    r := &BufferedUpdatesRepo{
        UpdatesRepo: ur,
    }
    r.tb = tb.New(buffSize, buffDuration, func(items []interface{}) {
        upds := make([]*repo.Update, 0, len(items))
        for _, item := range items {
            switch i := item.(type) {
            case *repo.Update:
                upds = append(upds, i)
            case repo.Update:
                upds = append(upds, &i)
            default:
                log.WithFields(log.Fields{
                    "context": "BufferedUpdatesRepo",
                    "action":  "FLUSH",
                }).Error(errors.New(fmt.Sprintf(
                    "can not flush %T",
                    i)))
            }
        }
        _, err := r.UpdatesRepo.Create(context.Background(), upds...)
        if err != nil {
            return
        }
        return
    })
    return r
}

func (r *BufferedUpdatesRepo) Create(_ context.Context, updates ...*repo.Update) ([]*repo.Update, error) {
    upds := make([]interface{}, 0, len(updates))
    for _, u := range updates {
        upds = append(upds, u)
    }
    r.tb.Put(upds...)
    return nil, nil
}

func (r *BufferedUpdatesRepo) Close() error {
    r.tb.Close()
    return nil
}
