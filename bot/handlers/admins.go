package handlers

import (
    "context"
    "fmt"
    "regexp"
    "strconv"

    "github.com/mitinarseny/telego/admins"
    "github.com/mitinarseny/telego/bot/chattools"
    "github.com/mitinarseny/telego/bot/filters"
    "github.com/mitinarseny/telego/log"
    "github.com/pkg/errors"
    tb "gopkg.in/tucnak/telebot.v2"
)

const (
    adminChosenEvent     = "adminChosen"
    adminRoleChosenEvent = "adminRoleChosen"
)

type AdminsStorage struct {
    Admins admins.AdminsRepo
    Roles  admins.RolesRepo
}

type Admins struct {
    log.UnsafeInfoErrorLogger
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

func (h *Admins) makeAdminsBtns(adms []*admins.Admin) ([][]tb.InlineButton, error) {
    keyboard := make([][]tb.InlineButton, 0, len(adms)/2+len(adms)%2)
    for i, admin := range adms {
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
            Unique: adminChosenEvent,
            Text:   adminName,
            Data:   strAdminID,
        }
        h.B.Handle(&btn, CallbackWithLog(h, WithCallbackFilters(&adminChosen{
            UnsafeInfoErrorLogger: h,
            b:                     h.B,
            storage: &chosenAdminStorage{
                Admins: h.Storage.Admins,
                Roles:  h.Storage.Roles,
            },
        }, filters.WithSender().IsAdminWithScopes(h.Storage.Admins, admins.AdminsReadScope))))
        keyboard[i/2] = append(keyboard[i/2], btn)
    }
    return keyboard, nil
}

type chosenAdminStorage struct {
    Admins admins.AdminsRepo
    Roles  admins.RolesRepo
}

type adminChosen struct {
    log.UnsafeInfoErrorLogger
    b       *tb.Bot
    storage *chosenAdminStorage
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
    sender, err := h.storage.Admins.GetByID(context.Background(), int64(c.Sender.ID))
    if err != nil {
        return err
    }
    return h.handle(c.Message, sender, admin)
}

func (h *adminChosen) handle(m *tb.Message, sender, admin *admins.Admin) error {
    var (
        text string
        opts = []interface{}{
            tb.ModeMarkdown,
        }
    )
    if sender.HadScopes(admins.AdminsScope) {
        txt, err := h.adminDescription(admin)
        if err != nil {
            return err
        }
        text = txt
        btns, err := h.makeAdminBtns(admin)
        if err != nil {
            return err
        }
        opts = append(opts, &tb.ReplyMarkup{
            InlineKeyboard: btns,
        })
    } else {
        txt, err := h.adminDescriptionWithRole(admin)
        if err != nil {
            return err
        }
        text = txt
    }
    _, err := h.b.Edit(m, text, opts...)
    return err
}

func (h *adminChosen) makeAdminBtns(admin *admins.Admin) ([][]tb.InlineButton, error) {
    roles, err := h.storage.Roles.GetAll(context.Background())
    if err != nil {
        return nil, err
    }
    keyboard := make([][]tb.InlineButton, 0, len(roles)/2+len(roles)%2)
    for i, role := range roles {
        if i%2 == 0 {
            keyboard = append(keyboard, make([]tb.InlineButton, 0, 2))
        }
        btn := tb.InlineButton{
            Unique: adminRoleChosenEvent,
            Data:   strconv.FormatInt(admin.ID, 10) + "@" + role.Name,
        }
        if admin.Role.Name == role.Name {
            btn.Text += "✅ "
        }
        btn.Text += role.Name
        h.b.Handle(&btn, CallbackWithLog(h, WithCallbackFilters(&roleChosen{
            b: h.b,
            storage: &roleChosenStorage{
                Admins: h.storage.Admins,
                Roles:  h.storage.Roles,
            },
            parent: h,
        },
            filters.WithSender().IsAdminWithScopes(h.storage.Admins, admins.AdminsScope),
            filters.DataShouldMatch(adminRoleRegex),
        )))
        keyboard[i/2] = append(keyboard[i/2], btn)
    }
    return keyboard, nil
}

const (
    adminTemplate = `*Admin:*
*ID:* %s
*From:* %s`
    adminWithRoleTemplate = adminTemplate + `
*Role:* %s`
)

func (h *adminChosen) adminDescription(admin *admins.Admin) (string, error) {
    strAdminID := strconv.FormatInt(admin.ID, 10)
    userLink, err := chattools.WithBot(h.b).UserLink(strAdminID)
    if err != nil {
        return "", err
    }
    return fmt.Sprintf(adminTemplate, strAdminID, userLink), nil
}

func (h *adminChosen) adminDescriptionWithRole(admin *admins.Admin) (string, error) {
    strAdminID := strconv.FormatInt(admin.ID, 10)
    userLink, err := chattools.WithBot(h.b).UserLink(strAdminID)
    if err != nil {
        return "", err
    }
    var role string
    if admin.Role != nil {
        role = admin.Role.Name
    }
    return fmt.Sprintf(adminWithRoleTemplate, strAdminID, userLink, role), nil
}

type roleChosenStorage struct {
    Admins admins.AdminsRepo
    Roles  admins.RolesRepo
}

type roleChosen struct {
    b       *tb.Bot
    storage *roleChosenStorage
    parent  *adminChosen
}

var (
    adminRoleRegex = regexp.MustCompile(`^(?P<adminID>\w+)@(?P<roleName>\w+)$`)
)

func (h *roleChosen) HandleCallback(c *tb.Callback) error {
    sender, err := h.storage.Admins.GetByID(context.Background(), int64(c.Sender.ID))
    if err != nil {
        return err
    }
    match := adminRoleRegex.FindStringSubmatch(c.Data)
    if len(match) != 3 {
        return errors.New(fmt.Sprintf("callback data %q does not match regexp", c.Data)) // TODO: already filtered with filter?
    }
    adminIDStr, roleName := match[1], match[2]
    adminID, err := strconv.ParseInt(adminIDStr, 10, 64)
    if err != nil {
        return errors.Wrap(err, "can not extract adminID from callback data")
    }
    admin, err := h.storage.Admins.GetByID(context.Background(), adminID)
    if err != nil {
        return err
    }
    strAdminID := strconv.FormatInt(admin.ID, 10)
    adminName, err := chattools.WithBot(h.b).GetTelegramName(strAdminID)
    if err != nil {
        return err
    }
    role, err := h.storage.Roles.GetByName(context.Background(), roleName)
    if err != nil {
        return err
    }
    if admin.Role != nil && admin.Role.Name == role.Name {
        return h.b.Respond(c, &tb.CallbackResponse{
            Text: fmt.Sprintf("%s is %s already", adminName, role.Name),
        })
    }
    admin, err = h.storage.Admins.AssignRoleByID(context.Background(), role.Name, admin.ID)
    if err != nil {
        return err
    }
    if err := h.b.Respond(c, &tb.CallbackResponse{
        Text: fmt.Sprintf("%s was assigned %s", adminName, role.Name),
    }); err != nil {
        return err
    }
    return h.parent.handle(c.Message, sender, admin)
}
