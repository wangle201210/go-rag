package main

import (
	"github.com/gogf/gf/v2/os/gctx"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
	"github.com/wangle201210/go-rag/server/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
