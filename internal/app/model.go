package app

import (
	"context"

	"github.com/limes-cloud/kratosx/pkg/value"

	"github.com/limes-cloud/ai-agent/internal/core"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/limes-cloud/ai-agent/api/model"
	kmodel "github.com/limes-cloud/kratosx/model"

	"github.com/limes-cloud/ai-agent/api/errors"
	"github.com/limes-cloud/ai-agent/internal/domain/entity"
	"github.com/limes-cloud/ai-agent/internal/domain/service"
	"github.com/limes-cloud/ai-agent/internal/infra/dbs"
	"github.com/limes-cloud/ai-agent/internal/types"
)

type Model struct {
	model.UnimplementedModelServer
	srv *service.Model
}

func NewModel() *Model {
	return &Model{
		srv: service.NewModel(
			dbs.NewModel(),
		),
	}
}

func init() {
	register(func(hs *http.Server, gs *grpc.Server) {
		srv := NewModel()
		model.RegisterModelServer(gs, srv)
		model.RegisterModelHTTPServer(hs, srv)
	})
}

// ListModelType 获取密钥分组列表
func (s *Model) ListModelType(_ context.Context, _ *model.ListModelTypeRequest) (*model.ListModelTypeReply, error) {
	list := s.srv.ListModelType()
	return &model.ListModelTypeReply{List: list}, nil
}

// ListModel 获取模型信息列表
func (s *Model) ListModel(c context.Context, req *model.ListModelRequest) (*model.ListModelReply, error) {
	result, total, err := s.srv.ListModel(core.MustContext(c), &types.ListModelRequest{
		Page:          req.Page,
		PageSize:      req.PageSize,
		Type:          req.Type,
		UUID:          req.Uuid,
		SecretGroupId: req.SecretGroupId,
	})
	if err != nil {
		return nil, err
	}

	reply := model.ListModelReply{Total: total}
	if err := value.Transform(result, &reply.List); err != nil {
		return nil, errors.TransformError()
	}
	return &reply, nil
}

// CreateModel 创建模型信息
func (s *Model) CreateModel(c context.Context, req *model.CreateModelRequest) (*model.CreateModelReply, error) {
	id, err := s.srv.CreateModel(core.MustContext(c), &entity.Model{
		Type:          req.Type,
		UUID:          req.Uuid,
		Keyword:       req.Keyword,
		Extra:         req.Extra,
		SecretGroupId: req.SecretGroupId,
		Description:   req.Description,
		Status:        req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &model.CreateModelReply{Id: id}, nil
}

// UpdateModel 更新模型信息
func (s *Model) UpdateModel(c context.Context, req *model.UpdateModelRequest) (*model.UpdateModelReply, error) {
	if err := s.srv.UpdateModel(core.MustContext(c), &entity.Model{
		BaseModel:     kmodel.BaseModel{Id: req.Id},
		Type:          req.Type,
		UUID:          req.Uuid,
		Keyword:       req.Keyword,
		Extra:         req.Extra,
		SecretGroupId: req.SecretGroupId,
		Description:   req.Description,
		Status:        req.Status,
	}); err != nil {
		return nil, err
	}

	return &model.UpdateModelReply{}, nil
}

// DeleteModel 删除模型信息
func (s *Model) DeleteModel(c context.Context, req *model.DeleteModelRequest) (*model.DeleteModelReply, error) {
	if err := s.srv.DeleteModel(core.MustContext(c), req.Id); err != nil {
		return nil, err
	}
	return &model.DeleteModelReply{}, nil
}
