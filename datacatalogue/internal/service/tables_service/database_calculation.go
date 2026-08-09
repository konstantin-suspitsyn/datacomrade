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

// GetDatabaseCalculationsSearchName возвращает страницу строк dc.database_calculation и её счётчики.
func (s *TablesService) GetDatabaseCalculationsSearchName(ctx context.Context, params tables_model.GetDatabaseCalculationsSearchNameParams) ([]tables_model.DcDatabaseCalculation, tables_model.CountGetDatabaseCalculationsSearchNameRow, error) {
	countParams := tables_model.CountGetDatabaseCalculationsSearchNameParams{
		PageLimit:  params.PageLimit,
		SearchName: params.SearchName,
	}
	count, err := s.TablesRepository.CountGetDatabaseCalculationsSearchName(ctx, countParams)
	if err != nil {
		return nil, tables_model.CountGetDatabaseCalculationsSearchNameRow{}, fmt.Errorf("count dc.database_calculation: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcDatabaseCalculation{}, count, nil
	}

	rows, err := s.TablesRepository.GetDatabaseCalculationsSearchName(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetDatabaseCalculationsSearchNameRow{}, fmt.Errorf("get dc.database_calculation page: %w", err)
	}

	return rows, count, nil
}

// GetDatabaseCalculationsFiltered возвращает страницу строк dc.database_calculation и её счётчики.
func (s *TablesService) GetDatabaseCalculationsFiltered(ctx context.Context, params tables_model.GetDatabaseCalculationsFilteredParams) ([]tables_model.DcDatabaseCalculation, tables_model.CountGetDatabaseCalculationsFilteredRow, error) {
	countParams := tables_model.CountGetDatabaseCalculationsFilteredParams{
		PageLimit:         params.PageLimit,
		DatabaseCatID:     params.DatabaseCatID,
		CalculationTypeID: params.CalculationTypeID,
	}
	count, err := s.TablesRepository.CountGetDatabaseCalculationsFiltered(ctx, countParams)
	if err != nil {
		return nil, tables_model.CountGetDatabaseCalculationsFilteredRow{}, fmt.Errorf("count dc.database_calculation: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcDatabaseCalculation{}, count, nil
	}

	rows, err := s.TablesRepository.GetDatabaseCalculationsFiltered(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetDatabaseCalculationsFilteredRow{}, fmt.Errorf("get dc.database_calculation page: %w", err)
	}

	return rows, count, nil
}

// GetDatabaseCalculationsDeleted возвращает страницу строк dc.database_calculation и её счётчики.
func (s *TablesService) GetDatabaseCalculationsDeleted(ctx context.Context, params tables_model.GetDatabaseCalculationsDeletedParams) ([]tables_model.DcDatabaseCalculation, tables_model.CountGetDatabaseCalculationsDeletedRow, error) {
	count, err := s.TablesRepository.CountGetDatabaseCalculationsDeleted(ctx, params.PageLimit)
	if err != nil {
		return nil, tables_model.CountGetDatabaseCalculationsDeletedRow{}, fmt.Errorf("count dc.database_calculation: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcDatabaseCalculation{}, count, nil
	}

	rows, err := s.TablesRepository.GetDatabaseCalculationsDeleted(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetDatabaseCalculationsDeletedRow{}, fmt.Errorf("get dc.database_calculation page: %w", err)
	}

	return rows, count, nil
}

// GetDatabaseCalculations возвращает страницу строк dc.database_calculation и её счётчики.
func (s *TablesService) GetDatabaseCalculations(ctx context.Context, params tables_model.GetDatabaseCalculationsParams) ([]tables_model.DcDatabaseCalculation, tables_model.CountGetDatabaseCalculationsRow, error) {
	count, err := s.TablesRepository.CountGetDatabaseCalculations(ctx, params.PageLimit)
	if err != nil {
		return nil, tables_model.CountGetDatabaseCalculationsRow{}, fmt.Errorf("count dc.database_calculation: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcDatabaseCalculation{}, count, nil
	}

	rows, err := s.TablesRepository.GetDatabaseCalculations(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetDatabaseCalculationsRow{}, fmt.Errorf("get dc.database_calculation page: %w", err)
	}

	return rows, count, nil
}

// CreateDatabaseCalculation вставляет строку dc.database_calculation и возвращает её целиком.
func (s *TablesService) CreateDatabaseCalculation(ctx context.Context, params tables_model.CreateDatabaseCalculationParams) (tables_model.DcDatabaseCalculation, error) {
	row, err := s.TablesRepository.CreateDatabaseCalculation(ctx, params)

	if err != nil {
		return tables_model.DcDatabaseCalculation{}, fmt.Errorf("%w: dc.database_calculation: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateDatabaseCalculationById обновляет строку dc.database_calculation.
//
// Запрос — :exec и не сообщает число затронутых строк, поэтому
// существование активной записи проверяется заранее.
func (s *TablesService) UpdateDatabaseCalculationById(ctx context.Context, params tables_model.UpdateDatabaseCalculationByIdParams) error {
	if _, err := s.GetDatabaseCalculationById(ctx, params.ID); err != nil {
		return fmt.Errorf("%w: dc.database_calculation: %w", customerrors.ErrUpdate, err)
	}

	if err := s.TablesRepository.UpdateDatabaseCalculationById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.database_calculation id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return nil
}

// DeleteDatabaseCalculationById мягко удаляет строку dc.database_calculation.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *TablesService) DeleteDatabaseCalculationById(ctx context.Context, params tables_model.DeleteDatabaseCalculationByIdParams) error {
	if _, err := s.GetDatabaseCalculationById(ctx, params.ID); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.TablesRepository.DeleteDatabaseCalculationById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.database_calculation id = %d: %w", customerrors.ErrDelete, params.ID, err)
	}

	return nil
}

// UndeleteDatabaseCalculationById восстанавливает мягко удалённую строку dc.database_calculation.
func (s *TablesService) UndeleteDatabaseCalculationById(ctx context.Context, params tables_model.UndeleteDatabaseCalculationByIdParams) error {
	if err := s.TablesRepository.UndeleteDatabaseCalculationById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.database_calculation id = %d: %w", customerrors.ErrUndelete, params.ID, err)
	}

	return nil
}
