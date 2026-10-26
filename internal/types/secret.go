package types

type ListSecretGroupRequest struct {
	Page     uint32  `json:"page"`
	PageSize uint32  `json:"pageSize"`
	Name     *string `json:"name"`
	Supplier *string `json:"supplier"`
}

type ListSecretRequest struct {
	Page     uint32  `json:"page"`
	PageSize uint32  `json:"pageSize"`
	Name     *string `json:"description"`
	GroupId  *uint32 `json:"groupId"`
	Status   *bool   `json:"status"`
}
