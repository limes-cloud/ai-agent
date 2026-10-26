package app

import (
	"context"

	"github.com/limes-cloud/ai-agent/internal/core"
	"github.com/limes-cloud/kratosx/model"
	"github.com/limes-cloud/kratosx/pkg/value"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/limes-cloud/ai-agent/api/secret"

	"github.com/limes-cloud/ai-agent/api/errors"
	"github.com/limes-cloud/ai-agent/internal/domain/entity"
	"github.com/limes-cloud/ai-agent/internal/domain/service"
	"github.com/limes-cloud/ai-agent/internal/infra/dbs"
	"github.com/limes-cloud/ai-agent/internal/types"
)

type Secret struct {
	secret.UnimplementedSecretServer
	srv *service.Secret
}

func NewSecret() *Secret {
	return &Secret{
		srv: service.NewSecret(
			dbs.NewSecret(),
		),
	}
}

func init() {
	register(func(hs *http.Server, gs *grpc.Server) {
		srv := NewSecret()
		secret.RegisterSecretServer(gs, srv)
		secret.RegisterSecretHTTPServer(hs, srv)
	})
}

// ListSupplier 获取密钥分组列表
func (s *Secret) ListSupplier(_ context.Context, _ *secret.ListSupplierRequest) (*secret.ListSupplierReply, error) {
	list := s.srv.ListSupplier()
	return &secret.ListSupplierReply{List: list}, nil
}

// ListSecretGroup 获取密钥分组列表
func (s *Secret) ListSecretGroup(c context.Context, req *secret.ListSecretGroupRequest) (*secret.ListSecretGroupReply, error) {
	list, total, err := s.srv.ListSecretGroup(core.MustContext(c), &types.ListSecretGroupRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Name:     req.Name,
		Supplier: req.Supplier,
	})
	if err != nil {
		return nil, err
	}
	reply := secret.ListSecretGroupReply{Total: total}
	for _, item := range list {
		reply.List = append(reply.List, &secret.ListSecretGroupReply_SecretGroup{
			Id:        item.Id,
			Name:      item.Name,
			Host:      item.Host,
			Supplier:  item.Supplier,
			CreatedAt: uint32(item.CreatedAt),
			UpdatedAt: uint32(item.UpdatedAt),
		})
	}
	return &reply, nil
}

// CreateSecretGroup 创建密钥分组
func (s *Secret) CreateSecretGroup(c context.Context, req *secret.CreateSecretGroupRequest) (*secret.CreateSecretGroupReply, error) {
	id, err := s.srv.CreateSecretGroup(core.MustContext(c), &entity.SecretGroup{
		Supplier: req.Supplier,
		Name:     req.Name,
		Host:     req.Host,
	})
	if err != nil {
		return nil, err
	}
	return &secret.CreateSecretGroupReply{Id: id}, nil
}

// UpdateSecretGroup 更新密钥分组
func (s *Secret) UpdateSecretGroup(c context.Context, req *secret.UpdateSecretGroupRequest) (*secret.UpdateSecretGroupReply, error) {
	if err := s.srv.UpdateSecretGroup(core.MustContext(c), &entity.SecretGroup{
		BaseModel: model.BaseModel{Id: req.Id},
		Supplier:  req.Supplier,
		Name:      req.Name,
		Host:      req.Host,
	}); err != nil {
		return nil, err
	}
	return &secret.UpdateSecretGroupReply{}, nil
}

// DeleteSecretGroup 删除密钥分组
func (s *Secret) DeleteSecretGroup(c context.Context, req *secret.DeleteSecretGroupRequest) (*secret.DeleteSecretGroupReply, error) {
	err := s.srv.DeleteSecretGroup(core.MustContext(c), req.Id)
	if err != nil {
		return nil, err
	}
	return &secret.DeleteSecretGroupReply{}, nil
}

// ListSecret 获取密钥信息列表（token 脱敏，仅尾 4 位）
func (s *Secret) ListSecret(c context.Context, req *secret.ListSecretRequest) (*secret.ListSecretReply, error) {
	result, total, err := s.srv.ListSecret(core.MustContext(c), &types.ListSecretRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Name:     req.Name,
		GroupId:  req.GroupId,
		Status:   req.Status,
	})
	if err != nil {
		return nil, err
	}

	reply := secret.ListSecretReply{Total: total}
	if err := value.Transform(result, &reply.List); err != nil {
		return nil, errors.TransformError()
	}
	// 脱敏：仅保留尾 4 位
	for _, item := range reply.List {
		if item != nil {
			item.Token = maskToken(item.Token)
		}
	}
	return &reply, nil
}

// CreateSecret 创建密钥信息
func (s *Secret) CreateSecret(c context.Context, req *secret.CreateSecretRequest) (*secret.CreateSecretReply, error) {
	id, err := s.srv.CreateSecret(core.MustContext(c), &entity.Secret{
		GroupId: req.GroupId,
		Token:   req.Token,
		Name:    req.Name,
	})
	if err != nil {
		return nil, err
	}
	return &secret.CreateSecretReply{Id: id}, nil
}

// UpdateSecret 更新密钥信息
func (s *Secret) UpdateSecret(c context.Context, req *secret.UpdateSecretRequest) (*secret.UpdateSecretReply, error) {
	if err := s.srv.UpdateSecret(core.MustContext(c), &entity.Secret{
		BaseModel: model.BaseModel{Id: req.Id},
		GroupId:   req.GroupId,
		Token:     req.Token,
		Status:    req.Status,
		Name:      req.Name,
	}); err != nil {
		return nil, err
	}

	return &secret.UpdateSecretReply{}, nil
}

// DeleteSecret 删除密钥信息
func (s *Secret) DeleteSecret(c context.Context, req *secret.DeleteSecretRequest) (*secret.DeleteSecretReply, error) {
	err := s.srv.DeleteSecret(core.MustContext(c), req.Id)
	if err != nil {
		return nil, err
	}
	return &secret.DeleteSecretReply{}, nil
}
