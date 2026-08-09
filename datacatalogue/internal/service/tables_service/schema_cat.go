package tablesservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	customerrors "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/utils/custom_errors"
)

// GetSchemaCatById возвращает активную строку dc.schema_cat по id.
func (s *TablesService) GetSchemaCatById(ctx context.Context, id int64) (tables_model.DcSchemaCat, error) {
	row, err := s.TablesRepository.GetSchemaCatById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcSchemaCat{}, fmt.Errorf("dc.schema_cat id = %d: %w", id, customerrors.ErrNotFound)
		}

		return tables_model.DcSchemaCat{}, fmt.Errorf("get dc.schema_cat id = %d: %w", id, err)
	}

	return row, nil
}

// GetSchemaCats возвращает страницу строк dc.schema_cat и её счётчики.
func (s *TablesService) GetSchemaCats(ctx context.Context, params tables_model.GetSchemaCatsParams) ([]tables_model.DcSchemaCat, tables_model.CountGetSchemaCatsRow, error) {
	count, err := s.TablesRepository.CountGetSchemaCats(ctx, params.PageLimit)
	if err != nil {
		return nil, tables_model.CountGetSchemaCatsRow{}, fmt.Errorf("count dc.schema_cat: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcSchemaCat{}, count, nil
	}

	rows, err := s.TablesRepository.GetSchemaCats(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetSchemaCatsRow{}, fmt.Errorf("get dc.schema_cat page: %w", err)
	}

	return rows, count, nil
}

// GetSchemaCatsDeleted возвращает страницу строк dc.schema_cat и её счётчики.
func (s *TablesService) GetSchemaCatsDeleted(ctx context.Context, params tables_model.GetSchemaCatsDeletedParams) ([]tables_model.DcSchemaCat, tables_model.CountGetSchemaCatsDeletedRow, error) {
	count, err := s.TablesRepository.CountGetSchemaCatsDeleted(ctx, params.PageLimit)
	if err != nil {
		return nil, tables_model.CountGetSchemaCatsDeletedRow{}, fmt.Errorf("count dc.schema_cat: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcSchemaCat{}, count, nil
	}

	rows, err := s.TablesRepository.GetSchemaCatsDeleted(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetSchemaCatsDeletedRow{}, fmt.Errorf("get dc.schema_cat page: %w", err)
	}

	return rows, count, nil
}

// CreateSchemaCat вставляет строку dc.schema_cat и возвращает её целиком.
func (s *TablesService) CreateSchemaCat(ctx context.Context, params tables_model.CreateSchemaCatParams) (tables_model.DcSchemaCat, error) {
	row, err := s.TablesRepository.CreateSchemaCat(ctx, params)

	if err != nil {
		return tables_model.DcSchemaCat{}, fmt.Errorf("%w: dc.schema_cat: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateSchemaCatById обновляет строку dc.schema_cat.
//
// Запрос — :exec и не сообщает число затронутых строк, поэтому
// существование активной записи проверяется заранее.
func (s *TablesService) UpdateSchemaCatById(ctx context.Context, params tables_model.UpdateSchemaCatByIdParams) error {
	if _, err := s.GetSchemaCatById(ctx, params.ID); err != nil {
		return fmt.Errorf("%w: dc.schema_cat: %w", customerrors.ErrUpdate, err)
	}

	if err := s.TablesRepository.UpdateSchemaCatById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.schema_cat id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return nil
}

// DeleteSchemaCatById мягко удаляет строку dc.schema_cat.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *TablesService) DeleteSchemaCatById(ctx context.Context, params tables_model.DeleteSchemaCatByIdParams) error {
	if _, err := s.GetSchemaCatById(ctx, params.ID); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.TablesRepository.DeleteSchemaCatById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.schema_cat id = %d: %w", customerrors.ErrDelete, params.ID, err)
	}

	return nil
}

// UndeleteSchemaCatById восстанавливает мягко удалённую строку dc.schema_cat.
func (s *TablesService) UndeleteSchemaCatById(ctx context.Context, params tables_model.UndeleteSchemaCatByIdParams) error {
	if err := s.TablesRepository.UndeleteSchemaCatById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.schema_cat id = %d: %w", customerrors.ErrUndelete, params.ID, err)
	}

	return nil
}
