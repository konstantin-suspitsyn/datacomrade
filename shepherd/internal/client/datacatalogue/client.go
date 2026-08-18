// Package datacatalogue оборачивает gRPC-соединение с Metadata Service
// (репозиторий datacatalogue) в один переиспользуемый клиент на весь процесс.
package datacatalogue

import (
	"fmt"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	userv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user/v1"
	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client держит одно gRPC-соединение и типизированные клиенты сервисов
// Metadata Service поверх него.
type Client struct {
	conn *grpc.ClientConn

	User         userv1.UserServiceClient
	Host         tablesv1.HostServiceClient
	DatabaseType tablesv1.DatabaseTypeServiceClient
	DatabaseCat  tablesv1.DatabaseCatServiceClient

	Alias                tablesv1.AliasServiceClient
	SchemaCat            tablesv1.SchemaCatServiceClient
	TableCat             tablesv1.TableCatServiceClient
	TableType            tablesv1.TableTypeServiceClient
	ColumnCat            tablesv1.ColumnCatServiceClient
	ColumnType           tablesv1.ColumnTypeServiceClient
	DomainCat            tablesv1.DomainCatServiceClient
	CalculationType      tablesv1.CalculationTypeServiceClient
	DatabaseCalculation  tablesv1.DatabaseCalculationServiceClient
	FollowingCalculation tablesv1.FollowingCalculationServiceClient
	GroupLevels          tablesv1.GroupLevelsServiceClient
	HasToGroup           tablesv1.HasToGroupServiceClient

	UserDomainRoles userdomainrolesv1.UserDomainRolesServiceClient
	UserTableRoles  userdomainrolesv1.UserTableRolesServiceClient
}

// Dial открывает gRPC-соединение с datacatalogue по addr (host:port).
// Соединение живёт на весь процесс Shepherd, не открывается заново на запрос.
func Dial(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("dial datacatalogue at %s: %w", addr, err)
	}

	return &Client{
		conn:         conn,
		User:         userv1.NewUserServiceClient(conn),
		Host:         tablesv1.NewHostServiceClient(conn),
		DatabaseType: tablesv1.NewDatabaseTypeServiceClient(conn),
		DatabaseCat:  tablesv1.NewDatabaseCatServiceClient(conn),

		Alias:                tablesv1.NewAliasServiceClient(conn),
		SchemaCat:            tablesv1.NewSchemaCatServiceClient(conn),
		TableCat:             tablesv1.NewTableCatServiceClient(conn),
		TableType:            tablesv1.NewTableTypeServiceClient(conn),
		ColumnCat:            tablesv1.NewColumnCatServiceClient(conn),
		ColumnType:           tablesv1.NewColumnTypeServiceClient(conn),
		DomainCat:            tablesv1.NewDomainCatServiceClient(conn),
		CalculationType:      tablesv1.NewCalculationTypeServiceClient(conn),
		DatabaseCalculation:  tablesv1.NewDatabaseCalculationServiceClient(conn),
		FollowingCalculation: tablesv1.NewFollowingCalculationServiceClient(conn),
		GroupLevels:          tablesv1.NewGroupLevelsServiceClient(conn),
		HasToGroup:           tablesv1.NewHasToGroupServiceClient(conn),

		UserDomainRoles: userdomainrolesv1.NewUserDomainRolesServiceClient(conn),
		UserTableRoles:  userdomainrolesv1.NewUserTableRolesServiceClient(conn),
	}, nil
}

// Close закрывает gRPC-соединение.
func (c *Client) Close() error {
	return c.conn.Close()
}
