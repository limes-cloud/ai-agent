package dbs

import (
	"errors"
	"strings"

	"github.com/limes-cloud/ai-agent/internal/core"

	"golang.org/x/exp/rand"

	"github.com/limes-cloud/ai-agent/internal/domain/entity"
	"github.com/limes-cloud/ai-agent/internal/types"
)

type Secret struct{}

func NewSecret() *Secret {
	return &Secret{}
}

// GetRandomSecretByGroupId 随机获取指定分组下面的密钥
func (t *Secret) GetRandomSecretByGroupId(ctx core.Context, id uint32) (*entity.Secret, error) {
	var list []*entity.Secret
	if err := ctx.DB().
		Where("group_id = ? and status = true", id).
		Preload("Group").
		Find(&list).Error; err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, errors.New("not exist available secret")
	}

	if len(list) == 1 {
		return list[0], nil
	}

	index := rand.Intn(len(list))

	return list[index], nil
}

// GetSecretGroup 获取指定的数据
func (t *Secret) GetSecretGroup(ctx core.Context, id uint32) (*entity.SecretGroup, error) {
	var (
		secretGroup = entity.SecretGroup{}
		fs          = []string{"*"}
	)
	return &secretGroup, ctx.DB().Select(fs).First(&secretGroup, id).Error
}

// ListSecretGroup 获取列表
func (t *Secret) ListSecretGroup(ctx core.Context, req *types.ListSecretGroupRequest) ([]*entity.SecretGroup, uint32, error) {
	var (
		list  []*entity.SecretGroup
		total int64
		fs    = []string{"*"}
	)

	db := ctx.DB().Model(entity.SecretGroup{}).Select(fs)
	if req.Name != nil {
		db = db.Where("name LIKE ?", *req.Name+"%")
	}
	if req.Supplier != nil {
		db = db.Where("supplier = ?", *req.Supplier)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	db = db.Offset(int((req.Page - 1) * req.PageSize)).Limit(int(req.PageSize))

	return list, uint32(total), db.Find(&list).Error
}

// CreateSecretGroup 创建数据
func (t *Secret) CreateSecretGroup(ctx core.Context, tg *entity.SecretGroup) (uint32, error) {
	return tg.Id, ctx.DB().Create(tg).Error
}

// UpdateSecretGroup 更新数据（全量更新，零值字段也写入）
func (t *Secret) UpdateSecretGroup(ctx core.Context, tg *entity.SecretGroup) error {
	return ctx.DB().Model(&entity.SecretGroup{}).Where("id = ?", tg.Id).
		Updates(map[string]interface{}{
			"name":     tg.Name,
			"host":     tg.Host,
			"supplier": tg.Supplier,
		}).Error
}

// DeleteSecretGroup 删除数据
func (t *Secret) DeleteSecretGroup(ctx core.Context, id uint32) error {
	return ctx.DB().Where("id = ?", id).Delete(&entity.SecretGroup{}).Error
}

// GetSecret 获取指定的数据
func (t *Secret) GetSecret(ctx core.Context, id uint32) (*entity.Secret, error) {
	var (
		secret = entity.Secret{}
		fs     = []string{"*"}
	)
	return &secret, ctx.DB().Select(fs).
		Preload("Group").
		First(&secret, id).Error
}

// ListSecret 获取列表
func (t *Secret) ListSecret(ctx core.Context, req *types.ListSecretRequest) ([]*entity.Secret, uint32, error) {
	var (
		list  []*entity.Secret
		total int64
		fs    = []string{"*"}
	)

	db := ctx.DB().Model(entity.Secret{}).Select(fs).Preload("Group")
	if req.GroupId != nil {
		db = db.Where("group_id = ?", *req.GroupId)
	}
	if req.Name != nil {
		db = db.Where("name like ?", *req.Name+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	db = db.Offset(int((req.Page - 1) * req.PageSize)).Limit(int(req.PageSize))
	return list, uint32(total), db.Find(&list).Error
}

// CreateSecret 创建数据
func (t *Secret) CreateSecret(ctx core.Context, secret *entity.Secret) (uint32, error) {
	return secret.Id, ctx.DB().Create(secret).Error
}

// UpdateSecret 更新数据（Updates(struct)：零值不更新，optional 指针字段非 nil 即可更新；token 为脱敏值（****）时置空跳过，防止回填覆盖真实密钥）
func (t *Secret) UpdateSecret(ctx core.Context, secret *entity.Secret) error {
	if strings.Contains(secret.Token, "*") {
		secret.Token = ""
	}
	return ctx.DB().Model(&entity.Secret{}).Where("id = ?", secret.Id).Updates(secret).Error
}

// DeleteSecret 删除数据
func (t *Secret) DeleteSecret(ctx core.Context, id uint32) error {
	return ctx.DB().Where("id = ?", id).Delete(&entity.Secret{}).Error
}
