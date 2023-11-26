package model

import (
	"comrades-mc/auth/internal/domain/constant"

	"github.com/comrades-mc/goutils/null"
)

type User struct {
	BaseModel
	Username   string              `db:"username"`
	Password   string              `db:"password"`
	Email      string              `db:"email"`
	Status     constant.UserStatus `db:"status"`
	VerifiedAt null.Time           `db:"verified_at"`

	// foreign key
	Role *Role `db:"ro"`
}
