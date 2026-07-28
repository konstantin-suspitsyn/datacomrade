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

// GetTableTypes возвращает все активные строки dc.table_type.
func (s *TablesService) GetTableTypes(ctx context.Context) ([]tables_model.DcTableType, error) {
	rows, err := s.TablesRepository.GetTableTypes(ctx)

	if err != nil {
		return nil, fmt.Errorf("get dc.table_type: %w", err)
	}

	return rows, nil
}

// GetDeletedTableTypeById возвращает мягко удалённую строку dc.table_type по id.
func (s *TablesService) GetDeletedTableTypeById(ctx context.Context, id int64) (tables_model.DcTableType, error) {
	row, err := s.TablesRepository.GetDeletedTableTypeById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcTableType{}, fmt.Errorf("deleted dc.table_type id = %d: %w", id, customerrors.ErrNotFound)
		}

		return tables_model.DcTableType{}, fmt.Errorf("get deleted dc.table_type id = %d: %w", id, err)
	}

	return row, nil
}

// GetDeletedTableTypes возвращает все мягко удалённые строки dc.table_type.
func (s *TablesService) GetDeletedTableTypes(ctx context.Context) ([]tables_model.DcTableType, error) {
	rows, err := s.TablesRepository.GetDeletedTableTypes(ctx)

	if err != nil {
		return nil, fmt.Errorf("get deleted dc.table_type: %w", err)
	}

	return rows, nil
}

// CreateTableType вставляет строку dc.table_type и возвращает её целиком.
func (s *TablesService) CreateTableType(ctx context.Context, params tables_model.CreateTableTypeParams) (tables_model.DcTableType, error) {
	row, err := s.TablesRepository.CreateTableType(ctx, params)

	if err != nil {
		return tables_model.DcTableType{}, fmt.Errorf("%w: dc.table_type: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateTableTypeById обновляет активную строку dc.table_type и возвращает её целиком.
//
// Запрос фильтрует по is_deleted = false, поэтому попытка обновить удалённую
// или несуществующую запись даёт sql.ErrNoRows — переводим его в ErrNotFound,
// чтобы api-слой ответил NotFound, а не Internal.
func (s *TablesService) UpdateTableTypeById(ctx context.Context, params tables_model.UpdateTableTypeByIdParams) (tables_model.DcTableType, error) {
	row, err := s.TablesRepository.UpdateTableTypeById(ctx, params)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcTableType{}, fmt.Errorf("dc.table_type id = %d: %w", params.ID, customerrors.ErrNotFound)
		}

		return tables_model.DcTableType{}, fmt.Errorf("%w: dc.table_type id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return row, nil
}

// DeleteTableTypeById мягко удаляет строку dc.table_type.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *TablesService) DeleteTableTypeById(ctx context.Context, id int64) error {
	if _, err := s.GetTableTypeById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.TablesRepository.DeleteTableTypeById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.table_type id = %d: %w", customerrors.ErrDelete, id, err)
	}

	return nil
}

// UndeleteTableTypeById восстанавливает мягко удалённую строку dc.table_type.
// Существование удалённой записи проверяется заранее по той же причине,
// что и в DeleteTableTypeById.
func (s *TablesService) UndeleteTableTypeById(ctx context.Context, id int64) error {
	if _, err := s.GetDeletedTableTypeById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrUndelete, err)
	}

	if err := s.TablesRepository.UndeleteTableTypeById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.table_type id = %d: %w", customerrors.ErrUndelete, id, err)
	}

	return nil
}
