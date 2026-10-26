package entity

import (
	"github.com/limes-cloud/kratosx/model"
)

const (
	SupplierARKBot   = "arkbot"
	SupplierDeepSeek = "deepseek"
	SupplierClaude   = "claude"
	SupplierGemini   = "gemini"
	SupplierOllama   = "ollama"
	SupplierOpenAl   = "openal"
	SupplierQianfan  = "qianfan"
	SupplierQwen     = "qwen"
)

type Supplier struct {
	Keyword string `json:"keyword"` // 供应商关键字
	Name    string `json:"name"`    // 供应商名称
}

type SecretGroup struct {
	Supplier string `json:"supplier" gorm:"column:supplier"` // 供应商
	Name     string `json:"name" gorm:"column:name"`         // 分组名称
	Host     string `json:"host" gorm:"column:host"`         // 服务地址
	model.BaseModel
}

type Secret struct {
	GroupId uint32       `json:"groupId" gorm:"column:group_id"` // 分组
	Token   string       `json:"token" gorm:"column:token"`      // 密钥
	Status  *bool        `json:"status" gorm:"column:status"`    // 状态
	Reason  *string      `json:"reason" gorm:"column:reason"`    // 禁用原因
	Name    string       `json:"name" gorm:"column:name"`        // 名称
	Group   *SecretGroup `json:"group" gorm:"foreignKey:GroupId;references:Id"`
	model.BaseModel
}
