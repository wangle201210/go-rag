package server

import (
	"context"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/wangle201210/go-rag/server/internal/cmd"
)

// Start 启动 go-rag 服务器（公开函数）
func Start(ctx context.Context) {
	if ctx == nil {
		ctx = gctx.GetInitCtx()
	}
	cmd.Main.Run(ctx)
}
