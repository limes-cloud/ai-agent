package service

import (
	"github.com/limes-cloud/ai-agent/api/errors"
	"github.com/limes-cloud/ai-agent/internal/core"
	"github.com/limes-cloud/ai-agent/internal/domain/entity"
	"github.com/limes-cloud/ai-agent/internal/domain/repository"
	"github.com/limes-cloud/ai-agent/internal/types"
)

type Prompt struct {
	repo repository.Prompt
}

func NewPrompt(repo repository.Prompt) *Prompt {
	return &Prompt{
		repo: repo,
	}
}

// ListPromptClassify 获取密钥分组列表
func (t *Prompt) ListPromptClassify(ctx core.Context, req *types.ListPromptClassifyRequest) ([]*entity.PromptClassify, uint32, error) {
	list, total, err := t.repo.ListPromptClassify(ctx, req)
	if err != nil {
		return nil, 0, errors.ListError(err.Error())
	}
	return list, total, nil
}

// CreatePromptClassify 创建密钥分组
func (t *Prompt) CreatePromptClassify(ctx core.Context, tg *entity.PromptClassify) (uint32, error) {
	id, err := t.repo.CreatePromptClassify(ctx, tg)
	if err != nil {
		return 0, errors.CreateError(err.Error())
	}
	return id, nil
}

// UpdatePromptClassify 更新密钥分组
func (t *Prompt) UpdatePromptClassify(ctx core.Context, tg *entity.PromptClassify) error {
	if err := t.repo.UpdatePromptClassify(ctx, tg); err != nil {
		return errors.UpdateError(err.Error())
	}
	return nil
}

// DeletePromptClassify 删除密钥分组（app 纬度隔离）
func (t *Prompt) DeletePromptClassify(ctx core.Context, app string, id uint32) error {
	err := t.repo.DeletePromptClassify(ctx, app, id)
	if err != nil {
		return errors.DeleteError(err.Error())
	}
	return nil
}

// GetPrompt 获取指定的密钥信息（app 纬度隔离）
func (t *Prompt) GetPrompt(ctx core.Context, app string, id uint32) (*entity.Prompt, error) {
	res, err := t.repo.GetPrompt(ctx, app, id)
	if err != nil {
		return nil, errors.GetError(err.Error())
	}
	return res, nil
}

// ListPrompt 获取密钥信息列表
func (t *Prompt) ListPrompt(ctx core.Context, req *types.ListPromptRequest) ([]*entity.Prompt, uint32, error) {
	list, total, err := t.repo.ListPrompt(ctx, req)
	if err != nil {
		return nil, 0, errors.ListError(err.Error())
	}
	return list, total, nil
}

// CreatePrompt 创建密钥信息
func (t *Prompt) CreatePrompt(ctx core.Context, task *entity.Prompt) (uint32, error) {
	id, err := t.repo.CreatePrompt(ctx, task)
	if err != nil {
		return 0, errors.CreateError(err.Error())
	}
	return id, nil
}

// UpdatePrompt 更新密钥信息
func (t *Prompt) UpdatePrompt(ctx core.Context, task *entity.Prompt) error {
	if err := t.repo.UpdatePrompt(ctx, task); err != nil {
		return errors.UpdateError(err.Error())
	}
	return nil
}

// DeletePrompt 删除密钥信息（app 纬度隔离）
func (t *Prompt) DeletePrompt(ctx core.Context, app string, id uint32) error {
	if err := t.repo.DeletePrompt(ctx, app, id); err != nil {
		return errors.DeleteError(err.Error())
	}
	return nil
}
