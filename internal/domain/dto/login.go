package dto

type RequestLogin struct {
	*RequestHeader
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}
