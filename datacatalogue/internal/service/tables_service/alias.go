package tablesservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	customerrors "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/utils/custom_errors"
)

// GetAliasById возвращает активную строку dc.alias по id.
func (s *TablesService) GetAliasById(ctx context.Context, id int64) (tables_model.DcAlias, error) {
	row, err := s.TablesRepository.GetAliasById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcAlias{}, fmt.Errorf("dc.alias id = %d: %w", id, customerrors.ErrNotFound)
		}

		return tables_model.DcAlias{}, fmt.Errorf("get dc.alias id = %d: %w", id, err)
	}

	return row, nil
}

// GetAliasesDeleted возвращает страницу строк dc.alias и её счётчики.
func (s *TablesService) GetAliasesDeleted(ctx context.Context, params tables_model.GetAliasesDeletedParams) ([]tables_model.DcAlias, tables_model.CountGetAliasesDeletedRow, error) {
	count, err := s.TablesRepository.CountGetAliasesDeleted(ctx, params.PageLimit)
	if err != nil {
		return nil, tables_model.CountGetAliasesDeletedRow{}, fmt.Errorf("count dc.alias: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcAlias{}, count, nil
	}

	rows, err := s.TablesRepository.GetAliasesDeleted(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetAliasesDeletedRow{}, fmt.Errorf("get dc.alias page: %w", err)
	}

	return rows, count, nil
}

// GetAliases возвращает страницу строк dc.alias и её счётчики.
func (s *TablesService) GetAliases(ctx context.Context, params tables_model.GetAliasesParams) ([]tables_model.DcAlias, tables_model.CountGetAliasesRow, error) {
	count, err := s.TablesRepository.CountGetAliases(ctx, params.PageLimit)
	if err != nil {
		return nil, tables_model.CountGetAliasesRow{}, fmt.Errorf("count dc.alias: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcAlias{}, count, nil
	}

	rows, err := s.TablesRepository.GetAliases(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetAliasesRow{}, fmt.Errorf("get dc.alias page: %w", err)
	}

	return rows, count, nil
}

// CreateAlias вставляет строку dc.alias и возвращает её целиком.
func (s *TablesService) CreateAlias(ctx context.Context, params tables_model.CreateAliasParams) (tables_model.DcAlias, error) {
	row, err := s.TablesRepository.CreateAlias(ctx, params)

	if err != nil {
		return tables_model.DcAlias{}, fmt.Errorf("%w: dc.alias: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateAliasById обновляет строку dc.alias.
//
// Запрос — :exec и не сообщает число затронутых строк, поэтому
// существование активной записи проверяется заранее.
func (s *TablesService) UpdateAliasById(ctx context.Context, params tables_model.UpdateAliasByIdParams) error {
	if _, err := s.GetAliasById(ctx, params.ID); err != nil {
		return fmt.Errorf("%w: dc.alias: %w", customerrors.ErrUpdate, err)
	}

	if err := s.TablesRepository.UpdateAliasById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.alias id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return nil
}

// DeleteAliasById мягко удаляет строку dc.alias.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *TablesService) DeleteAliasById(ctx context.Context, params tables_model.DeleteAliasByIdParams) error {
	if _, err := s.GetAliasById(ctx, params.ID); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.TablesRepository.DeleteAliasById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.alias id = %d: %w", customerrors.ErrDelete, params.ID, err)
	}

	return nil
}

// UndeleteAliasById восстанавливает мягко удалённую строку dc.alias.
func (s *TablesService) UndeleteAliasById(ctx context.Context, params tables_model.UndeleteAliasByIdParams) error {
	if err := s.TablesRepository.UndeleteAliasById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.alias id = %d: %w", customerrors.ErrUndelete, params.ID, err)
	}

	return nil
}
