package menu

import (
	"github.com/limes-cloud/kratosx/pkg/value"
	managerMenu "github.com/limes-cloud/manager/api/menu"
)

const (
	appKeyword = "ai-agent"

	menuTypeApi   = "A"
	menuTypePage  = "M"
	menuTypeGroup = "G"
)

var menus = []*managerMenu.RegisterMenuNode{
	{
		Key:      value.Pointer("AIAgentBaseApi"),
		Title:    "AI中台基础接口",
		Type:     menuTypeGroup,
		Icon:     value.Pointer("apps"),
		IsHidden: value.Pointer(true),
		Weight:   value.Pointer[int32](40),
		Children: []*managerMenu.RegisterMenuNode{
			api("文本生成", "/ai-agent/api/generate/text", "POST", ""),
			api("图片生成", "/ai-agent/api/generate/image", "POST", ""),
			api("提示词文本生成", "/ai-agent/api/generate/prompt/text", "POST", ""),
			api("提示词图片生成", "/ai-agent/api/generate/prompt/image", "POST", ""),
			api("查询模型类型", "/ai-agent/api/model/types", "GET", ""),
			api("查询密钥厂商", "/ai-agent/api/secret/suppliers", "GET", ""),
		},
	},
	{
		Key:       value.Pointer("AIAgentSecret"),
		Title:     "密钥管理",
		Type:      menuTypePage,
		Icon:      value.Pointer("iconsys-anquan_2"),
		Path:      value.Pointer("/ai-agent/secret"),
		Component: value.Pointer("/ai-agent/secret"),
		Weight:    value.Pointer[int32](30),
		IsHidden:  value.Pointer(false),
		IsCache:   value.Pointer(true),
		IsAffix:   value.Pointer(true),
		Children: []*managerMenu.RegisterMenuNode{
			group("密钥分组", "ai-agent:secret-group:group", []*managerMenu.RegisterMenuNode{
				api("查询密钥分组", "/ai-agent/api/secret/groups", "GET", "ai-agent:secret:group:query"),
				api("新增密钥分组", "/ai-agent/api/secret/group", "POST", "ai-agent:secret:group:add"),
				api("修改密钥分组", "/ai-agent/api/secret/group", "PUT", "ai-agent:secret:group:update"),
				api("删除密钥分组", "/ai-agent/api/secret/group", "DELETE", "ai-agent:secret:group:delete"),
			}),
			group("密钥", "ai-agent:secret:group", []*managerMenu.RegisterMenuNode{
				api("查询密钥", "/ai-agent/api/secrets", "GET", "ai-agent:secret:query"),
				api("新增密钥", "/ai-agent/api/secret", "POST", "ai-agent:secret:add"),
				api("修改密钥", "/ai-agent/api/secret", "PUT", "ai-agent:secret:update"),
				api("删除密钥", "/ai-agent/api/secret", "DELETE", "ai-agent:secret:delete"),
			}),
		},
	},
	{
		// 独立顶级菜单（不从属于密钥管理页）：/ai-agent/secret/group 需要作为一级路由渲染，
		// 若作为 AIAgentSecret 的子路由，父组件（index.vue）无 <router-view> 会导致子路由渲染时停留在父组件
		Key:       value.Pointer("AIAgentSecretGroup"),
		Title:     "密钥分组",
		Type:      menuTypePage,
		Path:      value.Pointer("/ai-agent/secret/group"),
		Component: value.Pointer("/ai-agent/secret/group"),
		IsHidden:  value.Pointer(true),
		Weight:    value.Pointer[int32](29),
	},
	{
		Key:       value.Pointer("AIAgentModel"),
		Title:     "模型管理",
		Type:      menuTypePage,
		Icon:      value.Pointer("iconsys-code-block"),
		Path:      value.Pointer("/ai-agent/model"),
		Component: value.Pointer("/ai-agent/model"),
		Weight:    value.Pointer[int32](20),
		IsHidden:  value.Pointer(false),
		IsCache:   value.Pointer(true),
		IsAffix:   value.Pointer(true),
		Children: []*managerMenu.RegisterMenuNode{
			api("查询模型", "/ai-agent/api/models", "GET", "ai-agent:model:query"),
			api("新增模型", "/ai-agent/api/model", "POST", "ai-agent:model:add"),
			api("修改模型", "/ai-agent/api/model", "PUT", "ai-agent:model:update"),
			api("删除模型", "/ai-agent/api/model", "DELETE", "ai-agent:model:delete"),
		},
	},
	{
		Key:       value.Pointer("AIAgentPrompt"),
		Title:     "Prompt管理",
		Type:      menuTypePage,
		Icon:      value.Pointer("iconsys-liebiaoxingshi"),
		Path:      value.Pointer("/ai-agent/prompt"),
		Component: value.Pointer("/ai-agent/prompt"),
		Weight:    value.Pointer[int32](10),
		IsHidden:  value.Pointer(false),
		IsCache:   value.Pointer(true),
		IsAffix:   value.Pointer(true),
		Children: []*managerMenu.RegisterMenuNode{
			group("提示词分类", "ai-agent:prompt-classify:group", []*managerMenu.RegisterMenuNode{
				api("查询提示词分类", "/ai-agent/api/prompt/classifies", "GET", "ai-agent:prompt:classify:query"),
				api("新增提示词分类", "/ai-agent/api/prompt/classify", "POST", "ai-agent:prompt:classify:add"),
				api("修改提示词分类", "/ai-agent/api/prompt/classify", "PUT", "ai-agent:prompt:classify:update"),
				api("删除提示词分类", "/ai-agent/api/prompt/classify", "DELETE", "ai-agent:prompt:classify:delete"),
			}),
			group("提示词", "ai-agent:prompt:group", []*managerMenu.RegisterMenuNode{
				api("查询提示词", "/ai-agent/api/prompts", "GET", "ai-agent:prompt:query"),
				api("新增提示词", "/ai-agent/api/prompt", "POST", "ai-agent:prompt:add"),
				api("修改提示词", "/ai-agent/api/prompt", "PUT", "ai-agent:prompt:update"),
				api("删除提示词", "/ai-agent/api/prompt", "DELETE", "ai-agent:prompt:delete"),
			}),
		},
	},
	{
		// 独立顶级菜单（不从属于提示词管理页）：/ai-agent/prompt/classify 需要作为一级路由渲染，
		// 若作为 AIAgentPrompt 的子路由，父组件（index.vue）无 <router-view> 会导致子路由渲染时停留在父组件
		Key:       value.Pointer("AIAgentPromptClassify"),
		Title:     "提示词分类",
		Type:      menuTypePage,
		Path:      value.Pointer("/ai-agent/prompt/classify"),
		Component: value.Pointer("/ai-agent/prompt/classify"),
		IsHidden:  value.Pointer(true),
		Weight:    value.Pointer[int32](9),
	},
}

func group(title string, permission string, children []*managerMenu.RegisterMenuNode) *managerMenu.RegisterMenuNode {
	return &managerMenu.RegisterMenuNode{
		Title:      title,
		Type:       menuTypeGroup,
		Permission: value.Pointer(permission),
		Weight:     value.Pointer[int32](0),
		Children:   children,
	}
}

func api(title string, path string, method string, permission string) *managerMenu.RegisterMenuNode {
	node := &managerMenu.RegisterMenuNode{
		Title:  title,
		Type:   menuTypeApi,
		Api:    value.Pointer(path),
		Method: value.Pointer(method),
		Weight: value.Pointer[int32](0),
	}
	if permission != "" {
		node.Permission = value.Pointer(permission)
	}
	return node
}
