package app

import (
	"context"

	"github.com/limes-cloud/ai-agent/internal/infra/rpcs"

	"google.golang.org/protobuf/types/known/structpb"

	mopenai "github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/eino-contrib/jsonschema"
	json "github.com/json-iterator/go"
	"github.com/limes-cloud/ai-agent/internal/types"
	"github.com/limes-cloud/ai-agent/api/errors"

	"github.com/limes-cloud/ai-agent/internal/core"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/limes-cloud/ai-agent/api/generate"
	"github.com/limes-cloud/ai-agent/internal/domain/service"
	"github.com/limes-cloud/ai-agent/internal/infra/dbs"
)

type Generate struct {
	generate.UnimplementedGenerateServer
	srv *service.Generate
}

func NewGenerate() *Generate {
	return &Generate{
		srv: service.NewGenerate(
			dbs.NewModel(),
			dbs.NewSecret(),
			dbs.NewPrompt(),
			rpcs.NewResource(),
			dbs.NewCallLog(),
		),
	}
}

func init() {
	register(func(hs *http.Server, gs *grpc.Server) {
		srv := NewGenerate()
		generate.RegisterGenerateServer(gs, srv)
		generate.RegisterGenerateHTTPServer(hs, srv)
	})
}

// Text 文本生产（app 取自 X-App-Key）
func (s *Generate) Text(ctx context.Context, req *generate.TextRequest) (*generate.TextReply, error) {
	app, err := currentApp(ctx)
	if err != nil {
		return nil, err
	}
	params := &types.TextRequest{
		App:         app,
		Model:       req.Model,
		System:      req.System,
		Prompt:      req.Prompt,
		Stream:      req.Stream,
		Temperature: req.Temperature,
		TopP:        req.TopP,
	}
	if req.ResponseFormat != nil {
		sm := jsonschema.Schema{}
		rf, _ := req.ResponseFormat.MarshalJSON()
		if err := json.Unmarshal(rf, &sm); err != nil {
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
	}
	if req.Variables != nil {
		params.Variables = req.Variables.AsMap()
	}
	if req.Attachment != nil {
		params.Attachment = &types.TextAttachment{
			Images: req.Attachment.Images,
		}
	}
	reply, err := s.srv.Text(core.MustContext(ctx), params)
	if err != nil {
		return nil, errors.SystemError(err.Error())
	}
	ft, _ := structpb.NewStruct(reply.Format)
	return &generate.TextReply{Content: reply.Content, Format: ft}, nil
}

// Image 图片生产（app 取自 X-App-Key）
func (s *Generate) Image(ctx context.Context, req *generate.ImageRequest) (*generate.ImageReply, error) {
	app, err := currentApp(ctx)
	if err != nil {
		return nil, err
	}
	params := &types.ImageRequest{
		App:    app,
		Model:  req.Model,
		Prompt: req.Prompt,
		Size:   req.Size,
		Store:  req.Store,
		Path:   req.DirectoryPath,
	}
	if req.Variables != nil {
		params.Variables = req.Variables.AsMap()
	}

	reply, err := s.srv.Image(core.MustContext(ctx), params)
	if err != nil {
		return nil, errors.SystemError(err.Error())
	}
	return &generate.ImageReply{List: reply.List}, nil
}

// PromptText 文本生产（按提示词模板；app 取自 X-App-Key）
func (s *Generate) PromptText(ctx context.Context, req *generate.PromptTextRequest) (*generate.PromptTextReply, error) {
	app, err := currentApp(ctx)
	if err != nil {
		return nil, err
	}
	params := &types.PromptTextRequest{
		App:    app,
		Prompt: req.Prompt,
	}
	if req.Variables != nil {
		params.Variables = req.Variables.AsMap()
	}
	if req.Attachment != nil {
		params.Attachment = &types.TextAttachment{
			Images: req.Attachment.Images,
		}
	}

	reply, err := s.srv.PromptText(core.MustContext(ctx), params)
	if err != nil {
		return nil, err
	}
	ft, _ := structpb.NewStruct(reply.Format)
	return &generate.PromptTextReply{Content: reply.Content, Format: ft}, nil
}

// PromptImage 图片生产（按提示词模板；app 取自 X-App-Key）
func (s *Generate) PromptImage(ctx context.Context, req *generate.PromptImageRequest) (*generate.PromptImageReply, error) {
	app, err := currentApp(ctx)
	if err != nil {
		return nil, err
	}
	params := &types.PromptImageRequest{
		App:    app,
		Prompt: req.Prompt,
	}
	if req.Variables != nil {
		params.Variables = req.Variables.AsMap()
	}
	reply, err := s.srv.PromptImage(core.MustContext(ctx), params)
	if err != nil {
		return nil, err
	}
	return &generate.PromptImageReply{List: reply.List}, nil
}
