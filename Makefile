.PHONY: proto build run test docker

# 生成 proto 代码
proto:
	kratosx proto client api/generate/proto/generate_service.proto
	kratosx proto client api/model/proto/model_service.proto
	kratosx proto client api/prompt/proto/prompt_service.proto
	kratosx proto client api/secret/proto/secret_service.proto
	kratosx proto client api/errors/proto/errors.proto

# 编译
build:
	go build -ldflags="-s -w" -o ai-agent main.go

# 本地运行
run:
	go run main.go -conf conf/

# 测试
test:
	go test ./... -v -count=1

# Docker 构建
docker:
	docker build -t ai-agent:latest .