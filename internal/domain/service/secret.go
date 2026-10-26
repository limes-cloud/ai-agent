package service

import (
	"github.com/limes-cloud/ai-agent/internal/core"
	"google.golang.org/protobuf/proto"

	"github.com/limes-cloud/ai-agent/api/errors"
	"github.com/limes-cloud/ai-agent/internal/domain/entity"
	"github.com/limes-cloud/ai-agent/internal/domain/repository"
	"github.com/limes-cloud/ai-agent/internal/types"
)

type Secret struct {
	repo repository.Secret
}

func NewSecret(repo repository.Secret) *Secret {
	return &Secret{
		repo: repo,
	}
}

// ListSupplier 获取供应商列表
func (t *Secret) ListSupplier() []string {
	return []string{
		entity.SupplierARKBot,
		entity.SupplierDeepSeek,
		entity.SupplierClaude,
		entity.SupplierGemini,
		entity.SupplierOllama,
		entity.SupplierOpenAl,
		entity.SupplierQianfan,
		entity.SupplierQwen,
	}
}

// GetSecretGroup 获取指定的密钥分组
func (t *Secret) GetSecretGroup(ctx core.Context, id uint32) (*entity.SecretGroup, error) {
	res, err := t.repo.GetSecretGroup(ctx, id)
	if err != nil {
		return nil, errors.GetError(err.Error())
	}
	return res, nil
}

// ListSecretGroup 获取密钥分组列表
func (t *Secret) ListSecretGroup(ctx core.Context, req *types.ListSecretGroupRequest) ([]*entity.SecretGroup, uint32, error) {
	list, total, err := t.repo.ListSecretGroup(ctx, req)
	if err != nil {
		return nil, 0, errors.ListError(err.Error())
	}
	return list, total, nil
}

// CreateSecretGroup 创建密钥分组
func (t *Secret) CreateSecretGroup(ctx core.Context, tg *entity.SecretGroup) (uint32, error) {
	id, err := t.repo.CreateSecretGroup(ctx, tg)
	if err != nil {
		return 0, errors.CreateError(err.Error())
	}
	return id, nil
}

// UpdateSecretGroup 更新密钥分组
func (t *Secret) UpdateSecretGroup(ctx core.Context, tg *entity.SecretGroup) error {
	if err := t.repo.UpdateSecretGroup(ctx, tg); err != nil {
		return errors.UpdateError(err.Error())
	}
	return nil
}

// DeleteSecretGroup 删除密钥分组
func (t *Secret) DeleteSecretGroup(ctx core.Context, id uint32) error {
	err := t.repo.DeleteSecretGroup(ctx, id)
	if err != nil {
		return errors.DeleteError(err.Error())
	}
	return nil
}

// GetSecret 获取指定的密钥信息
func (t *Secret) GetSecret(ctx core.Context, id uint32) (*entity.Secret, error) {
	res, err := t.repo.GetSecret(ctx, id)
	if err != nil {
		return nil, errors.GetError(err.Error())
	}
	return res, nil
}

// ListSecret 获取密钥信息列表
func (t *Secret) ListSecret(ctx core.Context, req *types.ListSecretRequest) ([]*entity.Secret, uint32, error) {
	list, total, err := t.repo.ListSecret(ctx, req)
	if err != nil {
		return nil, 0, errors.ListError(err.Error())
	}
	return list, total, nil
}

// CreateSecret 创建密钥信息
func (t *Secret) CreateSecret(ctx core.Context, task *entity.Secret) (uint32, error) {
	task.Status = proto.Bool(true)
	id, err := t.repo.CreateSecret(ctx, task)
	if err != nil {
		return 0, errors.CreateError(err.Error())
	}
	return id, nil
}

// UpdateSecret 更新密钥信息
func (t *Secret) UpdateSecret(ctx core.Context, task *entity.Secret) error {
	if err := t.repo.UpdateSecret(ctx, task); err != nil {
		return errors.UpdateError(err.Error())
	}
	return nil
}

// DeleteSecret 删除密钥信息
func (t *Secret) DeleteSecret(ctx core.Context, id uint32) error {
	if err := t.repo.DeleteSecret(ctx, id); err != nil {
		return errors.DeleteError(err.Error())
	}
	return nil
}
