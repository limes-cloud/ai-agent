package repository

import (
	"github.com/limes-cloud/ai-agent/internal/core"
	"github.com/limes-cloud/ai-agent/internal/domain/entity"
	"github.com/limes-cloud/ai-agent/internal/types"
)

type Model interface {
	// GetModel 获取指定的模型信息
	GetModel(ctx core.Context, id uint32) (*entity.Model, error)

	// GetModelByUUID 获取指定的模型信息
	GetModelByUUID(ctx core.Context, uuid string) (*entity.Model, error)

	// ListModel 获取模型信息列表
	ListModel(ctx core.Context, req *types.ListModelRequest) ([]*entity.Model, uint32, error)

	// CreateModel 创建模型信息（Status 缺省按启用）
	CreateModel(ctx core.Context, req *entity.Model) (uint32, error)

	// UpdateModel 更新模型信息
	UpdateModel(ctx core.Context, req *entity.Model) error

	// DeleteModel 删除模型信息
	DeleteModel(ctx core.Context, id uint32) error
}
