package repository

import (
	"github.com/limes-cloud/ai-agent/internal/core"
	"github.com/limes-cloud/ai-agent/internal/domain/entity"
	"github.com/limes-cloud/ai-agent/internal/types"
)

type Secret interface {
	// GetRandomSecretByGroupId 获取随机的密钥
	GetRandomSecretByGroupId(ctx core.Context, id uint32) (*entity.Secret, error)

	// GetSecretGroup 获取指定的密钥分组
	GetSecretGroup(ctx core.Context, id uint32) (*entity.SecretGroup, error)

	// ListSecretGroup 获取密钥分组列表
	ListSecretGroup(ctx core.Context, req *types.ListSecretGroupRequest) ([]*entity.SecretGroup, uint32, error)

	// CreateSecretGroup 创建密钥分组
	CreateSecretGroup(ctx core.Context, req *entity.SecretGroup) (uint32, error)

	// UpdateSecretGroup 更新密钥分组
	UpdateSecretGroup(ctx core.Context, req *entity.SecretGroup) error

	// DeleteSecretGroup 删除密钥分组
	DeleteSecretGroup(ctx core.Context, id uint32) error

	// GetSecret 获取指定的密钥信息
	GetSecret(ctx core.Context, id uint32) (*entity.Secret, error)

	// ListSecret 获取密钥信息列表
	ListSecret(ctx core.Context, req *types.ListSecretRequest) ([]*entity.Secret, uint32, error)

	// CreateSecret 创建密钥信息
	CreateSecret(ctx core.Context, req *entity.Secret) (uint32, error)

	// UpdateSecret 更新密钥信息
	UpdateSecret(ctx core.Context, req *entity.Secret) error

	// DeleteSecret 删除密钥信息
	DeleteSecret(ctx core.Context, id uint32) error
}
