package indexer

import (
	"context"
	"fmt"
	"log"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino-ext/components/indexer/es8"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/indexer"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	"github.com/wangle201210/go-rag/server/core/common"
	"github.com/wangle201210/go-rag/server/core/config"
)

var system = "知识库名字是：《deepchat使用文档》\n " +
	"请把用户输入内容拆解为QA对，数量控制在3-5个"

type QA struct {
	Item []*Item `json:"item" jsonschema:"description=拆解出来的QA问答对"`
}

type Item struct {
	Question string `json:"question" jsonschema:"description=问题"`
	Answer   string `json:"answer" jsonschema:"description=答案"`
}

// 处理函数
func insertQAFunc(ctx context.Context, params *QA) (string, error) {
	var chunkID string
	if value, ok := ctx.Value(common.DocQAChunks).(string); ok {
		chunkID = value
	} else {
		return "", fmt.Errorf("必须提供chunkID")
	}
	docs := make([]*schema.Document, 0, len(params.Item))
	for _, item := range params.Item {
		doc := &schema.Document{
			ID:      uuid.New().String(),
			Content: fmt.Sprintf("%s", item.Question),
			MetaData: map[string]any{
				common.DocQAChunks: chunkID,
				common.DocQAAnswer: item.Answer,
			},
		}
		docs = append(docs, doc)
	}
	_, err := getQAIndexer().Store(ctx, docs)
	if err != nil {
		return "", err
	}
	return "success", nil
}

func docQA(ctx context.Context, conf *config.Config, doc *schema.Document) (err error) {
	// 使用 InferTool 创建工具
	qaTool, err := utils.InferTool(
		"disassembleQA", // tool name
		"请把输入的内容拆解为QA对，数量控制在3-5个", // tool description
		insertQAFunc)
	tools := []tool.BaseTool{
		qaTool,
	}
	cm, err := common.GetChatModel(ctx, conf.GetChatModelConfig())
	if err != nil {
		return err
	}
	// 获取工具信息并绑定到 ChatModel
	toolInfos := make([]*schema.ToolInfo, 0, len(tools))
	for _, t := range tools {
		info, err := t.Info(ctx)
		if err != nil {
			log.Fatal(err)
		}
		toolInfos = append(toolInfos, info)
	}
	tcm, err := cm.(*openai.ChatModel).WithTools(toolInfos)
	if err != nil {
		log.Fatal(err)
	}
	// 创建 tools 节点
	toolsNode, err := compose.NewToolNode(context.Background(), &compose.ToolsNodeConfig{
		Tools: tools,
	})
	if err != nil {
		log.Fatal(err)
	}

	// 构建完整的处理链
	chain := compose.NewChain[[]*schema.Message, []*schema.Message]()
	chain.
		AppendChatModel(tcm, compose.WithNodeName("chat_model")).
		AppendToolsNode(toolsNode, compose.WithNodeName("tools"))

	// 编译并运行 chain
	agent, err := chain.Compile(ctx)
	if err != nil {
		log.Fatal(err)
	}
	ctx = context.WithValue(ctx, common.DocQAChunks, doc.ID)
	// 运行示例
	resp, err := agent.Invoke(ctx, []*schema.Message{
		{
			Role:    schema.System,
			Content: system,
		},
		{
			Role:    schema.User,
			Content: doc.Content,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// 输出结果
	for _, msg := range resp {
		fmt.Println(msg.Content)
	}
	return nil
}

var qaIndexer indexer.Indexer

func InitQaIndexer(ctx context.Context, conf *config.Config) error {
	idr, err := newQAIndexer(ctx, conf)
	if err != nil {
		return err
	}
	qaIndexer = idr
	return nil
}

func getQAIndexer() (idr indexer.Indexer) {
	return qaIndexer
}

func newQAIndexer(ctx context.Context, conf *config.Config) (idr indexer.Indexer, err error) {
	indexerConfig := &es8.IndexerConfig{
		Client:    conf.Client,
		Index:     conf.IndexName,
		BatchSize: 10,
		DocumentToFields: func(ctx context.Context, doc *schema.Document) (field2Value map[string]es8.FieldValue, err error) {
			var knowledgeName string
			if value, ok := ctx.Value(common.KnowledgeName).(string); ok {
				knowledgeName = value
			} else {
				err = fmt.Errorf("必须提供知识库名称")
				return
			}
			doc.ID = uuid.New().String()
			if doc.MetaData != nil {
				marshal, _ := sonic.Marshal(doc.MetaData)
				doc.MetaData[common.DocExtra] = string(marshal)
			}
			return map[string]es8.FieldValue{
				common.FieldContent: {
					Value:    doc.Content,
					EmbedKey: common.FieldContentVector, // vectorize doc content and save vector to field "content_vector"
				},
				common.DocQAChunks: {
					Value: doc.MetaData[common.DocQAChunks],
				},
				common.KnowledgeName: {
					Value: knowledgeName,
				},
				common.FieldExtra: {
					Value: doc.MetaData[common.DocExtra],
				},
			}, nil
		},
	}
	embeddingIns11, err := common.NewEmbedding(ctx, conf)
	if err != nil {
		return nil, err
	}
	indexerConfig.Embedding = embeddingIns11
	idr, err = es8.NewIndexer(ctx, indexerConfig)
	if err != nil {
		return nil, err
	}
	return idr, nil
}

func getQAContent(ctx context.Context, conf *config.Config, doc *schema.Document, knowledgeName string) (qaContent string, err error) {
	cm, err := common.GetNotThinkChatModel(ctx, nil)
	if err != nil {
		return
	}
	generate, err := cm.Generate(ctx, []*schema.Message{
		{
			Role: schema.System,
			Content: fmt.Sprintf("你是一个专业的问题生成助手，任务是从给定的文本中提取或生成可能的问题。你不需要回答这些问题，只需生成问题本身。\n"+
				"知识库名字是：《%s》\n\n"+
				"输出格式：\n"+
				"- 每个问题占一行\n"+
				"- 问题必须以问号结尾\n"+
				"- 避免重复或语义相似的问题\n\n"+
				"生成规则：\n"+
				"- 生成的问题必须严格基于文本内容，不能脱离文本虚构。\n"+
				"- 优先生成事实性问题（如谁、何时、何地、如何）。\n"+
				"- 对于复杂文本，可生成多层次问题（基础事实 + 推理问题）。\n"+
				"- 禁止生成主观或开放式问题（如“你认为...？”）。"+
				"- 数量控制在3-5个", knowledgeName),
		},
		{
			Role:    schema.User,
			Content: doc.Content,
		},
	})
	if err != nil {
		return
	}
	qaContent = generate.Content
	return
}
