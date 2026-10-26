package repository

import (
	"github.com/limes-cloud/ai-agent/internal/core"
	"github.com/limes-cloud/ai-agent/internal/domain/entity"
)

// CallLog 生成调用日志仓储接口（按 app 计量）。
type CallLog interface {
	// CreateCallLog 写入一条生成调用日志
	CreateCallLog(ctx core.Context, log *entity.CallLog) (uint32, error)
}