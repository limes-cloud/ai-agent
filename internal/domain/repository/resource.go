package repository

import (
	"github.com/limes-cloud/ai-agent/internal/core"
	"github.com/limes-cloud/resource/api/file"
)

type Resource interface {
	UploadFileByURL(ctx core.Context, req *file.UploadFileByURLRequest) (*file.UploadFileByURLResponse, error)
}
