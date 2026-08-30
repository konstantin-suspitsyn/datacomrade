module github.com/konstantin-suspitsyn/datacomrade/query_engine

go 1.26.5

require (
	github.com/joho/godotenv v1.5.1
	github.com/konstantin-suspitsyn/datacomrade/shared v0.0.0-00010101000000-000000000000
	github.com/lib/pq v1.12.3
	google.golang.org/grpc v1.83.0
	google.golang.org/protobuf v1.36.11
)

require (
	golang.org/x/net v0.56.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.40.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260526163538-3dc84a4a5aaa // indirect
)

replace github.com/konstantin-suspitsyn/datacomrade/shared => ../shared
