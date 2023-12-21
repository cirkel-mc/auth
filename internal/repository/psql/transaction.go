package psql

import (
	"context"

	"github.com/cirkel-mc/goutils/config/database/dbc"
)

func (p *psqlRepository) StartTransaction(ctx context.Context, txFunc func(context.Context, dbc.SqlDbc) error) error {
	return p.master.StartTransaction(ctx, txFunc)
}
