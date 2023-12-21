package model

import (
	"time"

	"github.com/cirkel-mc/goutils/null"
)

type UserDevice struct {
	Id        int         `db:"id"`
	CreatedAt time.Time   `db:"created_at"`
	UserId    int         `db:"user_id"`
	DeviceId  string      `db:"device_id"`
	Channel   string      `db:"channel"`
	UserAgent string      `db:"user_agent"`
	FcmToken  null.String `db:"fcm_token"`
}
