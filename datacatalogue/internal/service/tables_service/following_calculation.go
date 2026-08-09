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

// GetFollowingCalculations возвращает страницу строк dc.following_calculation и её счётчики.
func (s *TablesService) GetFollowingCalculations(ctx context.Context, params tables_model.GetFollowingCalculationsParams) ([]tables_model.DcFollowingCalculation, tables_model.CountGetFollowingCalculationsRow, error) {
	count, err := s.TablesRepository.CountGetFollowingCalculations(ctx, params.PageLimit)
	if err != nil {
		return nil, tables_model.CountGetFollowingCalculationsRow{}, fmt.Errorf("count dc.following_calculation: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcFollowingCalculation{}, count, nil
	}

	rows, err := s.TablesRepository.GetFollowingCalculations(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetFollowingCalculationsRow{}, fmt.Errorf("get dc.following_calculation page: %w", err)
	}

	return rows, count, nil
}

// GetFollowingCalculationsDeleted возвращает страницу строк dc.following_calculation и её счётчики.
func (s *TablesService) GetFollowingCalculationsDeleted(ctx context.Context, params tables_model.GetFollowingCalculationsDeletedParams) ([]tables_model.DcFollowingCalculation, tables_model.CountGetFollowingCalculationsDeletedRow, error) {
	count, err := s.TablesRepository.CountGetFollowingCalculationsDeleted(ctx, params.PageLimit)
	if err != nil {
		return nil, tables_model.CountGetFollowingCalculationsDeletedRow{}, fmt.Errorf("count dc.following_calculation: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcFollowingCalculation{}, count, nil
	}

	rows, err := s.TablesRepository.GetFollowingCalculationsDeleted(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetFollowingCalculationsDeletedRow{}, fmt.Errorf("get dc.following_calculation page: %w", err)
	}

	return rows, count, nil
}

// GetFollowingCalculationsFilters возвращает страницу строк dc.following_calculation и её счётчики.
func (s *TablesService) GetFollowingCalculationsFilters(ctx context.Context, params tables_model.GetFollowingCalculationsFiltersParams) ([]tables_model.DcFollowingCalculation, tables_model.CountGetFollowingCalculationsFiltersRow, error) {
	countParams := tables_model.CountGetFollowingCalculationsFiltersParams{
		PageLimit:         params.PageLimit,
		ColumnCatID:       params.ColumnCatID,
		CalculationTypeID: params.CalculationTypeID,
	}
	count, err := s.TablesRepository.CountGetFollowingCalculationsFilters(ctx, countParams)
	if err != nil {
		return nil, tables_model.CountGetFollowingCalculationsFiltersRow{}, fmt.Errorf("count dc.following_calculation: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcFollowingCalculation{}, count, nil
	}

	rows, err := s.TablesRepository.GetFollowingCalculationsFilters(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetFollowingCalculationsFiltersRow{}, fmt.Errorf("get dc.following_calculation page: %w", err)
	}

	return rows, count, nil
}

// CreateFollowingCalculation вставляет строку dc.following_calculation и возвращает её целиком.
func (s *TablesService) CreateFollowingCalculation(ctx context.Context, params tables_model.CreateFollowingCalculationParams) (tables_model.DcFollowingCalculation, error) {
	row, err := s.TablesRepository.CreateFollowingCalculation(ctx, params)

	if err != nil {
		return tables_model.DcFollowingCalculation{}, fmt.Errorf("%w: dc.following_calculation: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateFollowingCalculationById обновляет строку dc.following_calculation.
//
// Запрос — :exec и не сообщает число затронутых строк, поэтому
// существование активной записи проверяется заранее.
func (s *TablesService) UpdateFollowingCalculationById(ctx context.Context, params tables_model.UpdateFollowingCalculationByIdParams) error {
	if _, err := s.GetFollowingCalculationById(ctx, params.ID); err != nil {
		return fmt.Errorf("%w: dc.following_calculation: %w", customerrors.ErrUpdate, err)
	}

	if err := s.TablesRepository.UpdateFollowingCalculationById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.following_calculation id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return nil
}

// DeleteFollowingCalculationById мягко удаляет строку dc.following_calculation.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *TablesService) DeleteFollowingCalculationById(ctx context.Context, params tables_model.DeleteFollowingCalculationByIdParams) error {
	if _, err := s.GetFollowingCalculationById(ctx, params.ID); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.TablesRepository.DeleteFollowingCalculationById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.following_calculation id = %d: %w", customerrors.ErrDelete, params.ID, err)
	}

	return nil
}

// UndeleteFollowingCalculationById восстанавливает мягко удалённую строку dc.following_calculation.
func (s *TablesService) UndeleteFollowingCalculationById(ctx context.Context, params tables_model.UndeleteFollowingCalculationByIdParams) error {
	if err := s.TablesRepository.UndeleteFollowingCalculationById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.following_calculation id = %d: %w", customerrors.ErrUndelete, params.ID, err)
	}

	return nil
}
