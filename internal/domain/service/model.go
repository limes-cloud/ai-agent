package service

import (
	"github.com/limes-cloud/ai-agent/api/errors"
	"github.com/limes-cloud/ai-agent/internal/core"
	"github.com/limes-cloud/ai-agent/internal/domain/entity"
	"github.com/limes-cloud/ai-agent/internal/domain/repository"
	"github.com/limes-cloud/ai-agent/internal/types"
)

type Model struct {
	repo repository.Model
}

func NewModel(repo repository.Model) *Model {
	return &Model{
		repo: repo,
	}
}

// ListModelType 获取模型类型列表
func (t *Model) ListModelType() []string {
	return []string{
		entity.ModelTypeEmbedding,
		entity.ModelTypeChat,
		entity.ModelTypeImage,
		entity.ModelTypeVideo,
	}
}

// ListModel 获取模型信息列表
func (t *Model) ListModel(ctx core.Context, req *types.ListModelRequest) ([]*entity.Model, uint32, error) {
	list, total, err := t.repo.ListModel(ctx, req)
	if err != nil {
		return nil, 0, errors.ListError(err.Error())
	}
	return list, total, nil
}

// CreateModel 创建模型信息
func (t *Model) CreateModel(ctx core.Context, task *entity.Model) (uint32, error) {
	id, err := t.repo.CreateModel(ctx, task)
	if err != nil {
		return 0, errors.CreateError(err.Error())
	}
	return id, nil
}

// UpdateModel 更新模型信息
func (t *Model) UpdateModel(ctx core.Context, task *entity.Model) error {
	if err := t.repo.UpdateModel(ctx, task); err != nil {
		return errors.UpdateError(err.Error())
	}
	return nil
}

// DeleteModel 删除模型信息
func (t *Model) DeleteModel(ctx core.Context, id uint32) error {
	if err := t.repo.DeleteModel(ctx, id); err != nil {
		return errors.DeleteError(err.Error())
	}
	return nil
}
