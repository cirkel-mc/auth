package cmd

import (
	"time"

	"github.com/comrades-mc/goutils/abstract"
	"github.com/comrades-mc/goutils/config/database"
	"github.com/comrades-mc/goutils/env"
)

func sqlDatabase() (abstract.SQLDatabase, abstract.SQLDatabase) {
	opts := []database.SqlFuncOption{
		database.SetSqlMaxConnection(20),
		database.SetSqlMaxIdleConnection(0),
		database.SetSqlMaxIdleTime(10 * time.Second),
	}

	master, err := database.NewSqlxConnection(append(opts, database.SetSqlDSN(env.GetString("DSN_MASTER")))...)
	if err != nil {
		panic(err)
	}

	slave, err := database.NewSqlxConnection(append(opts, database.SetSqlDSN(env.GetString("DSN_SLAVE")))...)
	if err != nil {
		panic(err)
	}

	return master, slave
}

func redisDatabase(svcName string) abstract.RedisDatabase {
	return database.NewRedisConnection(
		database.SetRedisServiceName(svcName),
		database.SetRedisAddress(env.GetListString("REDIS")),
		database.SetRedisMaxIdleConnection(0),
		database.SetRedisMaxIdleTimeout(10*time.Second),
		database.SetRedisMinIdleConnection(0),
		database.SetRedisPoolSize(10),
	)
}
