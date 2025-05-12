package main

import (
	"github.com/gogf/gf/v2/os/gctx"

	"github.com/wangle201210/go-rag/server/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
