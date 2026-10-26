package dbs

import (
	"github.com/limes-cloud/kratosx/pkg/value"

	"github.com/limes-cloud/ai-agent/internal/core"
	"github.com/limes-cloud/ai-agent/internal/domain/entity"
	"github.com/limes-cloud/ai-agent/internal/types"
)

type Model struct{}

func NewModel() *Model {
	return &Model{}
}

// GetModel 获取指定的数据
func (t *Model) GetModel(ctx core.Context, id uint32) (*entity.Model, error) {
	var (
		model = entity.Model{}
		fs    = []string{"*"}
	)
	return &model, ctx.DB().Select(fs).
		Preload("SecretGroup").
		First(&model, id).Error
}

// GetModelByUUID 获取指定的数据
func (t *Model) GetModelByUUID(ctx core.Context, uuid string) (*entity.Model, error) {
	model := entity.Model{}
	return &model, ctx.DB().
		Preload("SecretGroup").
		Where("uuid = ?", uuid).
		First(&model).Error
}

// ListModel 获取列表
func (t *Model) ListModel(ctx core.Context, req *types.ListModelRequest) ([]*entity.Model, uint32, error) {
	var (
		list  []*entity.Model
		total int64
		fs    = []string{"*"}
	)

	db := ctx.DB().Model(entity.Model{}).Select(fs).Preload("SecretGroup")
	if req.Type != nil {
		db = db.Where("type = ?", *req.Type)
	}
	if req.UUID != nil {
		db = db.Where("uuid like ?", *req.UUID+"%")
	}
	if req.SecretGroupId != nil {
		db = db.Where("secret_group_id = ?", *req.SecretGroupId)
	}
	if len(req.Ids) != 0 {
		db = db.Where("id in ?", req.Ids)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	db = db.Offset(int((req.Page - 1) * req.PageSize)).Limit(int(req.PageSize))
	return list, uint32(total), db.Find(&list).Error
}

// CreateModel 创建数据（Status 缺省按启用处理）
func (t *Model) CreateModel(ctx core.Context, model *entity.Model) (uint32, error) {
	if model.Status == nil {
		model.Status = value.Pointer(uint32(1))
	}
	return model.Id, ctx.DB().Create(model).Error
}

// UpdateModel 更新数据（Updates(struct)：零值不更新，optional 指针字段非 nil 即可更新——符合 rules 约定）
func (t *Model) UpdateModel(ctx core.Context, model *entity.Model) error {
	return ctx.DB().Model(&entity.Model{}).Where("id = ?", model.Id).Updates(model).Error
}

// DeleteModel 删除数据
func (t *Model) DeleteModel(ctx core.Context, id uint32) error {
	return ctx.DB().Where("id = ?", id).Delete(&entity.Model{}).Error
}
