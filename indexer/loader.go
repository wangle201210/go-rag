package indexer

import (
	"context"

	"github.com/cloudwego/eino-ext/components/document/loader/file"
	"github.com/cloudwego/eino-ext/components/document/loader/url"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/schema"
	"github.com/wangle201210/go-rag/common"
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
