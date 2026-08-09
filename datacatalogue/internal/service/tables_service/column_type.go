package tablesservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	customerrors "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/utils/custom_errors"
)

// GetColumnTypeById возвращает активную строку dc.column_type по id.
func (s *TablesService) GetColumnTypeById(ctx context.Context, id int64) (tables_model.DcColumnType, error) {
	row, err := s.TablesRepository.GetColumnTypeById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcColumnType{}, fmt.Errorf("dc.column_type id = %d: %w", id, customerrors.ErrNotFound)
		}

		return tables_model.DcColumnType{}, fmt.Errorf("get dc.column_type id = %d: %w", id, err)
	}

	return row, nil
}

// GetColumnTypes возвращает страницу строк dc.column_type и её счётчики.
func (s *TablesService) GetColumnTypes(ctx context.Context, params tables_model.GetColumnTypesParams) ([]tables_model.DcColumnType, tables_model.CountGetColumnTypesRow, error) {
	count, err := s.TablesRepository.CountGetColumnTypes(ctx, params.PageLimit)
	if err != nil {
		return nil, tables_model.CountGetColumnTypesRow{}, fmt.Errorf("count dc.column_type: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcColumnType{}, count, nil
	}

	rows, err := s.TablesRepository.GetColumnTypes(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetColumnTypesRow{}, fmt.Errorf("get dc.column_type page: %w", err)
	}

	return rows, count, nil
}

// GetColumnTypesSearchName возвращает страницу строк dc.column_type и её счётчики.
func (s *TablesService) GetColumnTypesSearchName(ctx context.Context, params tables_model.GetColumnTypesSearchNameParams) ([]tables_model.DcColumnType, tables_model.CountGetColumnTypesSearchNameRow, error) {
	countParams := tables_model.CountGetColumnTypesSearchNameParams{
		PageLimit:  params.PageLimit,
		SearchName: params.SearchName,
	}
	count, err := s.TablesRepository.CountGetColumnTypesSearchName(ctx, countParams)
	if err != nil {
		return nil, tables_model.CountGetColumnTypesSearchNameRow{}, fmt.Errorf("count dc.column_type: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcColumnType{}, count, nil
	}

	rows, err := s.TablesRepository.GetColumnTypesSearchName(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetColumnTypesSearchNameRow{}, fmt.Errorf("get dc.column_type page: %w", err)
	}

	return rows, count, nil
}

// GetColumnTypesDeleted возвращает страницу строк dc.column_type и её счётчики.
func (s *TablesService) GetColumnTypesDeleted(ctx context.Context, params tables_model.GetColumnTypesDeletedParams) ([]tables_model.DcColumnType, tables_model.CountGetColumnTypesDeletedRow, error) {
	count, err := s.TablesRepository.CountGetColumnTypesDeleted(ctx, params.PageLimit)
	if err != nil {
		return nil, tables_model.CountGetColumnTypesDeletedRow{}, fmt.Errorf("count dc.column_type: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcColumnType{}, count, nil
	}

	rows, err := s.TablesRepository.GetColumnTypesDeleted(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetColumnTypesDeletedRow{}, fmt.Errorf("get dc.column_type page: %w", err)
	}

	return rows, count, nil
}

// CreateColumnType вставляет строку dc.column_type и возвращает её целиком.
func (s *TablesService) CreateColumnType(ctx context.Context, params tables_model.CreateColumnTypeParams) (tables_model.DcColumnType, error) {
	row, err := s.TablesRepository.CreateColumnType(ctx, params)

	if err != nil {
		return tables_model.DcColumnType{}, fmt.Errorf("%w: dc.column_type: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateColumnTypeById обновляет строку dc.column_type.
//
// Запрос — :exec и не сообщает число затронутых строк, поэтому
// существование активной записи проверяется заранее.
func (s *TablesService) UpdateColumnTypeById(ctx context.Context, params tables_model.UpdateColumnTypeByIdParams) error {
	if _, err := s.GetColumnTypeById(ctx, params.ID); err != nil {
		return fmt.Errorf("%w: dc.column_type: %w", customerrors.ErrUpdate, err)
	}

	if err := s.TablesRepository.UpdateColumnTypeById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.column_type id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return nil
}

// DeleteColumnTypeById мягко удаляет строку dc.column_type.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *TablesService) DeleteColumnTypeById(ctx context.Context, params tables_model.DeleteColumnTypeByIdParams) error {
	if _, err := s.GetColumnTypeById(ctx, params.ID); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.TablesRepository.DeleteColumnTypeById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.column_type id = %d: %w", customerrors.ErrDelete, params.ID, err)
	}

	return nil
}

// UndeleteColumnTypeById восстанавливает мягко удалённую строку dc.column_type.
func (s *TablesService) UndeleteColumnTypeById(ctx context.Context, params tables_model.UndeleteColumnTypeByIdParams) error {
	if err := s.TablesRepository.UndeleteColumnTypeById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.column_type id = %d: %w", customerrors.ErrUndelete, params.ID, err)
	}

	return nil
}
