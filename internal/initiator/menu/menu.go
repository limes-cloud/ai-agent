package menu

import (
	"context"

	"github.com/limes-cloud/ai-agent/internal/core"
	managerMenu "github.com/limes-cloud/manager/api/menu"
)

const registerModeAppend = "append"

func Register(ctx core.Context) {
	ctx.BeforeStart("register ai-agent menu", func() {
		conn, err := ctx.GrpcConn("Manager")
		if err != nil {
			ctx.Logger().Errorw("msg", "get manager grpc conn error", "err", err.Error())
			return
		}

		reply, err := managerMenu.NewMenuClient(conn).RegisterMenu(context.Background(), &managerMenu.RegisterMenuRequest{
			AppKeyword: appKeyword,
			Mode:       registerModeAppend,
			Menus:      menus,
		})
		if err != nil {
			ctx.Logger().Errorw("msg", "register ai-agent menu error", "err", err.Error())
			return
		}

		ctx.Logger().Infow(
			"msg", "register ai-agent menu success",
			"created", reply.Created,
			"updated", reply.Updated,
			"deleted", reply.Deleted,
		)
	})
}
