package handlers

import (
    "context"
    "errors"
    "fmt"

    "github.com/mitinarseny/telego/administration/repo"
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
    inlineKeys := h.makeAdminsBtns(adms)
    _, err = h.B.Send(m.Sender, "Here is the list of admins:", &tb.ReplyMarkup{
        InlineKeyboard: inlineKeys,
    })
    return err
}

func (h *Admins) makeAdminsBtns(admins []*repo.Admin) [][]tb.InlineButton {
    keyboard := make([][]tb.InlineButton, 0, len(admins)/2+len(admins)%2)
    for i, admin := range admins {
        if i%2 == 0 {
            keyboard = append(keyboard, make([]tb.InlineButton, 0, 2))
        }
        btn := tb.InlineButton{
            Unique: fmt.Sprintf("%d", admin.ID), // TODO: more unique
            Text:   fmt.Sprintf("%d (%s)", admin.ID, admin.Role.Name),
        }
        h.B.Handle(&btn, CallbackWithLog(h, &ChosenAdmin{
            b: h.B,
            storage: &ChosenAdminStorage{
                Admins: h.Storage.Admins,
            },
        }))
        keyboard[i/2] = append(keyboard[i/2], btn)
    }
    return keyboard
}

type ChosenAdminStorage struct {
    Admins repo.AdminsRepo
}

type ChosenAdmin struct {
    b       *tb.Bot
    storage *ChosenAdminStorage
}

func (h *ChosenAdmin) HandleCallback(c *tb.Callback) error {
    return errors.New("ChosenAdmin.HandleCallback is not implemented yet")
}
