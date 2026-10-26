package service

import (
	"encoding/json"
	"time"

	mopenai "github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/eino-contrib/jsonschema"

	"github.com/google/uuid"

	"github.com/limes-cloud/kratosx/pkg/value"

	"github.com/limes-cloud/resource/api/file"

	"github.com/limes-cloud/ai-agent/api/errors"
	"github.com/limes-cloud/ai-agent/internal/core"
	"github.com/limes-cloud/ai-agent/internal/domain/entity"
	"github.com/limes-cloud/ai-agent/internal/domain/repository"
	"github.com/limes-cloud/ai-agent/internal/infra/llm"
	"github.com/limes-cloud/ai-agent/internal/types"
)

type Generate struct {
	model    repository.Model
	secret   repository.Secret
	prompt   repository.Prompt
	resource repository.Resource
	callLog  repository.CallLog
}

func NewGenerate(
	model repository.Model,
	secret repository.Secret,
	prompt repository.Prompt,
	resource repository.Resource,
	callLog repository.CallLog,
) *Generate {
	return &Generate{
		model:    model,
		secret:   secret,
		prompt:   prompt,
		resource: resource,
		callLog:  callLog,
	}
}

// base 按模型 uuid 获取模型与密钥；模型被禁用时拒绝。
func (t *Generate) base(ctx core.Context, uuid string) (*entity.Model, *entity.Secret, error) {
	// 模型是全局资源，不按 app 隔离
	model, err := t.model.GetModelByUUID(ctx, uuid)
	if err != nil {
		return nil, nil, err
	}
	if value.Value(model.Status) == 0 {
		return nil, nil, errors.SystemError("模型已禁用")
	}

	// 获取密钥信息（全局密钥组，随机取启用密钥）
	secret, err := t.secret.GetRandomSecretByGroupId(ctx, model.SecretGroupId)
	if err != nil {
		return nil, nil, err
	}
	return model, secret, nil
}

// recordCallLog 记录一次生成调用（按 app 计量；失败不阻断主流程）。
func (t *Generate) recordCallLog(ctx core.Context, app string, model *entity.Model, supplier, scene string, start time.Time, status uint8, errMsg string) {
	if app == "" || model == nil {
		return
	}
	_, err := t.callLog.CreateCallLog(ctx, &entity.CallLog{
		App:          app,
		ModelId:      model.Id,
		ModelUUID:    model.UUID,
		Supplier:     supplier,
		Scene:        scene,
		LatencyMs:    time.Since(start).Milliseconds(),
		Status:       status,
		ErrorMessage: errMsg,
	})
	if err != nil {
		ctx.Logger().Warnw("msg", "create call log error", "app", app, "err", err.Error())
	}
}

// Text 文本生产
func (t *Generate) Text(ctx core.Context, req *types.TextRequest) (reply *types.TextResponse, err error) {
	start := time.Now()
	var lastModel *entity.Model
	defer func() {
		status, msg := uint8(1), ""
		if err != nil {
			status, msg = 0, err.Error()
		}
		if len(msg) > 1024 {
			msg = msg[:1024]
		}
		supplier := ""
		if lastModel != nil && lastModel.SecretGroup != nil {
			supplier = lastModel.SecretGroup.Supplier
		}
		t.recordCallLog(ctx, req.App, lastModel, supplier, "text", start, status, msg)
	}()

	model, secret, berr := t.base(ctx, req.Model)
	if berr != nil {
		return nil, berr
	}
	lastModel = model

	lc := llm.NewLLM(model, secret)

	if req.ResponseFormat != nil {
		dst := map[string]any{}
		if serr := lc.TextScan(ctx, req, &dst); serr != nil {
			return nil, errors.SystemError(serr.Error())
		}
		r, mErr := json.Marshal(dst)
		if mErr != nil {
			return nil, errors.SystemError(mErr.Error())
		}
		return &types.TextResponse{Content: string(r), Format: dst}, nil
	}
	res, terr := lc.Text(ctx, req)
	if terr != nil {
		return nil, errors.SystemError(terr.Error())
	}
	return &types.TextResponse{Content: res.Content}, nil
}

// Image 图片生产（生成后上传资源中心，返回 key 列表）
func (t *Generate) Image(ctx core.Context, req *types.ImageRequest) (reply *types.ImageResponse, err error) {
	start := time.Now()
	var lastModel *entity.Model
	defer func() {
		status, msg := uint8(1), ""
		if err != nil {
			status, msg = 0, err.Error()
		}
		if len(msg) > 1024 {
			msg = msg[:1024]
		}
		supplier := ""
		if lastModel != nil && lastModel.SecretGroup != nil {
			supplier = lastModel.SecretGroup.Supplier
		}
		t.recordCallLog(ctx, req.App, lastModel, supplier, "image", start, status, msg)
	}()

	model, secret, berr := t.base(ctx, req.Model)
	if berr != nil {
		return nil, berr
	}
	lastModel = model

	lc := llm.NewLLM(model, secret)
	res, gerr := lc.Image(ctx, req)
	if gerr != nil {
		return nil, errors.SystemError(gerr.Error())
	}
	if len(res.List) == 0 {
		return nil, errors.SystemError("图片生成失败")
	}

	// 存储到资源中心（X-App-Key 透传，按 app 的 AppConfig 落桶/计量）；返回对象 key（内容指纹，持久可引用）
	var list []string
	for _, v := range res.List {
		up, uerr := t.resource.UploadFileByURL(ctx, &file.UploadFileByURLRequest{
			DirectoryPath: value.Pointer(req.Path),
			Url:           v,
			Name:          uuid.New().String() + ".png",
			Store:         req.Store,
		})
		if uerr != nil {
			return nil, uerr
		}
		list = append(list, up.Key)
	}
	return &types.ImageResponse{List: list}, nil
}

