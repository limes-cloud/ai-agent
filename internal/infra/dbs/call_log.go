package dbs

import (
	"github.com/limes-cloud/ai-agent/internal/core"
	"github.com/limes-cloud/ai-agent/internal/domain/entity"
)

type CallLog struct{}

func NewCallLog() *CallLog {
	return &CallLog{}
}

// CreateCallLog 写入一条生成调用日志
func (t *CallLog) CreateCallLog(ctx core.Context, log *entity.CallLog) (uint32, error) {
	return log.Id, ctx.DB().Create(log).Error
}