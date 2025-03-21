module github.com/luvx21/coding-go/coding-usage

go 1.24.1

//replace github.com/luvx21/coding-go/coding-common => ../coding-common

require (
	github.com/Goldziher/go-utils v1.8.1
	github.com/IBM/sarama v1.45.1
	github.com/PuerkitoBio/goquery v1.10.2
	github.com/allegro/bigcache/v3 v3.1.0
	github.com/apache/pulsar-client-go v0.14.0
	github.com/avast/retry-go/v4 v4.6.1
	github.com/bits-and-blooms/bloom/v3 v3.7.0
	github.com/bytedance/sonic v1.13.1
	github.com/cch123/elasticsql v1.0.1
	github.com/cloudwego/fastpb v0.0.5
	github.com/cloudwego/hertz v0.9.6
	github.com/cloudwego/kitex v0.12.3
	github.com/deckarep/golang-set/v2 v2.8.0
	github.com/dgraph-io/ristretto v0.2.0
	github.com/docker/docker v28.0.1+incompatible
	github.com/dolthub/swiss v0.2.1
	github.com/eko/gocache/lib/v4 v4.2.0
	github.com/eko/gocache/store/bigcache/v4 v4.2.2
	github.com/eko/gocache/store/redis/v4 v4.2.2
	github.com/elastic/go-elasticsearch/v8 v8.17.1
	github.com/emirpasic/gods/v2 v2.0.0-alpha
	github.com/go-co-op/gocron/v2 v2.16.1
	github.com/go-resty/resty/v2 v2.16.5
	github.com/go-sql-driver/mysql v1.9.0
	github.com/gocolly/colly v1.2.0
	github.com/gofiber/fiber/v2 v2.52.6
	github.com/golang/glog v1.2.4
	github.com/google/uuid v1.6.0
	github.com/google/wire v0.6.0
	github.com/gookit/goutil v0.6.18
	github.com/gorilla/websocket v1.5.3
	github.com/jedib0t/go-pretty/v6 v6.6.7
	github.com/jmespath-community/go-jmespath v1.1.1
	github.com/jmoiron/sqlx v1.4.0
	github.com/json-iterator/go v1.1.12
	github.com/labstack/echo/v4 v4.13.3
	github.com/levigross/grequests v0.0.0-20231203190023-9c307ef1f48d
	github.com/lib/pq v1.10.9
	github.com/linvon/cuckoo-filter v0.4.0
	github.com/linxGnu/grocksdb v1.9.8
	github.com/lmittmann/tint v1.0.7
	github.com/logrusorgru/aurora v2.0.3+incompatible
	github.com/loov/hrtime v1.0.3
	github.com/luvx21/coding-go/coding-common v0.0.0-20250307100758-4a92a8615909
	github.com/luvx21/coding-go/infra/nosql/mongodb v0.0.0-20250307100758-4a92a8615909
	github.com/marcboeker/go-duckdb v1.8.5
	github.com/mattn/go-sqlite3 v1.14.24
	github.com/panjf2000/ants/v2 v2.11.2
	github.com/parnurzeal/gorequest v0.3.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/prometheus/client_golang v1.21.1
	github.com/pterm/pterm v0.12.80
	github.com/puzpuzpuz/xsync/v3 v3.5.1
	github.com/pywee/gobson-where v0.0.4
	github.com/redis/go-redis/v9 v9.7.1
	github.com/redis/rueidis v1.0.56
	github.com/robfig/cron/v3 v3.0.1
	github.com/samber/lo v1.49.1
	github.com/segmentio/kafka-go v0.4.47
	github.com/sethvargo/go-retry v0.3.0
	github.com/smallnest/exp v0.7.1
	github.com/sourcegraph/conc v0.3.0
	github.com/thedevsaddam/gojsonq/v2 v2.5.2
	github.com/tidwall/gjson v1.18.0
	github.com/tursodatabase/go-libsql v0.0.0-20250313100617-0ab5a1a61a71
	github.com/tursodatabase/libsql-client-go v0.0.0-20240902231107-85af5b9d094d
	github.com/valyala/fasthttp v1.59.0
	github.com/withlin/canal-go v1.1.2
	github.com/xxl-job/xxl-job-executor-go v1.2.0
	github.com/yanyiwu/gojieba v1.4.5
	go-micro.dev/v4 v4.11.0
	go.mongodb.org/mongo-driver v1.17.3
	go.uber.org/zap v1.27.0
	golang.org/x/sync v0.12.0
	golang.org/x/time v0.11.0
	google.golang.org/grpc v1.71.0
	google.golang.org/protobuf v1.36.5
)

require (
	github.com/bytedance/gopkg v0.1.1 // indirect
	github.com/cloudwego/netpoll v0.6.5 // indirect
	github.com/nyaruka/phonenumbers v1.0.55 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.33.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
)

