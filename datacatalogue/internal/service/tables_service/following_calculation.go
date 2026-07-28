package tablesservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	customerrors "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/utils/custom_errors"
)

// GetFollowingCalculationById возвращает активную строку dc.following_calculation по id.
func (s *TablesService) GetFollowingCalculationById(ctx context.Context, id int64) (tables_model.DcFollowingCalculation, error) {
	row, err := s.TablesRepository.GetFollowingCalculationById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcFollowingCalculation{}, fmt.Errorf("dc.following_calculation id = %d: %w", id, customerrors.ErrNotFound)
		}

		return tables_model.DcFollowingCalculation{}, fmt.Errorf("get dc.following_calculation id = %d: %w", id, err)
	}

	return row, nil
}

// GetFollowingCalculations возвращает все активные строки dc.following_calculation.
func (s *TablesService) GetFollowingCalculations(ctx context.Context) ([]tables_model.DcFollowingCalculation, error) {
	rows, err := s.TablesRepository.GetFollowingCalculations(ctx)

	if err != nil {
		return nil, fmt.Errorf("get dc.following_calculation: %w", err)
	}

	return rows, nil
}

// GetDeletedFollowingCalculationById возвращает мягко удалённую строку dc.following_calculation по id.
func (s *TablesService) GetDeletedFollowingCalculationById(ctx context.Context, id int64) (tables_model.DcFollowingCalculation, error) {
	row, err := s.TablesRepository.GetDeletedFollowingCalculationById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcFollowingCalculation{}, fmt.Errorf("deleted dc.following_calculation id = %d: %w", id, customerrors.ErrNotFound)
		}

		return tables_model.DcFollowingCalculation{}, fmt.Errorf("get deleted dc.following_calculation id = %d: %w", id, err)
	}

	return row, nil
}

// GetDeletedFollowingCalculations возвращает все мягко удалённые строки dc.following_calculation.
func (s *TablesService) GetDeletedFollowingCalculations(ctx context.Context) ([]tables_model.DcFollowingCalculation, error) {
	rows, err := s.TablesRepository.GetDeletedFollowingCalculations(ctx)

	if err != nil {
		return nil, fmt.Errorf("get deleted dc.following_calculation: %w", err)
	}

	return rows, nil
}

// CreateFollowingCalculation вставляет строку dc.following_calculation и возвращает её целиком.
func (s *TablesService) CreateFollowingCalculation(ctx context.Context, params tables_model.CreateFollowingCalculationParams) (tables_model.DcFollowingCalculation, error) {
	row, err := s.TablesRepository.CreateFollowingCalculation(ctx, params)

	if err != nil {
		return tables_model.DcFollowingCalculation{}, fmt.Errorf("%w: dc.following_calculation: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateFollowingCalculationById обновляет активную строку dc.following_calculation и возвращает её целиком.
//
// Запрос фильтрует по is_deleted = false, поэтому попытка обновить удалённую
// или несуществующую запись даёт sql.ErrNoRows — переводим его в ErrNotFound,
// чтобы api-слой ответил NotFound, а не Internal.
func (s *TablesService) UpdateFollowingCalculationById(ctx context.Context, params tables_model.UpdateFollowingCalculationByIdParams) (tables_model.DcFollowingCalculation, error) {
	row, err := s.TablesRepository.UpdateFollowingCalculationById(ctx, params)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcFollowingCalculation{}, fmt.Errorf("dc.following_calculation id = %d: %w", params.ID, customerrors.ErrNotFound)
		}

		return tables_model.DcFollowingCalculation{}, fmt.Errorf("%w: dc.following_calculation id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return row, nil
}

// DeleteFollowingCalculationById мягко удаляет строку dc.following_calculation.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *TablesService) DeleteFollowingCalculationById(ctx context.Context, id int64) error {
	if _, err := s.GetFollowingCalculationById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.TablesRepository.DeleteFollowingCalculationById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.following_calculation id = %d: %w", customerrors.ErrDelete, id, err)
	}

	return nil
}

// UndeleteFollowingCalculationById восстанавливает мягко удалённую строку dc.following_calculation.
// Существование удалённой записи проверяется заранее по той же причине,
// что и в DeleteFollowingCalculationById.
func (s *TablesService) UndeleteFollowingCalculationById(ctx context.Context, id int64) error {
	if _, err := s.GetDeletedFollowingCalculationById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrUndelete, err)
	}

	if err := s.TablesRepository.UndeleteFollowingCalculationById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.following_calculation id = %d: %w", customerrors.ErrUndelete, id, err)
	}

	return nil
}
