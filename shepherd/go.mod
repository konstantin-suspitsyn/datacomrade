module github.com/konstantin-suspitsyn/datacomrade/shepherd

go 1.26.5

require (
	github.com/go-chi/chi/v5 v5.3.1
	github.com/go-chi/cors v1.2.2
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/konstantin-suspitsyn/datacomrade/platform v0.0.0-00010101000000-000000000000
	github.com/konstantin-suspitsyn/datacomrade/shared v0.0.0-00010101000000-000000000000
	github.com/lestrrat-go/jwx/v3 v3.2.0
	github.com/oapi-codegen/runtime v1.6.0
	github.com/redis/go-redis/v9 v9.21.0
	go.uber.org/zap v1.27.1
	google.golang.org/grpc v1.83.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.4.1 // indirect
	github.com/goccy/go-json v0.10.6 // indirect
	github.com/lestrrat-go/blackmagic v1.0.4 // indirect
	github.com/lestrrat-go/dsig v1.3.0 // indirect
	github.com/lestrrat-go/dsig-secp256k1 v1.0.0 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/httprc/v3 v3.0.6 // indirect
	github.com/lestrrat-go/option/v2 v2.0.0 // indirect
	github.com/segmentio/asm v1.2.1 // indirect
	github.com/valyala/fastjson v1.6.10 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/crypto v0.54.0 // indirect
	golang.org/x/net v0.56.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.40.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260526163538-3dc84a4a5aaa // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace (
	github.com/konstantin-suspitsyn/datacomrade/platform => ../platform
	github.com/konstantin-suspitsyn/datacomrade/shared => ../shared
)
