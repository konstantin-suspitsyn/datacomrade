package authlogicservice

import (
	"context"
	"fmt"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/auth_logic"
)

// GetTableIdsByExternalUserIdAndRoles оборачивает sqlc-запрос auth_logic.GetTableIdsByExternalUserIdAndRoles.
func (s *AuthLogicService) GetTableIdsByExternalUserIdAndRoles(ctx context.Context, params auth_logic.GetTableIdsByExternalUserIdAndRolesParams) ([]int64, error) {
	rows, err := s.AuthLogicRepository.GetTableIdsByExternalUserIdAndRoles(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("GetTableIdsByExternalUserIdAndRoles: %w", err)
	}

	return rows, nil
}

// GetTableIdsByUserIdAndRoles оборачивает sqlc-запрос auth_logic.GetTableIdsByUserIdAndRoles.
func (s *AuthLogicService) GetTableIdsByUserIdAndRoles(ctx context.Context, params auth_logic.GetTableIdsByUserIdAndRolesParams) ([]int64, error) {
	rows, err := s.AuthLogicRepository.GetTableIdsByUserIdAndRoles(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("GetTableIdsByUserIdAndRoles: %w", err)
	}

	return rows, nil
}
