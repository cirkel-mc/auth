package model

type Role struct {
	BaseModel
	Name string `db:"name"`
	Key  string `db:"key"`
}
