package repository

import (
	"github.com/limes-cloud/ai-agent/internal/core"
	"github.com/limes-cloud/ai-agent/internal/domain/entity"
	"github.com/limes-cloud/ai-agent/internal/types"
)

type Prompt interface {
	// GetPromptClassify 获取指定的提示词分组（app 纬度隔离）
	GetPromptClassify(ctx core.Context, app string, id uint32) (*entity.PromptClassify, error)

	// ListPromptClassify 获取提示词分组列表（req.App 必填）
	ListPromptClassify(ctx core.Context, req *types.ListPromptClassifyRequest) ([]*entity.PromptClassify, uint32, error)

	// CreatePromptClassify 创建提示词分组（req.App 必填）
	CreatePromptClassify(ctx core.Context, req *entity.PromptClassify) (uint32, error)

	// UpdatePromptClassify 更新提示词分组
	UpdatePromptClassify(ctx core.Context, req *entity.PromptClassify) error

	// DeletePromptClassify 删除提示词分组（app 纬度隔离）
	DeletePromptClassify(ctx core.Context, app string, id uint32) error

	// GetPromptByUUID 获取指定的提示词信息（app 纬度隔离）
	GetPromptByUUID(ctx core.Context, app string, uuid string) (*entity.Prompt, error)

	// GetPrompt 获取指定的提示词信息（app 纬度隔离）
	GetPrompt(ctx core.Context, app string, id uint32) (*entity.Prompt, error)

	// ListPrompt 获取提示词信息列表（req.App 必填）
	ListPrompt(ctx core.Context, req *types.ListPromptRequest) ([]*entity.Prompt, uint32, error)

	// CreatePrompt 创建提示词信息（req.App 必填）
	CreatePrompt(ctx core.Context, req *entity.Prompt) (uint32, error)

	// UpdatePrompt 更新提示词信息
	UpdatePrompt(ctx core.Context, req *entity.Prompt) error

	// DeletePrompt 删除提示词信息（app 纬度隔离）
	DeletePrompt(ctx core.Context, app string, id uint32) error
}
