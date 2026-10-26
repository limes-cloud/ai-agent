package llm

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/limes-cloud/ai-agent/internal/types"

	"github.com/limes-cloud/ai-agent/internal/core"

	"github.com/limes-cloud/ai-agent/internal/domain/entity"
)

func TestNewLLM(t *testing.T) {
	if os.Getenv("AI_AGENT_LLM_INTEGRATION") != "1" {
		t.Skip("set AI_AGENT_LLM_INTEGRATION=1 to run external LLM integration tests")
	}
	tokenValue := os.Getenv("BAIDU_QIANFAN_API_KEY")
	if tokenValue == "" {
		t.Skip("set BAIDU_QIANFAN_API_KEY to run external LLM integration tests")
	}

	model := &entity.Model{
		Keyword:       "qwen-image",
		UUID:          "qwen-image",
		SecretGroupId: 1,
		SecretGroup: &entity.SecretGroup{
			Host:     "http://qianfan.baidubce.com/v2",
			Supplier: entity.SupplierQianfan,
		},
		Extra: `{
			"supportJson": false
		}`,
	}
	token := &entity.Secret{
		Token: tokenValue,
		Group: model.SecretGroup,
	}

	ctx := core.MustContext(context.Background())

	t.Run("test emb", func(t *testing.T) {
		model.Type = entity.ModelTypeEmbedding
		llm := NewLLM(model, token)
		resp, err := llm.Embedding(ctx, &types.EmbeddingRequest{
			Texts: []string{"hello world"},
		})
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 1, len(resp.Vectors))
	})

	t.Run("test chat", func(t *testing.T) {
		model.Type = entity.ModelTypeChat
		llm := NewLLM(model, token)
		resp, err := llm.Text(ctx, &types.TextRequest{
			System: "你是翻译大师，帮我把英文翻译为中文",
			Prompt: "hello world",
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Log(resp.Content)
		assert.NotEmpty(t, resp.Content)
	})

	t.Run("test chat stream", func(t *testing.T) {
		model.Type = entity.ModelTypeChat
		llm := NewLLM(model, token)
		resp, err := llm.Text(ctx, &types.TextRequest{
			System: "你是作诗大师，根据输入的文本，生成一首诗",
			Prompt: "春天",
			Stream: true,
			RecvFunc: func(chunk string) {
				fmt.Println(chunk)
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Log(resp.Content)
		assert.NotEmpty(t, resp.Content)
	})

	t.Run("test chat json", func(t *testing.T) {
		reply := struct {
			Content string `json:"content" jsonschema_description:"结果"`
		}{}

		model.Type = entity.ModelTypeChat
		llm := NewLLM(model, token)
		err := llm.TextScan(ctx, &types.TextRequest{
			System: "你是翻译大师，帮我把英文翻译为中文",
			Prompt: "hello world",
		}, &reply)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(reply.Content)
		assert.NotEmpty(t, reply.Content)
	})

	t.Run("test chat attachment", func(t *testing.T) {
		reply := struct {
			Content string `json:"content" jsonschema_description:"分析结果"`
		}{}

		model.Type = entity.ModelTypeChat
		llm := NewLLM(model, token)
		err := llm.TextScan(ctx, &types.TextRequest{
			System: "你是一个图片颜色分析大师，能够准确分析图片颜色出现次数",
			Attachment: &types.TextAttachment{
				Images: []string{
					"http://118.25.113.106:6081/assets/logo.46a0960c.png",
				},
			},
			Prompt: "图片上出现了几种颜色？",
		}, &reply)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(reply.Content)
		assert.NotEmpty(t, reply.Content)
	})

	t.Run("test image", func(t *testing.T) {
		model.Type = entity.ModelTypeImage
		llm := NewLLM(model, token)
		reply, err := llm.Image(ctx, &types.ImageRequest{
			Prompt: "一只可爱的猫在玩球？",
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Log(reply.List)
		assert.NotEmpty(t, reply.List)
	})
}
