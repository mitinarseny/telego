package bot

import "github.com/mitinarseny/telego/administration/repo"

type Storage struct {
    Admins repo.AdminsRepo
    Roles repo.RolesRepo
}
