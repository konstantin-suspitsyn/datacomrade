package authlogic

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/auth_logic"
	authlogicv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/auth_logic/v1"
)

// ToGetTableIdsByExternalUserIdAndRolesParams собирает параметры auth_logic.GetTableIdsByExternalUserIdAndRoles из запроса gRPC.
func ToGetTableIdsByExternalUserIdAndRolesParams(req *authlogicv1.GetTableIdsByExternalUserIdAndRolesRequest) auth_logic.GetTableIdsByExternalUserIdAndRolesParams {
	return auth_logic.GetTableIdsByExternalUserIdAndRolesParams{
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
		Name:       req.GetName(),
	}
}

// GetTableIdsByExternalUserIdAndRolesToProto переводит результат auth_logic.GetTableIdsByExternalUserIdAndRoles в ответ gRPC.
func GetTableIdsByExternalUserIdAndRolesToProto(tableIds []int64) *authlogicv1.GetTableIdsByExternalUserIdAndRolesResponse {
	return &authlogicv1.GetTableIdsByExternalUserIdAndRolesResponse{TableIds: tableIds}
}

// ToGetTableIdsByUserIdAndRolesParams собирает параметры auth_logic.GetTableIdsByUserIdAndRoles из запроса gRPC.
func ToGetTableIdsByUserIdAndRolesParams(req *authlogicv1.GetTableIdsByUserIdAndRolesRequest) auth_logic.GetTableIdsByUserIdAndRolesParams {
	return auth_logic.GetTableIdsByUserIdAndRolesParams{
		UserID: req.GetUserId(),
		Name:   req.GetName(),
	}
}

// GetTableIdsByUserIdAndRolesToProto переводит результат auth_logic.GetTableIdsByUserIdAndRoles в ответ gRPC.
func GetTableIdsByUserIdAndRolesToProto(tableIds []int64) *authlogicv1.GetTableIdsByUserIdAndRolesResponse {
	return &authlogicv1.GetTableIdsByUserIdAndRolesResponse{TableIds: tableIds}
}
