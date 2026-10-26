package entity

import (
	"encoding/json"

	"github.com/limes-cloud/kratosx/model"
)

type PromptClassify struct {
	App         string `json:"app" gorm:"column:app;index:idx_prompt_classify_app_name,priority:1"`   // 所属应用（app 纬度隔离）
	Name        string `json:"name" gorm:"column:name;index:idx_prompt_classify_app_name,priority:2"` // 分类名
	Description string `json:"description" gorm:"column:description"`
	model.BaseModel
}

type Prompt struct {
	App         string          `json:"app" gorm:"column:app;index:idx_prompt_app_uuid,priority:1"`   // 所属应用（app 纬度隔离）
	UUID        string          `json:"uuid" gorm:"column:uuid;index:idx_prompt_app_uuid,priority:2"` // 唯一标识
	Title       string          `json:"title" gorm:"column:title"`                                    // 提示词标题
	System      *string         `json:"system" gorm:"column:system"`                                  // 系统提示
	Prompt      string          `json:"prompt" gorm:"column:prompt"`                                  // 提示词模版
	Status      *uint32         `json:"status" gorm:"column:status"`                                  // 状态 0:禁用 1:启用
	Description string          `json:"description" gorm:"column:description"`
	Extra       string          `json:"extra" gorm:"column:extra"`                                    // 额外信息
	Version     string          `json:"version" gorm:"column:version"`                                // 版本号
	ModelId     uint32          `json:"modelId" gorm:"column:model_id"`                               // 模型ID
	ClassifyId  uint32          `json:"classifyId" gorm:"column:classify_id"`
	Model       *Model          `json:"model" gorm:"foreignKey:ModelId;references:Id"`
	Classify    *PromptClassify `json:"classify" gorm:"foreignKey:ClassifyId;references:Id"`
	model.BaseModel
}

func (p *Prompt) GetExtra() *PromptExtra {
	extra := &PromptExtra{}
	_ = json.Unmarshal([]byte(p.Extra), extra)
	return extra
}

type PromptExtra struct {
	Temperature  *float32 `json:"temperature"`  // 温度，默认0.8
	TopP         *float32 `json:"topP"`         // 核采样概率总和阈值
	FormatJson   *bool    `json:"formatJson"`   // 是否json格式返回
	OutputSchema *string  `json:"outputSchema"` // 输出格式

	Path  string `json:"path"`
	Store string `json:"store"`
	Size  string `json:"size"`
}

//type PromptHistory struct {
//	PromptID uint32  `json:"prompt_id"`
//	Version  string `json:"version"`
//	Data     string `json:"data"`
//	model.BaseModel
//}
