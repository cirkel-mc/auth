package cache

import "github.com/cirkel-mc/goutils/config/database/rdc"

const (
	prefixAccessToken  = "at:%s"
	prefixRefreshToken = "rt:%s"
	prefixCsrfToken    = "csrf:%s"
)

type cacheRepository struct {
	client rdc.Rdc
}

func New(c rdc.Rdc) *cacheRepository {
	return &cacheRepository{client: c}
}
