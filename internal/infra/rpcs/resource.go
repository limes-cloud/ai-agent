package rpcs

import (
	"github.com/limes-cloud/ai-agent/internal/core"
	"github.com/limes-cloud/resource/api/file"
)

type Resource struct{}

func NewResource() *Resource {
	return &Resource{}
}

func (r Resource) client(ctx core.Context) (file.FileServiceClient, error) {
	conn, err := ctx.GrpcConn("Resource")
	if err != nil {
		return nil, err
	}
	return file.NewFileServiceClient(conn), nil
}

func (r Resource) UploadFileByURL(ctx core.Context, req *file.UploadFileByURLRequest) (*file.UploadFileByURLResponse, error) {
	client, err := r.client(ctx)
	if err != nil {
		return nil, err
	}
	return client.UploadFileByURL(ctx, req)
}
