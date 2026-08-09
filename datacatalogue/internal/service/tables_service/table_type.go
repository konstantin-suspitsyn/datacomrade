package tablesservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	customerrors "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/utils/custom_errors"
)

// GetTableTypeById возвращает активную строку dc.table_type по id.
func (s *TablesService) GetTableTypeById(ctx context.Context, id int64) (tables_model.DcTableType, error) {
	row, err := s.TablesRepository.GetTableTypeById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcTableType{}, fmt.Errorf("dc.table_type id = %d: %w", id, customerrors.ErrNotFound)
		}

		return tables_model.DcTableType{}, fmt.Errorf("get dc.table_type id = %d: %w", id, err)
	}

	return row, nil
}

// GetTableTypes возвращает страницу строк dc.table_type и её счётчики.
func (s *TablesService) GetTableTypes(ctx context.Context, params tables_model.GetTableTypesParams) ([]tables_model.DcTableType, tables_model.CountGetTableTypesRow, error) {
	count, err := s.TablesRepository.CountGetTableTypes(ctx, params.PageLimit)
	if err != nil {
		return nil, tables_model.CountGetTableTypesRow{}, fmt.Errorf("count dc.table_type: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcTableType{}, count, nil
	}

	rows, err := s.TablesRepository.GetTableTypes(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetTableTypesRow{}, fmt.Errorf("get dc.table_type page: %w", err)
	}

	return rows, count, nil
}

// GetTableTypesDeleted возвращает страницу строк dc.table_type и её счётчики.
func (s *TablesService) GetTableTypesDeleted(ctx context.Context, params tables_model.GetTableTypesDeletedParams) ([]tables_model.DcTableType, tables_model.CountGetTableTypesDeletedRow, error) {
	count, err := s.TablesRepository.CountGetTableTypesDeleted(ctx, params.PageLimit)
	if err != nil {
		return nil, tables_model.CountGetTableTypesDeletedRow{}, fmt.Errorf("count dc.table_type: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcTableType{}, count, nil
	}

	rows, err := s.TablesRepository.GetTableTypesDeleted(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetTableTypesDeletedRow{}, fmt.Errorf("get dc.table_type page: %w", err)
	}

	return rows, count, nil
}

// CreateTableType вставляет строку dc.table_type и возвращает её целиком.
func (s *TablesService) CreateTableType(ctx context.Context, params tables_model.CreateTableTypeParams) (tables_model.DcTableType, error) {
	row, err := s.TablesRepository.CreateTableType(ctx, params)

	if err != nil {
		return tables_model.DcTableType{}, fmt.Errorf("%w: dc.table_type: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateTableTypeById обновляет строку dc.table_type.
//
// Запрос — :exec и не сообщает число затронутых строк, поэтому
// существование активной записи проверяется заранее.
func (s *TablesService) UpdateTableTypeById(ctx context.Context, params tables_model.UpdateTableTypeByIdParams) error {
	if _, err := s.GetTableTypeById(ctx, params.ID); err != nil {
		return fmt.Errorf("%w: dc.table_type: %w", customerrors.ErrUpdate, err)
	}

	if err := s.TablesRepository.UpdateTableTypeById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.table_type id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return nil
}

// DeleteTableTypeById мягко удаляет строку dc.table_type.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *TablesService) DeleteTableTypeById(ctx context.Context, params tables_model.DeleteTableTypeByIdParams) error {
	if _, err := s.GetTableTypeById(ctx, params.ID); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.TablesRepository.DeleteTableTypeById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.table_type id = %d: %w", customerrors.ErrDelete, params.ID, err)
	}

	return nil
}

// UndeleteTableTypeById восстанавливает мягко удалённую строку dc.table_type.
func (s *TablesService) UndeleteTableTypeById(ctx context.Context, params tables_model.UndeleteTableTypeByIdParams) error {
	if err := s.TablesRepository.UndeleteTableTypeById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.table_type id = %d: %w", customerrors.ErrUndelete, params.ID, err)
	}

	return nil
}
