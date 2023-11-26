package usecase

import (
	"comrades-mc/auth/internal/domain/dto"
	"comrades-mc/auth/internal/repository"
	"context"
	"time"

	"github.com/comrades-mc/goutils/zone"
)

type usecaseInstance struct {
	tz    *time.Location
	psql  repository.Psql
	cache repository.Cache
}

type Usecase interface {
	Register(ctx context.Context, req *dto.RequestRegister) (*dto.Token, error)
}

func New(p repository.Psql, c repository.Cache) *usecaseInstance {
	return &usecaseInstance{
		tz:    zone.TzJakarta(),
		psql:  p,
		cache: c,
	}
}
