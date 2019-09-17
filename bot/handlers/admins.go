package handlers

import (
    "context"
    "fmt"
    "strconv"

    "github.com/mitinarseny/telego/administration/repo"
    "github.com/mitinarseny/telego/bot/filters"
    "github.com/pkg/errors"
    tb "gopkg.in/tucnak/telebot.v2"
)

type AdminsStorage struct {
    Admins repo.AdminsRepo
}

type Admins struct {
    Logger
    B       *tb.Bot
    Storage *AdminsStorage
}

func (h *Admins) HandleMsg(m *tb.Message) error {
    adms, err := h.Storage.Admins.GetAll(context.Background())
    if err != nil {
        return err
    }
    inlineKeys, err := h.makeAdminsBtns(adms)
    if err != nil {
        return err
    }
    _, err = h.B.Send(m.Sender, "Here is the list of admins:", &tb.ReplyMarkup{
        InlineKeyboard: inlineKeys,
    })
    return err
}

func (h *Admins) makeAdminsBtns(admins []*repo.Admin) ([][]tb.InlineButton, error) {
    keyboard := make([][]tb.InlineButton, 0, len(admins)/2+len(admins)%2)
    for i, admin := range admins {
        if admin.Role == nil {
            return nil, errors.Errorf("admin %d has no role", admin.ID)
        }
        if i%2 == 0 {
            keyboard = append(keyboard, make([]tb.InlineButton, 0, 2))
        }
        strAdminID := strconv.FormatInt(admin.ID, 10)
        chat, err := h.B.ChatByID(strAdminID)
        if err != nil {
            return nil, errors.Wrapf(err, "can not get chat with %q", strAdminID)
        }
        adminName := chat.FirstName
        if chat.LastName != "" {
            adminName += " " + chat.LastName
        }
        adminName += " (" + admin.Role.Name + ")"
        btn := tb.InlineButton{
            Unique: "adminChosen",
            Text:   adminName,
            Data:   strAdminID,
        }
        h.B.Handle(&btn, CallbackWithLog(h, WithCallbackFilters(&adminChosen{
            b: h.B,
            storage: &ChosenAdminStorage{
                Admins: h.Storage.Admins,
            },
        }, filters.WithSender().IsAdminWithScopes(h.Storage.Admins, repo.AdminsReadScope))))
        keyboard[i/2] = append(keyboard[i/2], btn)
    }
    return keyboard, nil
}

type ChosenAdminStorage struct {
    Admins repo.AdminsRepo
}

type adminChosen struct {
    b       *tb.Bot
    storage *ChosenAdminStorage
}

func (h *adminChosen) HandleCallback(c *tb.Callback) error {
    adminID, err := strconv.ParseInt(c.Data, 10, 64)
    if err != nil {
        return errors.Wrap(err, "can not parse adminID")
    }
    admin, err := h.storage.Admins.GetByID(context.Background(), adminID)
    if err != nil {
        return err
    }
    text, err := h.adminDescription(admin)
    if err != nil {
        return err
    }
    _, err = h.b.Edit(c.Message, text, tb.ModeMarkdown)
    return err
}

const (
    AdminTemplate = `*Admin:*
*ID:* %s
*User:* %s`
)

func (h *adminChosen) adminDescription(admin *repo.Admin) (string, error) {
    strAdminID := strconv.FormatInt(admin.ID, 10)
    chat, err := h.b.ChatByID(strAdminID)
    if err != nil {
        return "", errors.Wrapf(err, "can not get chat with %q", strAdminID)
    }
    name := "[" + chat.FirstName
    if chat.LastName != "" {
        name += " " + chat.LastName
    }
    name += "](tg://user?id=" + strAdminID + ")"
    return fmt.Sprintf(AdminTemplate, strAdminID, name), nil
}
