package psql

import (
	"cirkel/auth/internal/domain/model"
	"context"
	"fmt"

	"github.com/cirkel-mc/goutils/logger"
	"github.com/cirkel-mc/goutils/tracer"
)

const queryUser = `
select
	us.id, us.username, us.email, us.password, us.status, us.verified_at,
	ro.id AS "ro.id", ro.name AS "ro.name", ro.key AS "ro.key", us.is_partner
from "user".users as us
inner join "user".roles as ro on ro.id = us.role_id
`

func (p *psqlRepository) GetUserNextVal(ctx context.Context) (resp int, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "PsqlRepository:GetUserNextVal")
	defer trace.Finish()

	query := `select nextval('"user".users_id_seq')`
	err = p.master.Get(ctx, &resp, query)
	if err != nil {
		trace.SetError(err)
		logger.Log.Errorf(ctx, "failed to get nextval for table user: %s", err)

		return
	}

	return
}

func (p *psqlRepository) FindUserByUsername(ctx context.Context, username string) (resp *model.User, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "PsqlRepository:FindUserByUsername")
	defer trace.Finish()

	resp = new(model.User)
	query := fmt.Sprintf("%s where us.username=$1", queryUser)
	err = p.slave.Get(ctx, resp, query, username)
	if err != nil {
		trace.SetError(err)
		logger.Log.Errorf(ctx, "failed to get user by username: %s", err)

		return
	}

	return
}

func (p *psqlRepository) FindUserByEmail(ctx context.Context, email string) (resp *model.User, err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "PsqlRepository:FindUserByEmail")
	defer trace.Finish()

	resp = new(model.User)
	query := fmt.Sprintf("%s where us.email=$1", queryUser)
	err = p.slave.Get(ctx, resp, query, email)
	if err != nil {
		trace.SetError(err)
		logger.Log.Errorf(ctx, "failed to get user by email: %s", err)

		return
	}

	return
}

func (p *psqlRepository) CreateUser(ctx context.Context, user *model.User) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "PsqlRepository:CreateUser")
	defer trace.Finish()

	query := `
	insert into "user".users 
	(id, created_at, username, email, password, status, role_id)
	values ($1, $2, $3, $4, $5, $6, $7)
	`

	err = p.master.Preparex(ctx, query,
		user.Id,
		user.CreatedAt,
		user.Username,
		user.Email,
		user.Password,
		user.Status,
		user.Role.Id,
	)
	if err != nil {
		trace.SetError(err)
		logger.Log.Errorf(ctx, "failed to create new_user: %s", err)

		return
	}

	return
}

func (p *psqlRepository) UpdateUser(ctx context.Context, user *model.User) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "PsqlRepository:UpdateUser")
	defer trace.Finish()

	query := `
	update "user".users set
		updated_at=$2, status=$3, verified_at=$4
	where id=$1
	`

	err = p.master.Preparex(ctx, query,
		user.Id,
		user.UpdatedAt,
		user.Status,
		user.VerifiedAt,
	)
	if err != nil {
		trace.SetError(err)
		logger.Log.Errorf(ctx, "failed to update user: %s", err)

		return
	}

	return
}

func (p *psqlRepository) DeleteUser(ctx context.Context, user *model.User) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "PsqlRepository:DeleteUser")
	defer trace.Finish()

	query := `
	update "user".users set
		deleted_at=$2, deleted_by=$3, status=$4
	where id=$1
	`

	err = p.master.Preparex(ctx, query,
		user.Id,
		user.DeletedAt,
		user.DeletedBy,
		user.Status,
	)
	if err != nil {
		trace.SetError(err)
		logger.Log.Errorf(ctx, "failed to delete user: %s", err)

		return
	}

	return
}
