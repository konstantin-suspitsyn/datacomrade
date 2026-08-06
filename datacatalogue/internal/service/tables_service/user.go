package tablesservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	customerrors "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/utils/custom_errors"
)

// GetUsers возвращает страницу строк dc.user и её счётчики.
func (s *TablesService) GetUsers(ctx context.Context, params tables_model.GetUsersParams) ([]tables_model.DcUser, tables_model.CountGetUsersRow, error) {
	count, err := s.TablesRepository.CountGetUsers(ctx, params.PageLimit)
	if err != nil {
		return nil, tables_model.CountGetUsersRow{}, fmt.Errorf("count dc.user: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcUser{}, count, nil
	}

	rows, err := s.TablesRepository.GetUsers(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetUsersRow{}, fmt.Errorf("get dc.user page: %w", err)
	}

	return rows, count, nil
}

// GetUserByExternalId возвращает активную строку dc.user по уникальной колонке external_id.
func (s *TablesService) GetUserByExternalId(ctx context.Context, externalID uuid.UUID) (tables_model.DcUser, error) {
	row, err := s.TablesRepository.GetUserByExternalId(ctx, externalID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcUser{}, fmt.Errorf("dc.user external_id = %v: %w", externalID, customerrors.ErrNotFound)
		}

		return tables_model.DcUser{}, fmt.Errorf("get dc.user external_id = %v: %w", externalID, err)
	}

	return row, nil
}

// CreateUser вставляет строку dc.user и возвращает её целиком.
func (s *TablesService) CreateUser(ctx context.Context, params tables_model.CreateUserParams) (tables_model.DcUser, error) {
	row, err := s.TablesRepository.CreateUser(ctx, params)

	if err != nil {
		return tables_model.DcUser{}, fmt.Errorf("%w: dc.user: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}
