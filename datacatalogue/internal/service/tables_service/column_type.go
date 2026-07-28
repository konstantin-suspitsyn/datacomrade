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

// GetColumnTypes возвращает все активные строки dc.column_type.
func (s *TablesService) GetColumnTypes(ctx context.Context) ([]tables_model.DcColumnType, error) {
	rows, err := s.TablesRepository.GetColumnTypes(ctx)

	if err != nil {
		return nil, fmt.Errorf("get dc.column_type: %w", err)
	}

	return rows, nil
}

// GetDeletedColumnTypeById возвращает мягко удалённую строку dc.column_type по id.
func (s *TablesService) GetDeletedColumnTypeById(ctx context.Context, id int64) (tables_model.DcColumnType, error) {
	row, err := s.TablesRepository.GetDeletedColumnTypeById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcColumnType{}, fmt.Errorf("deleted dc.column_type id = %d: %w", id, customerrors.ErrNotFound)
		}

		return tables_model.DcColumnType{}, fmt.Errorf("get deleted dc.column_type id = %d: %w", id, err)
	}

	return row, nil
}

// GetDeletedColumnTypes возвращает все мягко удалённые строки dc.column_type.
func (s *TablesService) GetDeletedColumnTypes(ctx context.Context) ([]tables_model.DcColumnType, error) {
	rows, err := s.TablesRepository.GetDeletedColumnTypes(ctx)

	if err != nil {
		return nil, fmt.Errorf("get deleted dc.column_type: %w", err)
	}

	return rows, nil
}

// CreateColumnType вставляет строку dc.column_type и возвращает её целиком.
func (s *TablesService) CreateColumnType(ctx context.Context, params tables_model.CreateColumnTypeParams) (tables_model.DcColumnType, error) {
	row, err := s.TablesRepository.CreateColumnType(ctx, params)

	if err != nil {
		return tables_model.DcColumnType{}, fmt.Errorf("%w: dc.column_type: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateColumnTypeById обновляет активную строку dc.column_type и возвращает её целиком.
//
// Запрос фильтрует по is_deleted = false, поэтому попытка обновить удалённую
// или несуществующую запись даёт sql.ErrNoRows — переводим его в ErrNotFound,
// чтобы api-слой ответил NotFound, а не Internal.
func (s *TablesService) UpdateColumnTypeById(ctx context.Context, params tables_model.UpdateColumnTypeByIdParams) (tables_model.DcColumnType, error) {
	row, err := s.TablesRepository.UpdateColumnTypeById(ctx, params)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcColumnType{}, fmt.Errorf("dc.column_type id = %d: %w", params.ID, customerrors.ErrNotFound)
		}

		return tables_model.DcColumnType{}, fmt.Errorf("%w: dc.column_type id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return row, nil
}

// DeleteColumnTypeById мягко удаляет строку dc.column_type.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *TablesService) DeleteColumnTypeById(ctx context.Context, id int64) error {
	if _, err := s.GetColumnTypeById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.TablesRepository.DeleteColumnTypeById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.column_type id = %d: %w", customerrors.ErrDelete, id, err)
	}

	return nil
}

// UndeleteColumnTypeById восстанавливает мягко удалённую строку dc.column_type.
// Существование удалённой записи проверяется заранее по той же причине,
// что и в DeleteColumnTypeById.
func (s *TablesService) UndeleteColumnTypeById(ctx context.Context, id int64) error {
	if _, err := s.GetDeletedColumnTypeById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrUndelete, err)
	}

	if err := s.TablesRepository.UndeleteColumnTypeById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.column_type id = %d: %w", customerrors.ErrUndelete, id, err)
	}

	return nil
}
