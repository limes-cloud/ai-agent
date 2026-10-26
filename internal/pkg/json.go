package pkg

import (
	"errors"

	"github.com/eino-contrib/jsonschema"
	json "github.com/json-iterator/go"
)

// Schema 生成JSONSchema
func Schema(st any) []byte {
	reflector := &jsonschema.Reflector{ExpandedStruct: false, DoNotReference: true}
	schema := reflector.Reflect(st)
	schemaJSON, _ := json.Marshal(schema)
	return schemaJSON
}

func MatchExtractJSON(input string) (string, error) {
	var depth int
	start := -1
	inString := false
	var runeStack []rune

	for i, r := range input {
		switch r {
		case '{':
			if !inString {
				if depth == 0 {
					// 记录最外层大括号的开始位置
					start = i
				}
				depth++
			}
			runeStack = append(runeStack, r)
		case '}':
			if !inString && depth > 0 {
				depth--
				runeStack = runeStack[:len(runeStack)-1]
				if depth == 0 {
					// 尝试解析JSON
					var js json.RawMessage
					subStr := input[start : i+1]
					if err := json.Unmarshal([]byte(subStr), &js); err == nil {
						// 如果解析成功，则返回提取的JSON字符串
						return subStr, nil
					}
				}
			} else if !inString && depth == 0 {
				// 处理额外的闭合大括号
				runeStack = runeStack[:len(runeStack)-1]
			}
		case '"':
			// 切换inString状态，并且忽略已转义的引号
			if i == 0 || input[i-1] != '\\' || (i > 1 && input[i-2] == '\\') {
				inString = !inString
			}
		}
	}

	if depth != 0 {
		// 如果大括号没有正确闭合，尝试查找一个闭合的子串
		for len(runeStack) > 0 && runeStack[len(runeStack)-1] == '{' {
			runeStack = runeStack[:len(runeStack)-1]
			start++
		}
		if start >= 0 && start < len(input) {
			return MatchExtractJSON(input[start:])
		}
	}

	return "", errors.New("no valid JSON found")
}
