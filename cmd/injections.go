package cmd

import (
	"context"
	"fmt"

	"github.com/comrades-mc/goutils/abstract"
	"github.com/comrades-mc/goutils/config"
)

func injectDependencies(cfg config.Config) (deps abstract.Dependency) {
	cfg.Injections(func(ctx context.Context) []abstract.Closer {
		master, slave := sqlDatabase()
		fmt.Println("sebelum redis")
		redis := redisDatabase(cfg.GetServiceName())

		deps = abstract.New(
			abstract.SetSQLDatabase(abstract.Master, master),
			abstract.SetSQLDatabase(abstract.Slave, slave),
			abstract.SetRedisDatabase(redis),
		)

		return []abstract.Closer{
			master, slave, redis,
		}
	})

	return
}
