package storage

import (
	"github.com/berkantay/user-management-service/model"
)

type UserRepository interface {
	AddUser(userInfo *model.UserInfo)
}
