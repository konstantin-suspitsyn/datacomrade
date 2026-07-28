package tablesservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	customerrors "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/utils/custom_errors"
)

// GetDatabaseTypeById возвращает активную строку dc.database_type по id.
func (s *TablesService) GetDatabaseTypeById(ctx context.Context, id int64) (tables_model.DcDatabaseType, error) {
	row, err := s.TablesRepository.GetDatabaseTypeById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcDatabaseType{}, fmt.Errorf("dc.database_type id = %d: %w", id, customerrors.ErrNotFound)
		}

		return tables_model.DcDatabaseType{}, fmt.Errorf("get dc.database_type id = %d: %w", id, err)
	}

	return row, nil
}

// GetDatabaseTypes возвращает все активные строки dc.database_type.
func (s *TablesService) GetDatabaseTypes(ctx context.Context) ([]tables_model.DcDatabaseType, error) {
	rows, err := s.TablesRepository.GetDatabaseTypes(ctx)

	if err != nil {
		return nil, fmt.Errorf("get dc.database_type: %w", err)
	}

	return rows, nil
}

// GetDeletedDatabaseTypeById возвращает мягко удалённую строку dc.database_type по id.
func (s *TablesService) GetDeletedDatabaseTypeById(ctx context.Context, id int64) (tables_model.DcDatabaseType, error) {
	row, err := s.TablesRepository.GetDeletedDatabaseTypeById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcDatabaseType{}, fmt.Errorf("deleted dc.database_type id = %d: %w", id, customerrors.ErrNotFound)
		}

		return tables_model.DcDatabaseType{}, fmt.Errorf("get deleted dc.database_type id = %d: %w", id, err)
	}

	return row, nil
}

// GetDeletedDatabaseTypes возвращает все мягко удалённые строки dc.database_type.
func (s *TablesService) GetDeletedDatabaseTypes(ctx context.Context) ([]tables_model.DcDatabaseType, error) {
	rows, err := s.TablesRepository.GetDeletedDatabaseTypes(ctx)

	if err != nil {
		return nil, fmt.Errorf("get deleted dc.database_type: %w", err)
	}

	return rows, nil
}

// CreateDatabaseType вставляет строку dc.database_type и возвращает её целиком.
func (s *TablesService) CreateDatabaseType(ctx context.Context, params tables_model.CreateDatabaseTypeParams) (tables_model.DcDatabaseType, error) {
	row, err := s.TablesRepository.CreateDatabaseType(ctx, params)

	if err != nil {
		return tables_model.DcDatabaseType{}, fmt.Errorf("%w: dc.database_type: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateDatabaseTypeById обновляет активную строку dc.database_type и возвращает её целиком.
//
// Запрос фильтрует по is_deleted = false, поэтому попытка обновить удалённую
// или несуществующую запись даёт sql.ErrNoRows — переводим его в ErrNotFound,
// чтобы api-слой ответил NotFound, а не Internal.
func (s *TablesService) UpdateDatabaseTypeById(ctx context.Context, params tables_model.UpdateDatabaseTypeByIdParams) (tables_model.DcDatabaseType, error) {
	row, err := s.TablesRepository.UpdateDatabaseTypeById(ctx, params)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcDatabaseType{}, fmt.Errorf("dc.database_type id = %d: %w", params.ID, customerrors.ErrNotFound)
		}

		return tables_model.DcDatabaseType{}, fmt.Errorf("%w: dc.database_type id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return row, nil
}

// DeleteDatabaseTypeById мягко удаляет строку dc.database_type.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *TablesService) DeleteDatabaseTypeById(ctx context.Context, id int64) error {
	if _, err := s.GetDatabaseTypeById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.TablesRepository.DeleteDatabaseTypeById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.database_type id = %d: %w", customerrors.ErrDelete, id, err)
	}

	return nil
}

// UndeleteDatabaseTypeById восстанавливает мягко удалённую строку dc.database_type.
// Существование удалённой записи проверяется заранее по той же причине,
// что и в DeleteDatabaseTypeById.
func (s *TablesService) UndeleteDatabaseTypeById(ctx context.Context, id int64) error {
	if _, err := s.GetDeletedDatabaseTypeById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrUndelete, err)
	}

	if err := s.TablesRepository.UndeleteDatabaseTypeById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.database_type id = %d: %w", customerrors.ErrUndelete, id, err)
	}

	return nil
}
