package app

import (
	"context"

	"github.com/limes-cloud/kratosx/pkg/crypto"

	"github.com/limes-cloud/ai-agent/internal/core"
	"github.com/limes-cloud/kratosx/model"
	"github.com/limes-cloud/kratosx/pkg/value"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/limes-cloud/ai-agent/api/prompt"

	"github.com/limes-cloud/ai-agent/api/errors"
	"github.com/limes-cloud/ai-agent/internal/domain/entity"
	"github.com/limes-cloud/ai-agent/internal/domain/service"
	"github.com/limes-cloud/ai-agent/internal/infra/dbs"
	"github.com/limes-cloud/ai-agent/internal/types"
)

type Prompt struct {
	prompt.UnimplementedPromptServer
	srv *service.Prompt
}

func NewPrompt() *Prompt {
	return &Prompt{
		srv: service.NewPrompt(
			dbs.NewPrompt(),
		),
	}
}

func init() {
	register(func(hs *http.Server, gs *grpc.Server) {
		srv := NewPrompt()
		prompt.RegisterPromptServer(gs, srv)
		prompt.RegisterPromptHTTPServer(hs, srv)
	})
}

// ListPromptClassify 获取密钥分组列表（app 纬度隔离）
func (s *Prompt) ListPromptClassify(c context.Context, req *prompt.ListPromptClassifyRequest) (*prompt.ListPromptClassifyReply, error) {
	app, err := currentApp(c)
	if err != nil {
		return nil, err
	}
	list, total, err := s.srv.ListPromptClassify(core.MustContext(c), &types.ListPromptClassifyRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Name:     req.Name,
		App:      app,
	})
	if err != nil {
		return nil, err
	}
	reply := prompt.ListPromptClassifyReply{Total: total}
	for _, item := range list {
		reply.List = append(reply.List, &prompt.ListPromptClassifyReply_PromptClassify{
			Id:          item.Id,
			Name:        item.Name,
			Description: item.Description,
			CreatedAt:   uint32(item.CreatedAt),
			UpdatedAt:   uint32(item.UpdatedAt),
		})
	}
	return &reply, nil
}

// CreatePromptClassify 创建密钥分组（app 取自 X-App-Key）
func (s *Prompt) CreatePromptClassify(c context.Context, req *prompt.CreatePromptClassifyRequest) (*prompt.CreatePromptClassifyReply, error) {
	app, err := currentApp(c)
	if err != nil {
		return nil, err
	}
	id, err := s.srv.CreatePromptClassify(core.MustContext(c), &entity.PromptClassify{
		App:         app,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}
	return &prompt.CreatePromptClassifyReply{Id: id}, nil
}

// UpdatePromptClassify 更新密钥分组
func (s *Prompt) UpdatePromptClassify(c context.Context, req *prompt.UpdatePromptClassifyRequest) (*prompt.UpdatePromptClassifyReply, error) {
	if err := s.srv.UpdatePromptClassify(core.MustContext(c), &entity.PromptClassify{
		BaseModel:   model.BaseModel{Id: req.Id},
		Description: req.Description,
		Name:        req.Name,
	}); err != nil {
		return nil, err
	}
	return &prompt.UpdatePromptClassifyReply{}, nil
}

// DeletePromptClassify 删除密钥分组（app 纬度隔离）
func (s *Prompt) DeletePromptClassify(c context.Context, req *prompt.DeletePromptClassifyRequest) (*prompt.DeletePromptClassifyReply, error) {
	app, err := currentApp(c)
	if err != nil {
		return nil, err
	}
	err = s.srv.DeletePromptClassify(core.MustContext(c), app, req.Id)
	if err != nil {
		return nil, err
	}
	return &prompt.DeletePromptClassifyReply{}, nil
}

// ListPrompt 获取密钥信息列表（app 纬度隔离）
func (s *Prompt) ListPrompt(c context.Context, req *prompt.ListPromptRequest) (*prompt.ListPromptReply, error) {
	app, err := currentApp(c)
	if err != nil {
		return nil, err
	}
	result, total, err := s.srv.ListPrompt(core.MustContext(c), &types.ListPromptRequest{
		Page:       req.Page,
		PageSize:   req.PageSize,
		Title:      req.Title,
		ClassifyId: req.ClassifyId,
		Status:     req.Status,
		App:        app,
	})
	if err != nil {
		return nil, err
	}

	reply := prompt.ListPromptReply{Total: total}
	if err := value.Transform(result, &reply.List); err != nil {
		return nil, errors.TransformError()
	}
	return &reply, nil
}

// CreatePrompt 创建密钥信息（app 取自 X-App-Key）
func (s *Prompt) CreatePrompt(c context.Context, req *prompt.CreatePromptRequest) (*prompt.CreatePromptReply, error) {
	app, err := currentApp(c)
	if err != nil {
		return nil, err
	}
	id, err := s.srv.CreatePrompt(core.MustContext(c), &entity.Prompt{
		App:         app,
		ClassifyId:  req.ClassifyId,
		ModelId:     req.ModelId,
		UUID:        req.Uuid,
		Title:       req.Title,
		System:      req.System,
		Prompt:      req.Prompt,
		Description: req.Description,
		Extra:       req.Extra,
		Version:     crypto.MD5([]byte(req.Prompt)),
	})
	if err != nil {
		return nil, err
	}
	return &prompt.CreatePromptReply{Id: id}, nil
}

// UpdatePrompt 更新密钥信息
func (s *Prompt) UpdatePrompt(c context.Context, req *prompt.UpdatePromptRequest) (*prompt.UpdatePromptReply, error) {
	if err := s.srv.UpdatePrompt(core.MustContext(c), &entity.Prompt{
		BaseModel:   model.BaseModel{Id: req.Id},
		ClassifyId:  req.ClassifyId,
		ModelId:     req.ModelId,
		Title:       req.Title,
		System:      req.System,
		Prompt:      req.Prompt,
		Description: req.Description,
		Extra:       req.Extra,
		Status:      req.Status,
		Version:     crypto.MD5([]byte(req.Prompt)),
	}); err != nil {
		return nil, err
	}

	return &prompt.UpdatePromptReply{}, nil
}

// DeletePrompt 删除密钥信息（app 纬度隔离）
func (s *Prompt) DeletePrompt(c context.Context, req *prompt.DeletePromptRequest) (*prompt.DeletePromptReply, error) {
	app, err := currentApp(c)
	if err != nil {
		return nil, err
	}
	err = s.srv.DeletePrompt(core.MustContext(c), app, req.Id)
	if err != nil {
		return nil, err
	}
	return &prompt.DeletePromptReply{}, nil
}
