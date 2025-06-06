package cmd

import (
	"context"
	"net/http"

	"github.com/ThinkInAIXYZ/go-mcp/server"
	"github.com/ThinkInAIXYZ/go-mcp/transport"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/wangle201210/go-rag/server/internal/controller/rag"
	"github.com/wangle201210/go-rag/server/internal/mcp"
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

func init() {
	Main.AddCommand(&gcmd.Command{
		Name:  "mcp",
		Usage: "mcp",
		Brief: "start mcp server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			trans, handler, err := transport.NewStreamableHTTPServerTransportAndHandler()
			if err != nil {
				g.Log().Panicf(ctx, "new sse transport and hander with error: %v", err)
			}
			// new mcp server
			mcpServer, _ := server.NewServer(trans)
			// register tool with mcpServer
			mcpServer.RegisterTool(mcp.GetIndexerByFileBase64ContentTool(), mcp.HandleIndexerByFileBase64Content)
			// mcpServer.RegisterTool(mcp.GetIndexerByFilePathTool(), mcp.HandleIndexerByFilePath)
			mcpServer.RegisterTool(mcp.GetRetrieverTool(), mcp.HandleRetriever)
			mcpServer.RegisterTool(mcp.GetKnowledgeBaseTool(), mcp.HandleKnowledgeBase)
			// start mcp Server
			go func() {
				mcpServer.Run()
			}()
			defer mcpServer.Shutdown(context.Background())
			http.Handle("/mcp", handler.HandleMCP())
			return http.ListenAndServe(":8089", nil)
		},
	})
}