// PromptText 文本生产（按提示词模板与 Extra 配置生成）
func (t *Generate) PromptText(ctx core.Context, req *types.PromptTextRequest) (reply *types.PromptTextResponse, err error) {
	start := time.Now()
	var lastModel *entity.Model
	defer func() {
		status, msg := uint8(1), ""
		if err != nil {
			status, msg = 0, err.Error()
		}
		if len(msg) > 1024 {
			msg = msg[:1024]
		}
		supplier := ""
		if lastModel != nil && lastModel.SecretGroup != nil {
			supplier = lastModel.SecretGroup.Supplier
		}
		t.recordCallLog(ctx, req.App, lastModel, supplier, "prompt_text", start, status, msg)
	}()

	// 获取提示词相关信息（app 纬度隔离）
	prompt, perr := t.prompt.GetPromptByUUID(ctx, req.App, req.Prompt)
	if perr != nil {
		return nil, errors.SystemError()
	}
	if value.Value(prompt.Status) != 1 {
		return nil, errors.PromptDisableError()
	}

	// 获取模型和密钥
	model, secret, berr := t.base(ctx, prompt.Model.UUID)
	if berr != nil {
		return nil, berr
	}
	lastModel = model

	// 初始化客户端
	lc := llm.NewLLM(model, secret)
	extra := prompt.GetExtra()

	params := &types.TextRequest{
		App:         req.App,
		Model:       model.Keyword,
		System:      value.Value(prompt.System),
		Prompt:      prompt.Prompt,
		Variables:   req.Variables,
		Temperature: extra.Temperature,
		TopP:        extra.TopP,
	}

	if extra.FormatJson != nil && *extra.FormatJson {
		sm := jsonschema.Schema{}
		if uerr := json.Unmarshal([]byte(value.Value(extra.OutputSchema)), &sm); uerr != nil {
			return nil, errors.ParamsError("responseFormat is invalid")
		}
		params.ResponseFormat = &mopenai.ChatCompletionResponseFormat{
			Type: mopenai.ChatCompletionResponseFormatTypeJSONSchema,
			JSONSchema: &mopenai.ChatCompletionResponseFormatJSONSchema{
				Name:       "object",
				JSONSchema: &sm,
				Strict:     true,
			},
		}

		dst := map[string]any{}
		if serr := lc.TextScan(ctx, params, &dst); serr != nil {
			return nil, errors.SystemError(serr.Error())
		}
		r, mErr := json.Marshal(dst)
		if mErr != nil {
			return nil, errors.SystemError(mErr.Error())
		}
		return &types.PromptTextResponse{Content: string(r), Format: dst}, nil
	}

	res, terr := lc.Text(ctx, params)
	if terr != nil {
		return nil, errors.SystemError(terr.Error())
	}
	return &types.PromptTextResponse{Content: res.Content}, nil
}

// PromptImage 图片生产（按提示词模板与 Extra 配置生成）
func (t *Generate) PromptImage(ctx core.Context, req *types.PromptImageRequest) (reply *types.PromptImageResponse, err error) {
	start := time.Now()
	var lastModel *entity.Model
	defer func() {
		status, msg := uint8(1), ""
		if err != nil {
			status, msg = 0, err.Error()
		}
		if len(msg) > 1024 {
			msg = msg[:1024]
		}
		supplier := ""
		if lastModel != nil && lastModel.SecretGroup != nil {
			supplier = lastModel.SecretGroup.Supplier
		}
		t.recordCallLog(ctx, req.App, lastModel, supplier, "prompt_image", start, status, msg)
	}()

	// 获取提示词相关信息（app 纬度隔离）
	prompt, perr := t.prompt.GetPromptByUUID(ctx, req.App, req.Prompt)
	if perr != nil {
		return nil, errors.SystemError()
	}

	// 获取模型和密钥
	model, secret, berr := t.base(ctx, prompt.Model.UUID)
	if berr != nil {
		return nil, berr
	}
	lastModel = model

	// 初始化客户端
	lc := llm.NewLLM(model, secret)
	extra := prompt.GetExtra()

	params := &types.ImageRequest{
		App:       req.App,
		Model:     model.Keyword,
		Prompt:    prompt.Prompt,
		Variables: req.Variables,
		Size:      extra.Size,
		Path:      extra.Path,
		Store:     extra.Store,
	}

	res, gerr := lc.Image(ctx, params)
	if gerr != nil {
		return nil, errors.SystemError(gerr.Error())
	}
	if len(res.List) == 0 {
		return nil, errors.SystemError("图片生成失败")
	}

	// 存储到资源中心；返回对象 key（内容指纹，持久可引用）
	var list []string
	for _, v := range res.List {
		up, uerr := t.resource.UploadFileByURL(ctx, &file.UploadFileByURLRequest{
			DirectoryPath: value.Pointer(params.Path),
			Url:           v,
			Name:          uuid.New().String() + ".png",
			Store:         params.Store,
		})
		if uerr != nil {
			return nil, uerr
		}
		list = append(list, up.Key)
	}

	return &types.PromptImageResponse{List: list}, nil
}