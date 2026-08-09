// Package datacatalogue оборачивает gRPC-соединение с Metadata Service
// (репозиторий datacatalogue) в один переиспользуемый клиент на весь процесс.
package datacatalogue

import (
	"fmt"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	userv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client держит одно gRPC-соединение и типизированные клиенты сервисов
// Metadata Service поверх него.
type Client struct {
	conn *grpc.ClientConn

	User userv1.UserServiceClient
	Host tablesv1.HostServiceClient
}

// Dial открывает gRPC-соединение с datacatalogue по addr (host:port).
// Соединение живёт на весь процесс Shepherd, не открывается заново на запрос.
func Dial(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("dial datacatalogue at %s: %w", addr, err)
	}

	return &Client{
		conn: conn,
		User: userv1.NewUserServiceClient(conn),
		Host: tablesv1.NewHostServiceClient(conn),
	}, nil
}

// Close закрывает gRPC-соединение.
func (c *Client) Close() error {
	return c.conn.Close()
}
