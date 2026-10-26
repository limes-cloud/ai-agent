package entity

import (
	"encoding/json"

	"github.com/limes-cloud/kratosx/model"
)

const (
	ModelTypeEmbedding = "embedding"
	ModelTypeChat      = "chat"
	ModelTypeImage     = "image"
	ModelTypeVideo     = "video"
)

type Model struct {
	Type          string       `json:"type" gorm:"column:type"`                           // 模型类型：embedding | chat | image | video
	UUID          string       `json:"uuid" gorm:"column:uuid;uniqueIndex:uk_model_uuid"` // 模型唯一标识（全局唯一，对外暴露的 ID）
	Keyword       string       `json:"keyword" gorm:"column:keyword"`                     // 模型关键字（供应商侧标识）
	SecretGroupId uint32       `json:"secretGroupId" gorm:"column:secret_group_id"`                 // 关联的密钥分组 ID（全局组）
	Description   string       `json:"description" gorm:"column:description"`                        // 模型描述
	Extra         string       `json:"extra" gorm:"column:extra"`                                   // 扩展配置（JSON 字符串，见 GetExtra()）
	Status        *uint32      `json:"status" gorm:"column:status"`                                 // 状态：0 禁用 | 1 启用（可空，按启用处理）
	SecretGroup   *SecretGroup `json:"secretGroup" gorm:"foreignKey:SecretGroupId;references:Id"`   // 关联的密钥分组（GORM Preload）
	model.BaseModel
}

func (m *Model) GetExtra() *ModelExtra {
	extra := &ModelExtra{}
	if m.Extra != "" {
		_ = json.Unmarshal([]byte(m.Extra), extra)
	}
	return extra
}

type ModelExtra struct {
	SupportJson *bool `json:"supportJson"`
	SupportTool *bool `json:"supportTool"`
}
