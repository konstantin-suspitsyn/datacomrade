// Package queryengine оборачивает gRPC-соединение с query_engine в один
// переиспользуемый клиент на весь процесс Shepherd.
package queryengine

import (
	"fmt"

	queryenginev1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/query_engine/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client держит одно gRPC-соединение с query_engine и типизированный клиент
// сервиса поверх него.
type Client struct {
	conn *grpc.ClientConn

	QueryEngine queryenginev1.QueryEngineServiceClient
}

// Dial открывает gRPC-соединение с query_engine по addr (host:port).
// Соединение живёт на весь процесс Shepherd, не открывается заново на запрос.
func Dial(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("dial query_engine at %s: %w", addr, err)
	}

	return &Client{
		conn:        conn,
		QueryEngine: queryenginev1.NewQueryEngineServiceClient(conn),
	}, nil
}

// Close закрывает gRPC-соединение.
func (c *Client) Close() error {
	return c.conn.Close()
}
