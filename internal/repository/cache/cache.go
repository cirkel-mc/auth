package cache

import "github.com/comrades-mc/goutils/config/database/rdc"

const (
	prefixAccessToken  = "at:%s"
	prefixRefreshToken = "rt:%s"
)

type cacheRepository struct {
	client rdc.Rdc
}

func New(c rdc.Rdc) *cacheRepository {
	return &cacheRepository{client: c}
}
