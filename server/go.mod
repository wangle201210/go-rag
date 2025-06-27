module github.com/wangle201210/go-rag/server

go 1.23.1

require (
	github.com/ThinkInAIXYZ/go-mcp v0.2.14
	github.com/bytedance/sonic v1.13.2
	github.com/cenkalti/backoff/v4 v4.3.0
	github.com/cloudwego/eino v0.3.31
	github.com/cloudwego/eino-ext/components/document/loader/file v0.0.0-20250424061409-ccd60fbc7c1c
	github.com/cloudwego/eino-ext/components/document/loader/url v0.0.0-20250610035057-2c4e7c8488a5
	github.com/cloudwego/eino-ext/components/document/parser/html v0.0.0-20250424061409-ccd60fbc7c1c
	github.com/cloudwego/eino-ext/components/document/parser/pdf v0.0.0-20250424061409-ccd60fbc7c1c
	github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown v0.0.0-20250610035057-2c4e7c8488a5
	github.com/cloudwego/eino-ext/components/document/transformer/splitter/recursive v0.0.0-20250424061409-ccd60fbc7c1c
	github.com/cloudwego/eino-ext/components/embedding/openai v0.0.0-20250424061409-ccd60fbc7c1c
	github.com/cloudwego/eino-ext/components/indexer/es8 v0.0.0-20250610035057-2c4e7c8488a5
	github.com/cloudwego/eino-ext/components/model/openai v0.0.0-20250610035057-2c4e7c8488a5
	github.com/cloudwego/eino-ext/components/model/qwen v0.0.0-20250610035057-2c4e7c8488a5
	github.com/cloudwego/eino-ext/components/retriever/es8 v0.0.0-20250610035057-2c4e7c8488a5
	github.com/elastic/go-elasticsearch/v8 v8.16.0
	github.com/gogf/gf/contrib/drivers/mysql/v2 v2.9.0
	github.com/gogf/gf/v2 v2.9.0
	github.com/google/uuid v1.6.0
	github.com/stretchr/testify v1.10.0
	github.com/wangle201210/chat-history v0.0.0-20250402104704-5eec15d5419e
	gorm.io/driver/mysql v1.5.7
	gorm.io/gorm v1.25.12
)

require (
	github.com/BurntSushi/toml v1.4.0 // indirect
	github.com/PuerkitoBio/goquery v1.8.1 // indirect
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/bytedance/sonic/loader v0.2.4 // indirect
	github.com/clbanning/mxj/v2 v2.7.0 // indirect
	github.com/cloudwego/base64x v0.1.5 // indirect
	github.com/cloudwego/eino-ext/libs/acl/openai v0.0.0-20250610035057-2c4e7c8488a5 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dslipak/pdf v0.0.2 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/elastic/elastic-transport-go/v8 v8.7.0 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/evanphx/json-patch v0.5.2 // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/getkin/kin-openapi v0.118.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/jsonpointer v0.21.1 // indirect
	github.com/go-openapi/swag v0.23.1 // indirect
	github.com/go-sql-driver/mysql v1.7.1 // indirect
	github.com/goph/emperror v0.17.2 // indirect
	github.com/gorilla/css v1.0.1 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/grokify/html-strip-tags-go v0.1.0 // indirect
	github.com/invopop/yaml v0.1.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.10 // indirect
	github.com/magiconair/properties v1.8.9 // indirect
	github.com/mailru/easyjson v0.9.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/meguminnnnnnnnn/go-openai v0.0.0-20250530094841-88286040d3c1 // indirect
	github.com/microcosm-cc/bluemonday v1.0.27 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/nikolalohinski/gonja v1.5.3 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/orcaman/concurrent-map/v2 v2.0.1 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/perimeterx/marshmallow v1.1.5 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/slongfield/pyfmt v0.0.0-20220222012616-ea85ff4c361f // indirect
	github.com/tidwall/gjson v1.18.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/yargevad/filepathx v1.0.0 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
	go.opentelemetry.io/otel v1.32.0 // indirect
	go.opentelemetry.io/otel/metric v1.32.0 // indirect
	go.opentelemetry.io/otel/sdk v1.32.0 // indirect
	go.opentelemetry.io/otel/trace v1.32.0 // indirect
	golang.org/x/arch v0.15.0 // indirect
	golang.org/x/exp v0.0.0-20250305212735-054e65f0b394 // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
