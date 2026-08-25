package tablesservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	customerrors "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/utils/custom_errors"
)

// GetTablesJoinsById возвращает активную строку dc.tables_joins по id.
func (s *TablesService) GetTablesJoinsById(ctx context.Context, id int64) (tables_model.GetTablesJoinsByIdRow, error) {
	row, err := s.TablesRepository.GetTablesJoinsById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.GetTablesJoinsByIdRow{}, fmt.Errorf("dc.tables_joins id = %d: %w", id, customerrors.ErrNotFound)
		}

		return tables_model.GetTablesJoinsByIdRow{}, fmt.Errorf("get dc.tables_joins id = %d: %w", id, err)
	}

	return row, nil
}

// GetTablesJoins возвращает страницу строк dc.tables_joins и её счётчики.
func (s *TablesService) GetTablesJoins(ctx context.Context, params tables_model.GetTablesJoinsParams) ([]tables_model.GetTablesJoinsRow, tables_model.CountGetTablesJoinsRow, error) {
	countParams := tables_model.CountGetTablesJoinsParams{
		PageLimit:         params.PageLimit,
		TableCatLeftName:  params.TableCatLeftName,
		TableCatRightName: params.TableCatRightName,
	}
	count, err := s.TablesRepository.CountGetTablesJoins(ctx, countParams)
	if err != nil {
		return nil, tables_model.CountGetTablesJoinsRow{}, fmt.Errorf("count dc.tables_joins: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.GetTablesJoinsRow{}, count, nil
	}

	rows, err := s.TablesRepository.GetTablesJoins(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetTablesJoinsRow{}, fmt.Errorf("get dc.tables_joins page: %w", err)
	}

	return rows, count, nil
}

// CreateTablesJoins вставляет строку dc.tables_joins и возвращает её целиком.
func (s *TablesService) CreateTablesJoins(ctx context.Context, params tables_model.CreateTablesJoinsParams) (tables_model.DcTablesJoin, error) {
	row, err := s.TablesRepository.CreateTablesJoins(ctx, params)

	if err != nil {
		return tables_model.DcTablesJoin{}, fmt.Errorf("%w: dc.tables_joins: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateTablesJoinsById обновляет строку dc.tables_joins.
//
// Запрос — :exec и не сообщает число затронутых строк, поэтому
// существование активной записи проверяется заранее.
func (s *TablesService) UpdateTablesJoinsById(ctx context.Context, params tables_model.UpdateTablesJoinsByIdParams) error {
	if _, err := s.GetTablesJoinsById(ctx, params.ID); err != nil {
		return fmt.Errorf("%w: dc.tables_joins: %w", customerrors.ErrUpdate, err)
	}

	if err := s.TablesRepository.UpdateTablesJoinsById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.tables_joins id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return nil
}

// DeleteTablesJoinsById мягко удаляет строку dc.tables_joins.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *TablesService) DeleteTablesJoinsById(ctx context.Context, params tables_model.DeleteTablesJoinsByIdParams) error {
	if _, err := s.GetTablesJoinsById(ctx, params.ID); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.TablesRepository.DeleteTablesJoinsById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.tables_joins id = %d: %w", customerrors.ErrDelete, params.ID, err)
	}

	return nil
}

// UndeleteTablesJoinsById восстанавливает мягко удалённую строку dc.tables_joins.
func (s *TablesService) UndeleteTablesJoinsById(ctx context.Context, id int64) error {
	if err := s.TablesRepository.UndeleteTablesJoinsById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.tables_joins id = %d: %w", customerrors.ErrUndelete, id, err)
	}

	return nil
}
