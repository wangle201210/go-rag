package websearch

import (
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/wangle201210/go-rag/server/core/websearch"
)

var (
	webSearchSvr     *websearch.YouComRetriever
	webSearchSvrOnce sync.Once
)

// GetWebSearchSvr 惰性初始化 You.com 联网检索 retriever
// 与 internal/logic/rag 的 init 模式不同：websearch 是可选功能，未配置时不应导致进程启动失败，
// 因此这里延迟到第一次调用时才读取配置并构建 retriever。
func GetWebSearchSvr() *websearch.YouComRetriever {
	webSearchSvrOnce.Do(func() {
		ctx := gctx.New()

		cfg := &websearch.Config{
			APIKey:     g.Cfg().MustGet(ctx, "websearch.youcom.apiKey").String(),
			Count:      g.Cfg().MustGet(ctx, "websearch.youcom.count").Int(),
			TimeoutSec: g.Cfg().MustGet(ctx, "websearch.youcom.timeout").Int(),
		}

		svr, err := websearch.NewRetriever(cfg)
		if err != nil {
			// websearch 是可选功能，构建失败只记录日志，不影响其他功能正常使用
			g.Log().Warningf(ctx, "NewRetriever of websearch failed, err=%v", err)
			return
		}
		webSearchSvr = svr
	})
	return webSearchSvr
}
