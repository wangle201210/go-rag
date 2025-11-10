package main

import (
	"context"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/google/uuid"
	"github.com/wangle201210/go-rag/server/chat-history/eino"
)

// 初始化一个
var eh = eino.NewEinoHistory("root:123456@tcp(127.0.0.1:3306)/chat_history?charset=utf8mb4&parseTime=True&loc=Local")

func createOpenAIChatModel(ctx context.Context) model.ChatModel {
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		Model:   os.Getenv("OPENAI_MODEL_NAME"),
		APIKey:  os.Getenv("OPENAI_API_KEY"),
		BaseURL: os.Getenv("OPENAI_BASE_URL"),
	})
	if err != nil {
		log.Fatalf("create openai chat model failed, err=%v", err)
	}
	return chatModel
}

func main() {
	ctx := context.Background()
	var convID = uuid.NewString() // 模拟一个会话id
	cm := createOpenAIChatModel(ctx)
	// 模拟用户连续问问题
	messList := []string{
		"我数学不好",
		"数学不好可以编程么",
		"我语文好有没有什么优势",
		"刚才我说我什么科目学得不好来着？",
	}
	for _, s := range messList {
		messages, err := createMessagesFromTemplate(ctx, convID, s)
		if err != nil {
			log.Fatalf("create messages failed: %v", err)
			return
		}
		result := generate(ctx, cm, messages)
		/* add start */
		err = eh.SaveMessage(result, convID)
		if err != nil {
			log.Fatalf("save assistant message err: %v", err)
			return
		}
		/* add end */
		log.Printf("result: %+v\n\n", result)
	}
}
