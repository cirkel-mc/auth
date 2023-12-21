package model

import (
	"time"

	"github.com/cirkel-mc/goutils/null"
)

type BaseModel struct {
	Id        int       `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	CreatedBy null.Int  `db:"created_by"`
	UpdatedAt null.Time `db:"updated_at"`
	UpdatedBy null.Int  `db:"updated_by"`
	DeletedAt null.Time `db:"deleted_at"`
	DeletedBy null.Int  `db:"deleted_by"`
}
