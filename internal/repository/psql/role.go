package psql

import (
	"cirkel/auth/internal/domain/model"
	"context"
	"fmt"

	"github.com/cirkel-mc/goutils/logger"
	"github.com/cirkel-mc/goutils/tracer"
)

const queryRole = `
select id, name, key
from "user".roles
`

func (p *psqlRepository) FindRoleById(ctx context.Context, id int) (resp *model.Role, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "PsqlRepository:FindRoleById")
	defer trace.Finish()

	resp = new(model.Role)

	query := fmt.Sprintf("%s where id=$1", queryRole)
	err = p.slave.Get(ctx, resp, query, id)
	if err != nil {
		trace.SetError(err)
		logger.Log.Error(ctx, err)

		return
	}

	return
}
