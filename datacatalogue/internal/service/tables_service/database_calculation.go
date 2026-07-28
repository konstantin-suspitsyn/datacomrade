package tablesservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	customerrors "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/utils/custom_errors"
)

// GetDatabaseCalculationById возвращает активную строку dc.database_calculation по id.
func (s *TablesService) GetDatabaseCalculationById(ctx context.Context, id int64) (tables_model.DcDatabaseCalculation, error) {
	row, err := s.TablesRepository.GetDatabaseCalculationById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcDatabaseCalculation{}, fmt.Errorf("dc.database_calculation id = %d: %w", id, customerrors.ErrNotFound)
		}

		return tables_model.DcDatabaseCalculation{}, fmt.Errorf("get dc.database_calculation id = %d: %w", id, err)
	}

	return row, nil
}

// GetDatabaseCalculations возвращает все активные строки dc.database_calculation.
func (s *TablesService) GetDatabaseCalculations(ctx context.Context) ([]tables_model.DcDatabaseCalculation, error) {
	rows, err := s.TablesRepository.GetDatabaseCalculations(ctx)

	if err != nil {
		return nil, fmt.Errorf("get dc.database_calculation: %w", err)
	}

	return rows, nil
}

// GetDeletedDatabaseCalculationById возвращает мягко удалённую строку dc.database_calculation по id.
func (s *TablesService) GetDeletedDatabaseCalculationById(ctx context.Context, id int64) (tables_model.DcDatabaseCalculation, error) {
	row, err := s.TablesRepository.GetDeletedDatabaseCalculationById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcDatabaseCalculation{}, fmt.Errorf("deleted dc.database_calculation id = %d: %w", id, customerrors.ErrNotFound)
		}

		return tables_model.DcDatabaseCalculation{}, fmt.Errorf("get deleted dc.database_calculation id = %d: %w", id, err)
	}

	return row, nil
}

// GetDeletedDatabaseCalculations возвращает все мягко удалённые строки dc.database_calculation.
func (s *TablesService) GetDeletedDatabaseCalculations(ctx context.Context) ([]tables_model.DcDatabaseCalculation, error) {
	rows, err := s.TablesRepository.GetDeletedDatabaseCalculations(ctx)

	if err != nil {
		return nil, fmt.Errorf("get deleted dc.database_calculation: %w", err)
	}

	return rows, nil
}

// CreateDatabaseCalculation вставляет строку dc.database_calculation и возвращает её целиком.
func (s *TablesService) CreateDatabaseCalculation(ctx context.Context, params tables_model.CreateDatabaseCalculationParams) (tables_model.DcDatabaseCalculation, error) {
	row, err := s.TablesRepository.CreateDatabaseCalculation(ctx, params)

	if err != nil {
		return tables_model.DcDatabaseCalculation{}, fmt.Errorf("%w: dc.database_calculation: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateDatabaseCalculationById обновляет активную строку dc.database_calculation и возвращает её целиком.
//
// Запрос фильтрует по is_deleted = false, поэтому попытка обновить удалённую
// или несуществующую запись даёт sql.ErrNoRows — переводим его в ErrNotFound,
// чтобы api-слой ответил NotFound, а не Internal.
func (s *TablesService) UpdateDatabaseCalculationById(ctx context.Context, params tables_model.UpdateDatabaseCalculationByIdParams) (tables_model.DcDatabaseCalculation, error) {
	row, err := s.TablesRepository.UpdateDatabaseCalculationById(ctx, params)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcDatabaseCalculation{}, fmt.Errorf("dc.database_calculation id = %d: %w", params.ID, customerrors.ErrNotFound)
		}

		return tables_model.DcDatabaseCalculation{}, fmt.Errorf("%w: dc.database_calculation id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return row, nil
}

// DeleteDatabaseCalculationById мягко удаляет строку dc.database_calculation.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *TablesService) DeleteDatabaseCalculationById(ctx context.Context, id int64) error {
	if _, err := s.GetDatabaseCalculationById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.TablesRepository.DeleteDatabaseCalculationById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.database_calculation id = %d: %w", customerrors.ErrDelete, id, err)
	}

	return nil
}

// UndeleteDatabaseCalculationById восстанавливает мягко удалённую строку dc.database_calculation.
// Существование удалённой записи проверяется заранее по той же причине,
// что и в DeleteDatabaseCalculationById.
func (s *TablesService) UndeleteDatabaseCalculationById(ctx context.Context, id int64) error {
	if _, err := s.GetDeletedDatabaseCalculationById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrUndelete, err)
	}

	if err := s.TablesRepository.UndeleteDatabaseCalculationById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.database_calculation id = %d: %w", customerrors.ErrUndelete, id, err)
	}

	return nil
}
