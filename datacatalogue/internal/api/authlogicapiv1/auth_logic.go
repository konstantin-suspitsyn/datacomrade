package authlogicapiv1

import (
	"context"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/api/apierror"
	authlogicconv "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter/authlogic"
	authlogicvalidation "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation/authlogic"
	authlogicv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/auth_logic/v1"
)

// GetTableIdsByExternalUserIdAndRoles оборачивает sqlc-запрос auth_logic.GetTableIdsByExternalUserIdAndRoles.
func (a *AuthLogicApiV1) GetTableIdsByExternalUserIdAndRoles(ctx context.Context, req *authlogicv1.GetTableIdsByExternalUserIdAndRolesRequest) (*authlogicv1.GetTableIdsByExternalUserIdAndRolesResponse, error) {
	if err := authlogicvalidation.ValidateGetTableIdsByExternalUserIdAndRoles(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	tableIds, err := a.services.AuthLogicService.GetTableIdsByExternalUserIdAndRoles(ctx, authlogicconv.ToGetTableIdsByExternalUserIdAndRolesParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return authlogicconv.GetTableIdsByExternalUserIdAndRolesToProto(tableIds), nil
}

// GetTableIdsByUserIdAndRoles оборачивает sqlc-запрос auth_logic.GetTableIdsByUserIdAndRoles.
func (a *AuthLogicApiV1) GetTableIdsByUserIdAndRoles(ctx context.Context, req *authlogicv1.GetTableIdsByUserIdAndRolesRequest) (*authlogicv1.GetTableIdsByUserIdAndRolesResponse, error) {
	if err := authlogicvalidation.ValidateGetTableIdsByUserIdAndRoles(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	tableIds, err := a.services.AuthLogicService.GetTableIdsByUserIdAndRoles(ctx, authlogicconv.ToGetTableIdsByUserIdAndRolesParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return authlogicconv.GetTableIdsByUserIdAndRolesToProto(tableIds), nil
}
