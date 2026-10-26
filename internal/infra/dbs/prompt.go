package dbs

import (
	"github.com/limes-cloud/ai-agent/internal/core"

	"github.com/limes-cloud/ai-agent/internal/domain/entity"
	"github.com/limes-cloud/ai-agent/internal/types"
)

type Prompt struct{}

func NewPrompt() *Prompt {
	return &Prompt{}
}

// GetPromptClassify 获取指定的数据（app 纬度隔离）
func (t *Prompt) GetPromptClassify(ctx core.Context, app string, id uint32) (*entity.PromptClassify, error) {
	classify := entity.PromptClassify{}
	return &classify, ctx.DB().Where("app = ?", app).First(&classify, id).Error
}

// ListPromptClassify 获取列表（app 纬度隔离）
func (t *Prompt) ListPromptClassify(ctx core.Context, req *types.ListPromptClassifyRequest) ([]*entity.PromptClassify, uint32, error) {
	var (
		list  []*entity.PromptClassify
		total int64
		fs    = []string{"*"}
	)

	db := ctx.DB().Model(entity.PromptClassify{}).Select(fs).Where("app = ?", req.App)
	if req.Name != nil {
		db = db.Where("name LIKE ?", *req.Name+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	db = db.Offset(int((req.Page - 1) * req.PageSize)).Limit(int(req.PageSize))

	return list, uint32(total), db.Find(&list).Error
}

// CreatePromptClassify 创建数据
func (t *Prompt) CreatePromptClassify(ctx core.Context, tg *entity.PromptClassify) (uint32, error) {
	return tg.Id, ctx.DB().Create(tg).Error
}

// UpdatePromptClassify 更新数据（app 不可改）
func (t *Prompt) UpdatePromptClassify(ctx core.Context, tg *entity.PromptClassify) error {
	return ctx.DB().Model(&entity.PromptClassify{}).Where("id = ?", tg.Id).
		Updates(map[string]interface{}{
			"name":        tg.Name,
			"description": tg.Description,
		}).Error
}

// DeletePromptClassify 删除数据（app 纬度隔离）
func (t *Prompt) DeletePromptClassify(ctx core.Context, app string, id uint32) error {
	return ctx.DB().Where("id = ? AND app = ?", id, app).Delete(&entity.PromptClassify{}).Error
}

// GetPrompt 获取指定的数据（app 纬度隔离）
func (t *Prompt) GetPrompt(ctx core.Context, app string, id uint32) (*entity.Prompt, error) {
	prompt := entity.Prompt{}
	return &prompt, ctx.DB().
		Preload("Classify").
		Preload("Model").
		Where("app = ?", app).
		First(&prompt, id).Error
}

// GetPromptByUUID 获取指定的数据（app 纬度隔离）
func (t *Prompt) GetPromptByUUID(ctx core.Context, app string, uuid string) (*entity.Prompt, error) {
	prompt := entity.Prompt{}
	return &prompt, ctx.DB().
		Preload("Classify").
		Preload("Model").
		Where("app = ? AND uuid = ?", app, uuid).
		First(&prompt).Error
}

// ListPrompt 获取列表（app 纬度隔离）
func (t *Prompt) ListPrompt(ctx core.Context, req *types.ListPromptRequest) ([]*entity.Prompt, uint32, error) {
	var (
		list  []*entity.Prompt
		total int64
	)

	db := ctx.DB().Model(entity.Prompt{}).Preload("Classify").Preload("Model").
		Where("app = ?", req.App)
	if req.ClassifyId != nil {
		db = db.Where("classify_id = ?", *req.ClassifyId)
	}
	if req.ModelId != nil {
		db = db.Where("model_id = ?", *req.ModelId)
	}
	if req.Title != nil {
		db = db.Where("title like ?", *req.Title+"%")
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

// CreatePrompt 创建数据
func (t *Prompt) CreatePrompt(ctx core.Context, prompt *entity.Prompt) (uint32, error) {
	return prompt.Id, ctx.DB().Create(prompt).Error
}

// UpdatePrompt 更新数据（Updates(struct)：零值不更新，optional 指针字段非 nil 即可更新——符合 rules 约定）
func (t *Prompt) UpdatePrompt(ctx core.Context, prompt *entity.Prompt) error {
	return ctx.DB().Model(&entity.Prompt{}).Where("id = ?", prompt.Id).Updates(prompt).Error
}

// DeletePrompt 删除数据（app 纬度隔离）
func (t *Prompt) DeletePrompt(ctx core.Context, app string, id uint32) error {
	return ctx.DB().Where("id = ? AND app = ?", id, app).Delete(&entity.Prompt{}).Error
}
