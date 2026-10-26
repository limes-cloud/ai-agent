package types

type ListPromptClassifyRequest struct {
	Page     uint32  `json:"page"`
	PageSize uint32  `json:"pageSize"`
	Name     *string `json:"name"`

	// 中转参数（由 app 层根据 X-App-Key 推导设置，不入 proto）
	App string `json:"-"`
}

type ListPromptRequest struct {
	Page       uint32  `json:"page"`
	PageSize   uint32  `json:"pageSize"`
	ClassifyId *uint32 `json:"classifyId"`
	ModelId    *uint32 `json:"modelId"`
	Title      *string `json:"title"`
	Status     *bool   `json:"status"`

	// 中转参数（由 app 层根据 X-App-Key 推导设置，不入 proto）
	App string `json:"-"`
}
