package types

type ListModelRequest struct {
	Page          uint32   `json:"page"`
	PageSize      uint32   `json:"pageSize"`
	Type          *string  `json:"type"`
	UUID          *string  `json:"uuid"`
	SecretGroupId *uint32  `json:"secretGroupId"`
	Ids           []uint32 `json:"ids"`
}
