package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/wangle201210/go-rag/server/internal/controller/rag"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			s.Group("/", func(group *ghttp.RouterGroup) {
				s.AddStaticPath("", "./static/fe/")
				s.SetIndexFiles([]string{"index.html"})
				group.Middleware(ghttp.MiddlewareHandlerResponse, ghttp.MiddlewareCORS)
				group.Bind(
					rag.NewV1(),
				)
			})
			s.Run()
			return nil
		},
	}
)
