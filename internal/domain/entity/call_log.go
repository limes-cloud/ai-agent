package entity

import (
	"github.com/limes-cloud/kratosx/model"
)

// CallLog 生成调用日志（按 app 计量的最小单元）。
type CallLog struct {
	App              string `json:"app" gorm:"column:app;index:idx_call_log_app_created,priority:1"`              // 应用标识（app 纬度计量）
	ModelId          uint32 `json:"modelId" gorm:"column:model_id"`                                              // 模型 ID
	ModelUUID        string `json:"modelUuid" gorm:"column:model_uuid;index:idx_call_log_app_created,priority:2"` // 模型 UUID
	Supplier         string `json:"supplier" gorm:"column:supplier"`                                              // 供应商（deepseek/claude/...）
	Scene            string `json:"scene" gorm:"column:scene"`                                                    // 生成场景：text | image | prompt_text | prompt_image | embedding
	PromptTokens     uint32 `json:"promptTokens" gorm:"column:prompt_tokens"`                                     // 输入 token 数（eino 未提供时为 0）
	CompletionTokens uint32 `json:"completionTokens" gorm:"column:completion_tokens"`                             // 输出 token 数
	LatencyMs        int64  `json:"latencyMs" gorm:"column:latency_ms"`                                           // 生成耗时（毫秒）
	Status           uint8  `json:"status" gorm:"column:status"`                                                  // 状态：0 失败 | 1 成功
	ErrorMessage     string `json:"errorMessage" gorm:"column:error_message"`                                     // 失败原因（成功为空）
	model.BaseModel
}