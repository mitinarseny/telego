package bot

import (
    "context"

    "github.com/mitinarseny/telego/administration/repo"
    "github.com/mitinarseny/telego/bot/filters"
    "github.com/mitinarseny/telego/bot/handlers"
    tb "gopkg.in/tucnak/telebot.v2"
)

func Configure(b *tb.Bot, storage *Storage, superUserID int64) (*tb.Bot, error) {
    h := handlers.Handler{
        Bot: b,
        Storage: &handlers.Storage{
            Admins: storage.Admins,
            Roles:  storage.Roles,
        },
    }
    if _, err := storage.Admins.CreateIfNotExists(context.Background(), &repo.Admin{
        ID: superUserID,
        Role: &repo.Role{
            Name: repo.SuperUserRoleName,
        },
    }); err != nil {
        return nil, err
    }
    f := filters.NewFilters(&filters.Storage{
        Admins: storage.Admins,
        Roles:  storage.Roles,
    })
    b.Handle("/hello", f.WithLog(h.HandleHello))
    b.Handle("/stats", f.WithLog(f.OnlyFromSuperUsers(h.HandleStats)))
    b.Handle("/admins", f.WithLog(f.OnlyFromSuperUsers(h.HandleAdmins)))
    b.Handle("/addadmin", f.WithLog(h.HandleAddAdmin))
    return b, nil
}