require (
	atomicgo.dev/cursor v0.2.0 // indirect
	atomicgo.dev/keyboard v0.2.9 // indirect
	atomicgo.dev/schedule v0.1.0 // indirect
	dario.cat/mergo v1.0.1 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4 // indirect
	github.com/99designs/keyring v1.2.1 // indirect
	github.com/AthenZ/athenz v1.10.39 // indirect
	github.com/DataDog/zstd v1.5.0 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/ProtonMail/go-crypto v1.1.6 // indirect
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/andybalholm/cascadia v1.3.3 // indirect
	github.com/antchfx/htmlquery v1.3.4 // indirect
	github.com/antchfx/xmlquery v1.4.4 // indirect
	github.com/antchfx/xpath v1.3.3 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/apache/arrow-go/v18 v18.2.0 // indirect
	github.com/ardielle/ardielle-go v1.5.2 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bitly/go-simplejson v0.5.1 // indirect
	github.com/bits-and-blooms/bitset v1.22.0 // indirect
	github.com/blockloop/scan v1.3.0 // indirect
	github.com/bufbuild/protocompile v0.14.1 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/bytedance/sonic/loader v0.2.4 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudflare/circl v1.6.0 // indirect
	github.com/cloudwego/base64x v0.1.5 // indirect
	github.com/cloudwego/configmanager v0.2.2 // indirect
	github.com/cloudwego/dynamicgo v0.5.2 // indirect
	github.com/cloudwego/frugal v0.2.3 // indirect
	github.com/cloudwego/gopkg v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/cloudwego/localsession v0.1.2 // indirect
	github.com/cloudwego/runtimex v0.1.1 // indirect
	github.com/cloudwego/thriftgo v0.3.18 // indirect
	github.com/coder/websocket v1.8.13 // indirect
	github.com/containerd/console v1.0.4 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.6 // indirect
	github.com/cyphar/filepath-securejoin v0.4.1 // indirect
	github.com/danieljoos/wincred v1.1.2 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dgryski/go-metro v0.0.0-20250106013310-edb8663e5e33 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/docker/go-connections v0.5.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/dolthub/maphash v0.1.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/dvsekhvalnov/jose2go v1.6.0 // indirect
	github.com/eapache/go-resiliency v1.7.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20230731223053-c322873962e3 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/elastic/elastic-transport-go/v8 v8.6.1 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/evanphx/json-patch/v5 v5.5.0 // indirect
	github.com/fatih/structtag v1.2.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/go-acme/lego/v4 v4.4.0 // indirect
	github.com/go-basic/ipv4 v1.0.0 // indirect
	github.com/go-git/gcfg v1.5.1-0.20230307220236-3a3c6141e376 // indirect
	github.com/go-git/go-billy/v5 v5.6.2 // indirect
	github.com/go-git/go-git/v5 v5.14.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/go-zookeeper/zk v1.0.3 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.2.1 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/godbus/dbus v0.0.0-20190726142602-4481cbc300e2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/golang/groupcache v0.0.0-20241129210726-2c02b8208cf8 // indirect
	github.com/golang/mock v1.7.0-rc.1 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/google/flatbuffers v25.2.10+incompatible // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/pprof v0.0.0-20250208200701-d0013a598941 // indirect
	github.com/gookit/color v1.5.4 // indirect
	github.com/gorilla/handlers v1.5.1 // indirect
	github.com/gsterjov/go-libsecret v0.0.0-20161001094733-a6f4afe4910c // indirect
	github.com/hamba/avro/v2 v2.28.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/iancoleman/strcase v0.3.0 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.4 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/jhump/protoreflect v1.17.0 // indirect
	github.com/jonboulle/clockwork v0.5.0 // indirect
	github.com/kennygrant/sanitize v1.2.4 // indirect
	github.com/kevinburke/ssh_config v1.2.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.10 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible // indirect
	github.com/lestrrat-go/strftime v1.1.0 // indirect
	github.com/libsql/sqlite-antlr4-parser v0.0.0-20240721121621-c0bdc870f11c // indirect
	github.com/lithammer/fuzzysearch v1.1.8 // indirect
	github.com/luvx12/logrus-prefixed-formatter v0.5.6 // indirect
	github.com/mailru/easyjson v0.9.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/miekg/dns v1.1.63 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/moby/sys/userns v0.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/moul/http2curl v1.0.0 // indirect
	github.com/mtibben/percent v0.2.1 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/nxadm/tail v1.4.11 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0 // indirect
	github.com/oxtoacart/bpool v0.0.0-20190530202638-03653db5a59c // indirect
	github.com/pierrec/lz4 v2.0.5+incompatible // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/pjbgf/sha1cd v0.3.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.63.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/saintfish/chardet v0.0.0-20230101081208-5e3ef4b5456d // indirect
	github.com/sergi/go-diff v1.3.2-0.20230802210424-5b0b94c5c0d3 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/skeema/knownhosts v1.3.1 // indirect
	github.com/smartystreets/goconvey v1.8.1 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	github.com/temoto/robotstxt v1.1.2 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/urfave/cli/v2 v2.27.6 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/xanzy/ssh-agent v0.3.3 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	github.com/xrash/smetrics v0.0.0-20240521201337-686a1a2994c1 // indirect
	github.com/xwb1989/sqlparser v0.0.0-20180606152119-120387863bf2 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	github.com/zeebo/xxh3 v1.0.2 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.57.0 // indirect
	go.opentelemetry.io/otel v1.35.0 // indirect
	go.opentelemetry.io/otel/metric v1.35.0 // indirect
	go.opentelemetry.io/otel/trace v1.35.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/mock v0.5.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/arch v0.15.0 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/exp v0.0.0-20250305212735-054e65f0b394 // indirect
	golang.org/x/mod v0.24.0 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/oauth2 v0.25.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/term v0.30.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	golang.org/x/tools v0.31.0 // indirect
	golang.org/x/xerrors v0.0.0-20240903120638-7835f813f4da // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250106144421-5f5ef82da422 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250313205543-e70fdf4c4cb4 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gotest.tools/v3 v3.5.1 // indirect
)
