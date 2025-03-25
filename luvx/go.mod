module luvx

go 1.24.1

replace (
	//github.com/luvx21/coding-go/coding-common => ../coding-common
	//github.com/luvx21/coding-go/infra/logs => ../infra/logs
	//github.com/luvx21/coding-go/infra/nosql/mongodb => ../infra/nosql/mongodb
)

require (
	github.com/PuerkitoBio/goquery v1.10.2
	github.com/allegro/bigcache/v3 v3.1.0
	github.com/bytedance/sonic v1.13.2
	github.com/eko/gocache/lib/v4 v4.2.0
	github.com/eko/gocache/store/bigcache/v4 v4.2.2
	github.com/gin-gonic/gin v1.10.0
	github.com/go-co-op/gocron/v2 v2.16.1
	github.com/go-sql-driver/mysql v1.9.1
	github.com/gocolly/colly v1.2.0
	github.com/google/uuid v1.6.0
	github.com/google/wire v0.6.0
	github.com/luvx21/coding-go/coding-common v0.0.0-20250323084229-6da5d10efe82
	github.com/luvx21/coding-go/infra/ai v0.0.0-00010101000000-000000000000
	github.com/luvx21/coding-go/infra/infra_sql v0.0.0-20250323084229-6da5d10efe82
	github.com/luvx21/coding-go/infra/logs v0.0.0-20250323084229-6da5d10efe82
	github.com/luvx21/coding-go/infra/nosql/mongodb v0.0.0-20250323084229-6da5d10efe82
	github.com/parnurzeal/gorequest v0.3.0
	github.com/redis/go-redis/v9 v9.7.3
	github.com/tidwall/gjson v1.18.0
	github.com/tursodatabase/libsql-client-go v0.0.0-20240902231107-85af5b9d094d
	github.com/valyala/fasthttp v1.59.0
	go.mongodb.org/mongo-driver v1.17.3
	golang.org/x/crypto v0.36.0
	golang.org/x/exp v0.0.0-20250305212735-054e65f0b394
	golang.org/x/sync v0.12.0
	golang.org/x/time v0.11.0
	gorm.io/driver/mysql v1.5.7
	gorm.io/gen v0.3.26
	gorm.io/gorm v1.25.12
	gorm.io/plugin/dbresolver v1.5.3
	modernc.org/sqlite v1.36.3
)

require (
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/smarty/assertions v1.16.0 // indirect
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/andybalholm/cascadia v1.3.3 // indirect
	github.com/antchfx/htmlquery v1.3.4 // indirect
	github.com/antchfx/xmlquery v1.4.4 // indirect
	github.com/antchfx/xpath v1.3.3 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blockloop/scan v1.3.0 // indirect
	github.com/bytedance/sonic/loader v0.2.4 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.5 // indirect
	github.com/coder/websocket v1.8.13 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/gin-contrib/sse v1.0.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.26.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/golang/groupcache v0.0.0-20241129210726-2c02b8208cf8 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmoiron/sqlx v1.4.0 // indirect
	github.com/jonboulle/clockwork v0.5.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kennygrant/sanitize v1.2.4 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.10 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible // indirect
	github.com/lestrrat-go/strftime v1.1.0 // indirect
	github.com/luvx12/logrus-prefixed-formatter v0.5.6 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/moul/http2curl v1.0.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/ncruces/go-strftime v0.1.9 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.36.3 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.21.1 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.63.0 // indirect
	github.com/prometheus/procfs v0.16.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/sagikazarmark/locafero v0.9.0 // indirect
	github.com/saintfish/chardet v0.0.0-20230101081208-5e3ef4b5456d // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.14.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/spf13/viper v1.20.1 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/temoto/robotstxt v1.1.2 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	go.uber.org/mock v0.5.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/arch v0.15.0 // indirect
	golang.org/x/mod v0.24.0 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/term v0.30.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	golang.org/x/tools v0.31.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/datatypes v1.2.5 // indirect
	gorm.io/hints v1.1.2 // indirect
	modernc.org/libc v1.62.0 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.9.1 // indirect
)
