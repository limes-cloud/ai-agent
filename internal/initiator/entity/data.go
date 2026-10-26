package entity

import managerEntity "github.com/limes-cloud/manager/api/entity"

const appKeyword = "ai-agent"

var registerTables = map[string]struct{}{
	"model":           {},
	"prompt":          {},
	"prompt_classify": {},
	"secret":          {},
	"secret_group":    {},
}

var ignoreFields = map[string]struct{}{
	"deleted_at": {},
}

func shouldRegisterTable(name string) bool {
	_, ok := registerTables[name]
	return ok
}

func shouldRegisterField(name string) bool {
	_, ok := ignoreFields[name]
	return !ok
}

func registerField(name string, comment string, index uint32) *managerEntity.RegisterEntityField {
	return &managerEntity.RegisterEntityField{
		Name:    name,
		Comment: comment,
		Index:   index,
	}
}
