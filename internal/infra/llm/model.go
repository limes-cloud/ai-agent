package llm

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/limes-cloud/kratosx/library/request"

	"github.com/go-resty/resty/v2"

	"github.com/eino-contrib/jsonschema"

	json "github.com/json-iterator/go"

	"github.com/limes-cloud/ai-agent/internal/pkg"

	"github.com/cloudwego/eino/components/model"
	"github.com/limes-cloud/kratosx/pkg/value"

	"github.com/cloudwego/eino-ext/components/embedding/ollama"
	"github.com/cloudwego/eino-ext/components/embedding/openai"
	mollama "github.com/cloudwego/eino-ext/components/model/ollama"
	mopenai "github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
	"github.com/limes-cloud/ai-agent/internal/core"
	"github.com/limes-cloud/ai-agent/internal/domain/entity"
	"github.com/limes-cloud/ai-agent/internal/types"
)

var reg = regexp.MustCompile(`\{\{(\w|\.)+}}`)

type LLM struct {
	model  *entity.Model
	secret *entity.Secret
}

func NewLLM(model *entity.Model, secret *entity.Secret) *LLM {
	return &LLM{
		model:  model,
		secret: secret,
	}
}

// Embedding 嵌入向量
func (llm *LLM) Embedding(ctx core.Context, req *types.EmbeddingRequest) (*types.EmbeddingResponse, error) {
	if llm.model.Type != entity.ModelTypeEmbedding {
		return nil, fmt.Errorf("model type is not embedding")
	}

	switch llm.secret.Group.Supplier {
	case entity.SupplierOllama:
		embedder, err := ollama.NewEmbedder(ctx, &ollama.EmbeddingConfig{
			BaseURL: llm.secret.Group.Host,
			Model:   llm.model.Keyword, // 供应商 API 模型名（与 Text/Image 一致，勿用业务标识 uuid）
		})
		if err != nil {
			return nil, err
		}
		vectors, err := embedder.EmbedStrings(ctx, req.Texts)
		if err != nil {
			return nil, err
		}
		return &types.EmbeddingResponse{Vectors: vectors}, nil
	default:
		embedder, err := openai.NewEmbedder(ctx, &openai.EmbeddingConfig{
			// OpenAI API 配置
			APIKey:  llm.secret.Token,
			Model:   llm.model.Keyword, // 供应商 API 模型名（与 Text/Image 一致，勿用业务标识 uuid）
			BaseURL: llm.secret.Group.Host,
		})
		if err != nil {
			return nil, err
		}
		vectors, err := embedder.EmbedStrings(ctx, req.Texts)
		if err != nil {
			return nil, err
		}
		return &types.EmbeddingResponse{Vectors: vectors}, nil
	}
}

// prompt 替换模板中的变量
func (llm *LLM) replace(prompt string, variables map[string]any) (string, error) {
	keys := reg.FindAllString(prompt, -1)
	for _, key := range keys {
		oriKey := key[2 : len(key)-2]
		// 兼容{{.key}}语法
		oriKey = strings.ReplaceAll(oriKey, ".", "")
		if _, ok := variables[oriKey]; !ok {
			return "", fmt.Errorf("key %s not found in variables", key)
		}
		prompt = strings.ReplaceAll(prompt, key, fmt.Sprint(variables[oriKey]))
	}
	return prompt, nil
}

