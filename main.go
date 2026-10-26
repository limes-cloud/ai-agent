package main

import (
	"context"

	"github.com/limes-cloud/manager/api/authorize"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/limes-cloud/ai-agent/api/errors"
	"github.com/limes-cloud/ai-agent/internal/app"
	"github.com/limes-cloud/ai-agent/internal/core"
	initentity "github.com/limes-cloud/ai-agent/internal/initiator/entity"
	initmenu "github.com/limes-cloud/ai-agent/internal/initiator/menu"
	"github.com/limes-cloud/kratosx"
	_ "go.uber.org/automaxprocs"
)

func main() {
	srv := core.InitApp(
		kratosx.WithRegistrarServer(app.Register),
		kratosx.WithValidateErrHook(func(ctx context.Context, err error) error {
			c := core.MustContext(ctx)
			c.Logger().Warnw("msg", "params validate error", "err", err)
			return errors.ParamsError()
		}),
		kratosx.WithMiddleware(authorize.AuthMiddleware()),
	)

	initentity.Register(core.MustContext(context.Background()))
	initmenu.Register(core.MustContext(context.Background()))

	if err := srv.App().Run(); err != nil {
		log.Fatal(err)
	}
}
