package common

import (
	"testing"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func TestCfg(t *testing.T) {
	ctx := gctx.New()
	cfg := &openai.ChatModelConfig{}
	err := g.Cfg().MustGet(ctx, "qa").Scan(cfg)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(cfg)
}
