# chat-history
聊天历史记录管理

## 为什么需要
目前eino提供了快捷接入大模型的方法但是聊天历史记录的功能没有提供，导致没法进行对轮对话，及对话管理

## 怎么接入
可以参照 [example](./example)
1. 安装依赖 `go get github.com/wangle201210/chat-history` 

2. 初始化 EinoHistory
```go
var eh = eino.NewEinoHistory("root:123456@tcp(127.0.0.1:3306)/chat_history")
```

3. 在原来获取messages的地方添加 add start -- add end 之间的代码 
```go
func createMessagesFromTemplate(ctx context.Context, convID, question string) (messages []*schema.Message, err error) {
	template := createTemplate(ctx)
	/* add start */
	chatHistory, err := eh.GetHistory(convID, 100)
	if err != nil {
		return
	}
	// 插入一条用户数据
	err = eh.SaveMessage(&schema.Message{
		Role:    schema.User,
		Content: question,
	}, convID)
	if err != nil {
		return
	}
	/* add end */
	// 使用模板生成消息
	messages, err = template.Format(context.Background(), map[string]any{
		"role":         "程序员鼓励师",
		"style":        "积极、温暖且专业",
		"question":     question,
		"chat_history": chatHistory,
	})
	if err != nil {
		return
	}
	return
}
```
4. 在返回消息的地方添加 add start -- add end 之间的代码 
```go
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
```