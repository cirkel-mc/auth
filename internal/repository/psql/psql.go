package psql

import "github.com/comrades-mc/goutils/config/database/dbc"

type psqlRepository struct {
	master, slave dbc.SqlDbc
}

func New(m, s dbc.SqlDbc) *psqlRepository {
	return &psqlRepository{master: m, slave: s}
}