func (llm *LLM) message(req *types.TextRequest) ([]*schema.Message, error) {
	var msgs []*schema.Message
	// 追加系统提示
	if req.System != "" {
		system, err := llm.replace(req.System, req.Variables)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, &schema.Message{
			Role:    schema.System,
			Content: system,
		})
	}

	// 追加用户输入
	prompt, err := llm.replace(req.Prompt, req.Variables)
	if err != nil {
		return nil, err
	}
	// 判断是否支持JSON
	extra := llm.model.GetExtra()
	if !value.Value(extra.SupportJson) && req.ResponseFormat != nil {
		schemaByte, _ := json.Marshal(req.ResponseFormat.JSONSchema.JSONSchema)
		prompt += "# 输出格式要求\n 请严格按照json格式输出，不要有任何多余的内容。json中的key参考properties中的内容：" + string(schemaByte)
	}

	msg := &schema.Message{
		Role:    schema.User,
		Content: prompt,
	}

	// 追加附件
	if req.Attachment != nil {
		msg = &schema.Message{
			Role: schema.User,
			UserInputMultiContent: []schema.MessageInputPart{
				{Type: schema.ChatMessagePartTypeText, Text: prompt},
			},
		}
		for _, image := range req.Attachment.Images {
			msg.UserInputMultiContent = append(msg.UserInputMultiContent, schema.MessageInputPart{
				Type: schema.ChatMessagePartTypeImageURL,
				Image: &schema.MessageInputImage{
					MessagePartCommon: schema.MessagePartCommon{
						URL: value.Pointer(image),
					},
					Detail: schema.ImageURLDetailAuto,
				},
			})
		}
	}

	msgs = append(msgs, msg)
	return msgs, nil
}

// Text 文本生成
func (llm *LLM) Text(ctx core.Context, req *types.TextRequest) (*types.TextResponse, error) {
	if llm.model.Type != entity.ModelTypeChat {
		return nil, fmt.Errorf("model type is not chat")
	}
	// 处理消息
	msgList, err := llm.message(req)
	if err != nil {
		return nil, err
	}

	switch llm.secret.Group.Supplier {
	case entity.SupplierOllama:
		chatModel, err := mollama.NewChatModel(ctx, &mollama.ChatModelConfig{
			BaseURL: llm.secret.Group.Host,
			Model:   llm.model.Keyword, // 供应商 API 模型名（与 openai 分支一致，勿用业务标识 uuid）
		})
		if err != nil {
			return nil, err
		}

		var opts []model.Option
		if req.TopP != nil {
			opts = append(opts, model.WithTopP(value.Value(req.TopP)))
		}
		if req.Temperature != nil {
			opts = append(opts, model.WithTemperature(value.Value(req.Temperature)))
		}

		if req.Stream {
			sr, err := chatModel.Stream(ctx, msgList, opts...)
			if err != nil {
				return nil, err
			}
			text := ""
			for {
				msg, err := sr.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					return nil, err
				}
				text += msg.Content
				if req.RecvFunc != nil {
					req.RecvFunc(msg.Content)
				}
			}
			return &types.TextResponse{Content: text}, nil
		} else {
			msg, err := chatModel.Generate(ctx, msgList, opts...)
			if err != nil {
				return nil, err
			}
			return &types.TextResponse{Content: msg.Content}, nil
		}
	default:
		chatModel, err := mopenai.NewChatModel(ctx, &mopenai.ChatModelConfig{
			APIKey:         llm.secret.Token,
			Model:          llm.model.Keyword,
			BaseURL:        llm.secret.Group.Host,
			TopP:           req.TopP,
			Temperature:    req.Temperature,
			ResponseFormat: req.ResponseFormat,
		})
		if err != nil {
			return nil, err
		}

		if req.Stream {
			sr, err := chatModel.Stream(ctx, msgList)
			if err != nil {
				return nil, err
			}
			text := ""
			for {
				msg, err := sr.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					return nil, err
				}
				text += msg.Content
				if req.RecvFunc != nil {
					req.RecvFunc(msg.Content)
				}
			}
			return &types.TextResponse{Content: text}, nil
		} else {
			msg, err := chatModel.Generate(ctx, msgList)
			if err != nil {
				return nil, err
			}
			return &types.TextResponse{Content: msg.Content}, nil
		}
	}
}

// TextScan 文本生成并解析
func (llm *LLM) TextScan(ctx core.Context, req *types.TextRequest, dst any) error {
	// Scan模式强制关闭流模式
	req.Stream = false
	// 如果未指定返回格式，则自动生成json schema
	if req.ResponseFormat == nil {
		schemaByte := pkg.Schema(dst)
		sm := jsonschema.Schema{}
		if err := json.Unmarshal(schemaByte, &sm); err != nil {
			return err
		}
		req.ResponseFormat = &mopenai.ChatCompletionResponseFormat{
			Type: mopenai.ChatCompletionResponseFormatTypeJSONSchema,
			JSONSchema: &mopenai.ChatCompletionResponseFormatJSONSchema{
				Name:       "object",
				JSONSchema: &sm,
				Strict:     true,
			},
		}
	}
	reply, err := llm.Text(ctx, req)
	if err != nil {
		return err
	}

	// 去除多余的描述
	text := reply.Content
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimSuffix(text, "```")

	// 尝试直接解析
	if err := json.Unmarshal([]byte(text), dst); err == nil {
		return nil
	}

	// 尝试解析提取
	text, err = pkg.MatchExtractJSON(text)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(text), dst)
}

