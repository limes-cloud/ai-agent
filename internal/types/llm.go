package types

import (
	"github.com/cloudwego/eino-ext/libs/acl/openai"
)

type EmbeddingRequest struct {
	Texts []string `json:"texts"`
}

type EmbeddingResponse struct {
	Vectors [][]float64 `json:"vectors"`
}

type TextAttachment struct {
	Images []string `json:"images"`
}

type TextRequest struct {
	App            string                               // 应用标识（app 层注入，按 app 隔离）
	Model          string
	System         string                               `json:"system"`         // 系统提示
	Prompt         string                               `json:"prompt"`         // 输入的文本
	Variables      map[string]any                       `json:"variable"`       // 变量
	Stream         bool                                 `json:"stream"`         // 流式输出
	ResponseFormat *openai.ChatCompletionResponseFormat `json:"responseFormat"` // 响应格式
	Temperature    *float32                             `json:"temperature"`    // 温度，默认0.8
	TopP           *float32                             `json:"topP"`           // 核采样概率总和阈值
	Attachment     *TextAttachment                      `json:"attachment"`     // 附件
	RecvFunc       func(chunk string)                   // 接收流式数据
}

type TextResponse struct {
	Content string         `json:"content"`
	Format  map[string]any `json:"format"`
}

type ImageRequest struct {
	App   string // 应用标识（app 层注入，按 app 隔离）
	Model string
	Path  string
	Store string

	Prompt    string         `json:"prompt"`   // 输入的文本
	Variables map[string]any `json:"variable"` // 变量
	Size      string         `json:"size"`     // 图片尺寸
	Refs      []string       `json:"refs"`     // 参考图
}

type ImageResponse struct {
	List []string `json:"list"`
}

type PromptTextRequest struct {
	App        string             // 应用标识（app 层注入，按 app 隔离）
	Prompt     string             `json:"prompt"`     // 提示词
	Variables  map[string]any     `json:"variable"`   // 变量
	Stream     bool               `json:"stream"`     // 流式输出
	Attachment *TextAttachment    `json:"attachment"` // 附件
	RecvFunc   func(chunk string) // 接收流式数据
}

type PromptTextResponse struct {
	Content string         `json:"content"`
	Format  map[string]any `json:"format"`
}

type PromptImageRequest struct {
	App       string         // 应用标识（app 层注入，按 app 隔离）
	Prompt    string         `json:"prompt"`   // 提示词
	Variables map[string]any `json:"variable"` // 变量
}

type PromptImageResponse struct {
	List []string `json:"list"`
}
