module github.com/fastenmind/fastener-api

go 1.23.0

toolchain go1.24.5

require (
	github.com/99designs/gqlgen v0.17.78
	github.com/aws/aws-sdk-go v1.44.327
	github.com/go-playground/validator/v10 v10.27.0
	github.com/go-redis/redis/v8 v8.11.5
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/google/uuid v1.6.0
	github.com/gorilla/websocket v1.5.3
	github.com/jackc/pgx/v5 v5.6.0
	github.com/joho/godotenv v1.5.1
	github.com/jung-kurt/gofpdf v1.16.2
	github.com/labstack/echo/v4 v4.11.4
	github.com/labstack/gommon v0.4.2
	github.com/newrelic/go-agent/v3 v3.40.1
	github.com/prometheus/client_golang v1.23.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/shirou/gopsutil/v3 v3.24.5
	github.com/shopspring/decimal v1.3.1
	github.com/sony/gobreaker v1.0.0
	github.com/streadway/amqp v1.1.0
	github.com/stretchr/testify v1.10.0
	github.com/swaggo/echo-swagger v1.4.1
	github.com/xuri/excelize/v2 v2.9.1
	go.elastic.co/apm v1.15.0
	go.elastic.co/apm/module/apmecho v1.15.0
	go.elastic.co/apm/module/apmgorm v1.15.0
	go.elastic.co/apm/module/apmhttp v1.15.0
	go.opentelemetry.io/otel v1.37.0
	go.opentelemetry.io/otel/exporters/jaeger v1.17.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.37.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.37.0
	go.opentelemetry.io/otel/sdk v1.37.0
	go.opentelemetry.io/otel/trace v1.37.0
	go.uber.org/zap v1.27.0
	golang.org/x/crypto v0.40.0
	golang.org/x/text v0.27.0
	golang.org/x/time v0.11.0
	google.golang.org/grpc v1.73.0
	gopkg.in/DataDog/dd-trace-go.v1 v1.74.3
	gopkg.in/yaml.v3 v3.0.1
	gorm.io/datatypes v1.2.0
	gorm.io/driver/postgres v1.5.5
	gorm.io/gorm v1.26.0
	gorm.io/plugin/dbresolver v1.6.2
)

require (
	github.com/DataDog/appsec-internal-go v1.13.0 // indirect
	github.com/DataDog/datadog-agent/comp/core/tagger/origindetection v0.66.1 // indirect
	github.com/DataDog/datadog-agent/pkg/obfuscate v0.66.1 // indirect
	github.com/DataDog/datadog-agent/pkg/proto v0.66.1 // indirect
	github.com/DataDog/datadog-agent/pkg/remoteconfig/state v0.66.1 // indirect
	github.com/DataDog/datadog-agent/pkg/trace v0.66.1 // indirect
	github.com/DataDog/datadog-agent/pkg/util/log v0.66.1 // indirect
	github.com/DataDog/datadog-agent/pkg/util/scrubber v0.66.1 // indirect
	github.com/DataDog/datadog-agent/pkg/version v0.66.1 // indirect
	github.com/DataDog/datadog-go/v5 v5.6.0 // indirect
	github.com/DataDog/dd-trace-go/v2 v2.1.0 // indirect
	github.com/DataDog/go-libddwaf/v4 v4.3.0 // indirect
	github.com/DataDog/go-runtime-metrics-internal v0.0.4-0.20250603194815-7edb7c2ad56a // indirect
	github.com/DataDog/go-sqllexer v0.1.6 // indirect
	github.com/DataDog/go-tuf v1.1.0-0.5.2 // indirect
	github.com/DataDog/opentelemetry-mapping-go/pkg/otlp/attributes v0.26.0 // indirect
	github.com/DataDog/sketches-go v1.4.7 // indirect
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/Masterminds/semver/v3 v3.3.1 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/agnivade/levenshtein v1.2.1 // indirect
	github.com/armon/go-radix v1.0.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff/v5 v5.0.2 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/eapache/queue/v2 v2.0.0-20230407133247-75960ed334e4 // indirect
	github.com/ebitengine/purego v0.8.3 // indirect
	github.com/elastic/go-licenser v0.3.1 // indirect
	github.com/elastic/go-sysinfo v1.1.1 // indirect
	github.com/elastic/go-windows v1.0.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/spec v0.20.4 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.1 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jcchavezs/porto v0.1.0 // indirect
	github.com/jinzhu/gorm v1.9.16 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/joeshaw/multierror v0.0.0-20140124173710-69b34d4ec901 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/labstack/echo v3.3.10+incompatible // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/lufia/plan9stats v0.0.0-20250317134145-8bc96cf8fc35 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/onsi/gomega v1.36.2 // indirect
	github.com/outcaste-io/ristretto v0.2.3 // indirect
	github.com/philhofer/fwd v1.1.3-0.20240916144458-20a13a1f6b7c // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/power-devops/perfstat v0.0.0-20240221224432-82ca36839d55 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.65.0 // indirect
	github.com/prometheus/procfs v0.16.1 // indirect
	github.com/puzpuzpuz/xsync/v3 v3.5.1 // indirect
	github.com/richardlehane/mscfb v1.0.4 // indirect
	github.com/richardlehane/msoleps v1.0.4 // indirect
	github.com/santhosh-tekuri/jsonschema v1.2.4 // indirect
	github.com/secure-systems-lab/go-securesystemslib v0.9.0 // indirect
	github.com/shirou/gopsutil/v4 v4.25.3 // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/sosodev/duration v1.3.1 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/swaggo/files/v2 v2.0.0 // indirect
	github.com/swaggo/swag v1.8.12 // indirect
	github.com/tiendc/go-deepcopy v1.6.0 // indirect
	github.com/tinylib/msgp v1.2.5 // indirect
	github.com/tklauser/go-sysconf v0.3.15 // indirect
	github.com/tklauser/numcpus v0.10.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/vektah/gqlparser/v2 v2.5.30 // indirect
	github.com/xuri/efp v0.0.1 // indirect
	github.com/xuri/nfp v0.0.1 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.elastic.co/apm/module/apmsql v1.15.0 // indirect
	go.elastic.co/fastjson v1.1.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/collector/component v1.28.1 // indirect
	go.opentelemetry.io/collector/pdata v1.28.1 // indirect
	go.opentelemetry.io/collector/pdata/pprofile v0.122.1 // indirect
	go.opentelemetry.io/collector/semconv v0.123.0 // indirect
	go.opentelemetry.io/otel/metric v1.37.0 // indirect
	go.opentelemetry.io/proto/otlp v1.7.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20250210185358-939b2ce775ac // indirect
	golang.org/x/lint v0.0.0-20201208152925-83fdc39ff7b5 // indirect
	golang.org/x/mod v0.26.0 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/tools v0.35.0 // indirect
	golang.org/x/xerrors v0.0.0-20231012003039-104605ab7028 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250603155806-513f23925822 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250603155806-513f23925822 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gorm.io/driver/mysql v1.5.7 // indirect
	howett.net/plist v0.0.0-20181124034731-591f970eefbb // indirect
)
