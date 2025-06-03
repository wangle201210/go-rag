package indexer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/eino-ext/components/document/loader/file"
	"github.com/cloudwego/eino-ext/components/document/loader/url"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	"github.com/wangle201210/go-rag/server/core/common"
	"github.com/wangle201210/go-rag/server/core/model"
	"github.com/wangle201210/go-rag/server/core/server"
)

// newLoader component initialization function of node 'Loader1' in graph 'rag'
func newLoader(ctx context.Context) (ldr document.Loader, err error) {
	mldr := &multiLoader{}
	parser, err := newParser(ctx)
	if err != nil {
		return nil, err
	}
	fldr, err := file.NewFileLoader(ctx, &file.FileLoaderConfig{
		UseNameAsID: true,
		Parser:      parser,
	})
	if err != nil {
		return nil, err
	}
	mldr.fileLoader = fldr
	uldr, err := url.NewLoader(ctx, &url.LoaderConfig{})
	if err != nil {
		return nil, err
	}
	mldr.urlLoader = uldr
	return mldr, nil
}

type multiLoader struct {
	fileLoader document.Loader
	urlLoader  document.Loader
}

func (x *multiLoader) Load(ctx context.Context, src document.Source, opts ...document.LoaderOption) ([]*schema.Document, error) {
	if common.IsURL(src.URI) {
		return x.urlLoader.Load(ctx, src, opts...)
	}
	return x.fileLoader.Load(ctx, src, opts...)
}

func GetLoaderCallback() callbacks.Handler {
	return callbacks.NewHandlerBuilder().OnEndFn(
		func(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
			o := output.(document.LoaderCallbackOutput)
			for _, doc := range o.Docs {
				server.GetDocumentMapping().Create(ctx, &model.DocumentMappingCreateReq{
					KnowledgeBaseId: uuid.NewString(),
					DocumentId:      uuid.NewString(),
					DocumentName:    doc.MetaData["_file_name"].(string),
					DocumentType:    doc.MetaData["_extension"].(string),
					DocumentPath:    doc.MetaData["_source"].(string),
					DocumentSize:    int64(len(doc.Content)),
				})
			}

			marshal, err := json.Marshal(o)
			if err != nil {
				return ctx
			}
			fmt.Printf("output: %s", marshal)
			return ctx
		}).Build()
}

//
// type Loader struct{}
//
// func NewLoader() *Loader {
// 	return &Loader{}
// }
//
// func (l *Loader) Load(filePath string) (string, error) {
// 	// 检查文件是否存在
// 	if _, err := os.Stat(filePath); os.IsNotExist(err) {
// 		return "", fmt.Errorf("文件不存在: %s", filePath)
// 	}
//
// 	// 读取文件内容
// 	content, err := ioutil.ReadFile(filePath)
// 	if err != nil {
// 		return "", fmt.Errorf("读取文件失败: %v", err)
// 	}
//
// 	// 获取文件扩展名
// 	ext := filepath.Ext(filePath)
// 	if ext == "" {
// 		return "", fmt.Errorf("文件没有扩展名: %s", filePath)
// 	}
//
// 	return string(content), nil
// }
