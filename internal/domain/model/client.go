package model

type Client struct {
	BaseModel
	Name         string `db:"name"`
	ClientId     string `db:"client_id"`
	ClientSecret string `db:"client_secret"`
	PublicKey    string `db:"public_key"`
	Channel      string `db:"channel"`
}
