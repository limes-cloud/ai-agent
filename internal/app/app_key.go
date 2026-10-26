package app

import (
	"context"

	"github.com/limes-cloud/ai-agent/api/errors"
	"github.com/limes-cloud/ai-agent/internal/core"
)

// currentApp 从请求上下文取应用标识（X-App-Key），缺失即拒绝（app 纬度隔离的入口）。
func currentApp(ctx context.Context) (string, error) {
	app := core.MustContext(ctx).AppKey()
	if app == "" {
		return "", errors.ParamsError("缺少应用标识（X-App-Key）")
	}
	return app, nil
}

// maskToken 密钥脱敏：仅保留尾 4 位，其余以 **** 代替。
func maskToken(token string) string {
	if len(token) <= 4 {
		return "****"
	}
	return "****" + token[len(token)-4:]
}