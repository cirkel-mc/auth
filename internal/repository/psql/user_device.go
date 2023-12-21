package psql

import (
	"cirkel/auth/internal/domain/model"
	"context"

	"github.com/cirkel-mc/goutils/logger"
	"github.com/cirkel-mc/goutils/tracer"
)

func (p *psqlRepository) CreateUserDevice(ctx context.Context, ud *model.UserDevice) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "PsqlRepository:CreateUserDevice")
	defer trace.Finish()

	query := `
	insert into user.user_device
	(created_at, user_id, device_id, channel, user_agent, fcm_token)
	values ($1, $2, $3, $4, $5, $6)
	`

	err = p.master.Preparex(ctx, query,
		ud.CreatedAt,
		ud.UserId,
		ud.DeviceId,
		ud.Channel,
		ud.UserAgent,
		ud.FcmToken,
	)
	if err != nil {
		trace.SetError(err)
		logger.Log.Errorf(ctx, "failed insert into user_device: %s", err)

		return
	}

	return
}
