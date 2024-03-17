package usecase

import (
	"cirkel/auth/internal/domain/dto"
	"cirkel/auth/internal/repository"
	"context"
	"time"

	"github.com/cirkel-mc/goutils/zone"
)

type usecaseInstance struct {
	tz    *time.Location
	psql  repository.Psql
	cache repository.Cache
}

type Usecase interface {
	Register(ctx context.Context, req *dto.RequestRegister) (*dto.Token, error)
	Login(ctx context.Context, req *dto.RequestLogin) (resp *dto.Token, err error)
	ValidateAuth(ctx context.Context, req *dto.RequestHeader) (resp *dto.Auth, err error)
	ValidateCsrfToken(ctx context.Context, csrfToken string) (string, error)
}

func New(p repository.Psql, c repository.Cache) *usecaseInstance {
	return &usecaseInstance{
		tz:    zone.TzJakarta(),
		psql:  p,
		cache: c,
	}
}