// Image 图片生成
func (llm *LLM) Image(ctx core.Context, req *types.ImageRequest) (*types.ImageResponse, error) {
	if llm.model.Type != entity.ModelTypeImage {
		return nil, errors.New("model does not support image generation")
	}

	if req.Size == "" {
		req.Size = "1024x1024"
	}

	// 格式化输入
	prompt, err := llm.replace(req.Prompt, req.Variables)
	if err != nil {
		return nil, err
	}

	baseURL := strings.TrimSuffix(llm.secret.Group.Host, "/")

	switch llm.secret.Group.Supplier {
	case entity.SupplierQianfan, entity.SupplierARKBot:
		if strings.Contains(llm.model.Keyword, "musesteamer") {
			baseURL += "/musesteamer"
		}
		response, err := ctx.Request(request.WithTimeout(1*time.Minute)).Option(func(request *resty.Request) {
			request.Header.Set("Authorization", "Bearer "+llm.secret.Token)
		}).Post(baseURL+"/images/generations", map[string]any{
			"model":  llm.model.Keyword,
			"prompt": prompt,
			"size":   req.Size,
		})
		if err != nil {
			return nil, err
		}

		reply := struct {
			Error struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
			Id      string `json:"id"`
			Created int    `json:"created"`
			Data    []struct {
				Url string `json:"url"`
			} `json:"data"`
		}{}
		if err := response.Result(&reply); err != nil {
			return nil, err
		}
		if reply.Error.Code != "" {
			return nil, fmt.Errorf("image generation error: %s", reply.Error.Message)
		}
		var list []string
		for _, v := range reply.Data {
			list = append(list, v.Url)
		}
		return &types.ImageResponse{
			List: list,
		}, nil

	case entity.SupplierQwen:
		// 千问AI平台（qianwenai）：百炼原生多模态生成（同步），非 OpenAI 兼容。
		// host 配置为 https://maas.qianwenaiapi.com/api/v1，协议：
		//   POST {host}/services/aigc/multimodal-generation/generation
		//   {"model","messages":[{"role":"user","content":[{"text"}]}],"parameters":{"size":"1024*1024"}}
		//   响应：output.choices[0].message.content[0].image（size 用 * 分隔，与 OpenAI 的 x 不同）
		response, err := ctx.Request(request.WithTimeout(2*time.Minute)).Option(func(request *resty.Request) {
			request.Header.Set("Authorization", "Bearer "+llm.secret.Token)
		}).Post(baseURL+"/services/aigc/multimodal-generation/generation", map[string]any{
			"model": llm.model.Keyword,
			"messages": []map[string]any{{
				"role":    "user",
				"content": []map[string]any{{"text": prompt}},
			}},
			"parameters": map[string]any{
				"size": strings.ReplaceAll(req.Size, "x", "*"),
				"n":    1,
			},
		})
		if err != nil {
			return nil, err
		}

		reply := struct {
			Code    string `json:"code"`
			Message string `json:"message"`
			Output  struct {
				Choices []struct {
					Message struct {
						Content []struct {
							Image string `json:"image"`
						} `json:"content"`
					} `json:"message"`
				} `json:"choices"`
			} `json:"output"`
		}{}
		if err := response.Result(&reply); err != nil {
			return nil, err
		}
		if reply.Code != "" {
			return nil, fmt.Errorf("image generation error: %s", reply.Message)
		}
		var list []string
		for _, c := range reply.Output.Choices {
			for _, v := range c.Message.Content {
				if v.Image != "" {
					list = append(list, v.Image)
				}
			}
		}
		return &types.ImageResponse{List: list}, nil
	}
	return nil, errors.New("model does not support image generation")
}
