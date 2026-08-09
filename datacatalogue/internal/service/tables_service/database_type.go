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

// GetDatabaseTypesSearchName возвращает страницу строк dc.database_type и её счётчики.
func (s *TablesService) GetDatabaseTypesSearchName(ctx context.Context, params tables_model.GetDatabaseTypesSearchNameParams) ([]tables_model.DcDatabaseType, tables_model.CountGetDatabaseTypesSearchNameRow, error) {
	countParams := tables_model.CountGetDatabaseTypesSearchNameParams{
		PageLimit:  params.PageLimit,
		SearchName: params.SearchName,
	}
	count, err := s.TablesRepository.CountGetDatabaseTypesSearchName(ctx, countParams)
	if err != nil {
		return nil, tables_model.CountGetDatabaseTypesSearchNameRow{}, fmt.Errorf("count dc.database_type: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcDatabaseType{}, count, nil
	}

	rows, err := s.TablesRepository.GetDatabaseTypesSearchName(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetDatabaseTypesSearchNameRow{}, fmt.Errorf("get dc.database_type page: %w", err)
	}

	return rows, count, nil
}

// GetDatabaseTypes возвращает страницу строк dc.database_type и её счётчики.
func (s *TablesService) GetDatabaseTypes(ctx context.Context, params tables_model.GetDatabaseTypesParams) ([]tables_model.DcDatabaseType, tables_model.CountGetDatabaseTypesRow, error) {
	count, err := s.TablesRepository.CountGetDatabaseTypes(ctx, params.PageLimit)
	if err != nil {
		return nil, tables_model.CountGetDatabaseTypesRow{}, fmt.Errorf("count dc.database_type: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcDatabaseType{}, count, nil
	}

	rows, err := s.TablesRepository.GetDatabaseTypes(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetDatabaseTypesRow{}, fmt.Errorf("get dc.database_type page: %w", err)
	}

	return rows, count, nil
}

// GetDatabaseTypesDeleted возвращает страницу строк dc.database_type и её счётчики.
func (s *TablesService) GetDatabaseTypesDeleted(ctx context.Context, params tables_model.GetDatabaseTypesDeletedParams) ([]tables_model.DcDatabaseType, tables_model.CountGetDatabaseTypesDeletedRow, error) {
	count, err := s.TablesRepository.CountGetDatabaseTypesDeleted(ctx, params.PageLimit)
	if err != nil {
		return nil, tables_model.CountGetDatabaseTypesDeletedRow{}, fmt.Errorf("count dc.database_type: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcDatabaseType{}, count, nil
	}

	rows, err := s.TablesRepository.GetDatabaseTypesDeleted(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetDatabaseTypesDeletedRow{}, fmt.Errorf("get dc.database_type page: %w", err)
	}

	return rows, count, nil
}

// CreateDatabaseType вставляет строку dc.database_type и возвращает её целиком.
func (s *TablesService) CreateDatabaseType(ctx context.Context, params tables_model.CreateDatabaseTypeParams) (tables_model.DcDatabaseType, error) {
	row, err := s.TablesRepository.CreateDatabaseType(ctx, params)

	if err != nil {
		return tables_model.DcDatabaseType{}, fmt.Errorf("%w: dc.database_type: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateDatabaseTypeById обновляет строку dc.database_type.
//
// Запрос — :exec и не сообщает число затронутых строк, поэтому
// существование активной записи проверяется заранее.
func (s *TablesService) UpdateDatabaseTypeById(ctx context.Context, params tables_model.UpdateDatabaseTypeByIdParams) error {
	if _, err := s.GetDatabaseTypeById(ctx, params.ID); err != nil {
		return fmt.Errorf("%w: dc.database_type: %w", customerrors.ErrUpdate, err)
	}

	if err := s.TablesRepository.UpdateDatabaseTypeById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.database_type id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return nil
}

// DeleteDatabaseTypeById мягко удаляет строку dc.database_type.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *TablesService) DeleteDatabaseTypeById(ctx context.Context, params tables_model.DeleteDatabaseTypeByIdParams) error {
	if _, err := s.GetDatabaseTypeById(ctx, params.ID); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.TablesRepository.DeleteDatabaseTypeById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.database_type id = %d: %w", customerrors.ErrDelete, params.ID, err)
	}

	return nil
}

// UndeleteDatabaseTypeById восстанавливает мягко удалённую строку dc.database_type.
func (s *TablesService) UndeleteDatabaseTypeById(ctx context.Context, params tables_model.UndeleteDatabaseTypeByIdParams) error {
	if err := s.TablesRepository.UndeleteDatabaseTypeById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.database_type id = %d: %w", customerrors.ErrUndelete, params.ID, err)
	}

	return nil
}
