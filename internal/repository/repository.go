package repository

import (
	"cirkel/auth/internal/domain/dto"
	"cirkel/auth/internal/domain/model"
	"context"

	"github.com/cirkel-mc/goutils/config/database/dbc"
	"github.com/cirkel-mc/goutils/types"
)

type Cache interface {
	GetAccessToken(ctx context.Context, accessToken string) (*types.TokenClaim, error)
	GetRefreshToken(ctx context.Context, refreshToken string) (*types.TokenClaim, error)
	SetAccessToken(ctx context.Context, tc *types.TokenClaim) (*dto.Token, error)
	GetCsrfToken(ctx context.Context, csrfToken string) error
	SetCsrfToken(ctx context.Context) (string, error)
	DeleteCsrfToken(ctx context.Context, csrfToken string) error
}

type Psql interface {
	StartTransaction(ctx context.Context, txFunc func(context.Context, dbc.SqlDbc) error) error
	GetUserNextVal(ctx context.Context) (int, error)
	FindUserById(ctx context.Context, id int) (resp *model.User, err error)
	FindUserByUsername(ctx context.Context, username string) (*model.User, error)
	FindUserByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, user *model.User) error
	CreateUserDevice(ctx context.Context, ud *model.UserDevice) error
	FindRoleById(ctx context.Context, id int) (resp *model.Role, err error)
	FindClientByClientId(ctx context.Context, clientId string) (resp *model.Client, err error)
}
