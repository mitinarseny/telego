package admins

import (
    "context"
    "fmt"
    "strconv"
)

type NotifierType string

const (
    TelegramNotifier NotifierType = "telegram"
)

type NotificationType string

const (
    StatusNotificationType NotificationType = "status"
)

type Notifications struct {
    Status []NotifierType `bson:"status,omitempty"`
}

type Admin struct {
    ID            int64         `bson:"id, omitempty"`
    Role          *Role         `bson:"role,omitempty"`
    Notifications *Notifications `bson:"notifications,omitempty"`
}

func (a *Admin) HadScopes(scopes ...Scope) bool {
    return a.Role.HasScopes(scopes...)
}

func (a *Admin) Recipient() string {
    return strconv.FormatInt(a.ID, 10)
}

type AdminsRepo interface {
    Create(ctx context.Context, admins ...*Admin) ([]*Admin, error)
    CreateIfNotExists(ctx context.Context, admins ...*Admin) ([]*Admin, error)
    AssignRoleByID(ctx context.Context, roleName string, adminID int64) (*Admin, error)
    AssignRoleByIDs(ctx context.Context, roleName string, adminIDs ...int64) ([]*Admin, error)
    GetAll(ctx context.Context) ([]*Admin, error)
    GetAllShouldBeNotifiedAbout(ctx context.Context, notificationType NotificationType) ([]*Admin, error)
    GetByID(ctx context.Context, adminID int64) (*Admin, error)
    GetByIDs(ctx context.Context, adminIDs ...int64) ([]*Admin, error)
    GetRoleByID(ctx context.Context, adminID int64) (*Role, error)
    HasScopesByID(ctx context.Context, adminID int64, scopes ...Scope) (bool, error)
    DeleteByIDs(ctx context.Context, adminIDs ...int64) error
}

type AdminNotFound int64

func (e AdminNotFound) Error() string {
    return fmt.Sprintf("Admin %q not found", strconv.FormatInt(int64(e), 10))
}
