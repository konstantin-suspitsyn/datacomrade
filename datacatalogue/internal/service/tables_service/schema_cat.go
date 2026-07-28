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

// GetSchemaCats возвращает все активные строки dc.schema_cat.
func (s *TablesService) GetSchemaCats(ctx context.Context) ([]tables_model.DcSchemaCat, error) {
	rows, err := s.TablesRepository.GetSchemaCats(ctx)

	if err != nil {
		return nil, fmt.Errorf("get dc.schema_cat: %w", err)
	}

	return rows, nil
}

// GetDeletedSchemaCatById возвращает мягко удалённую строку dc.schema_cat по id.
func (s *TablesService) GetDeletedSchemaCatById(ctx context.Context, id int64) (tables_model.DcSchemaCat, error) {
	row, err := s.TablesRepository.GetDeletedSchemaCatById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcSchemaCat{}, fmt.Errorf("deleted dc.schema_cat id = %d: %w", id, customerrors.ErrNotFound)
		}

		return tables_model.DcSchemaCat{}, fmt.Errorf("get deleted dc.schema_cat id = %d: %w", id, err)
	}

	return row, nil
}

// GetDeletedSchemaCats возвращает все мягко удалённые строки dc.schema_cat.
func (s *TablesService) GetDeletedSchemaCats(ctx context.Context) ([]tables_model.DcSchemaCat, error) {
	rows, err := s.TablesRepository.GetDeletedSchemaCats(ctx)

	if err != nil {
		return nil, fmt.Errorf("get deleted dc.schema_cat: %w", err)
	}

	return rows, nil
}

// GetSchemaCatsByDatabaseId возвращает активные строки dc.schema_cat, отобранные по database_id.
func (s *TablesService) GetSchemaCatsByDatabaseId(ctx context.Context, databaseID int64) ([]tables_model.DcSchemaCat, error) {
	rows, err := s.TablesRepository.GetSchemaCatsByDatabaseId(ctx, databaseID)

	if err != nil {
		return nil, fmt.Errorf("get dc.schema_cat by database_id = %d: %w", databaseID, err)
	}

	return rows, nil
}

// CreateSchemaCat вставляет строку dc.schema_cat и возвращает её целиком.
func (s *TablesService) CreateSchemaCat(ctx context.Context, params tables_model.CreateSchemaCatParams) (tables_model.DcSchemaCat, error) {
	row, err := s.TablesRepository.CreateSchemaCat(ctx, params)

	if err != nil {
		return tables_model.DcSchemaCat{}, fmt.Errorf("%w: dc.schema_cat: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateSchemaCatById обновляет активную строку dc.schema_cat и возвращает её целиком.
//
// Запрос фильтрует по is_deleted = false, поэтому попытка обновить удалённую
// или несуществующую запись даёт sql.ErrNoRows — переводим его в ErrNotFound,
// чтобы api-слой ответил NotFound, а не Internal.
func (s *TablesService) UpdateSchemaCatById(ctx context.Context, params tables_model.UpdateSchemaCatByIdParams) (tables_model.DcSchemaCat, error) {
	row, err := s.TablesRepository.UpdateSchemaCatById(ctx, params)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcSchemaCat{}, fmt.Errorf("dc.schema_cat id = %d: %w", params.ID, customerrors.ErrNotFound)
		}

		return tables_model.DcSchemaCat{}, fmt.Errorf("%w: dc.schema_cat id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return row, nil
}

// DeleteSchemaCatById мягко удаляет строку dc.schema_cat.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *TablesService) DeleteSchemaCatById(ctx context.Context, id int64) error {
	if _, err := s.GetSchemaCatById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.TablesRepository.DeleteSchemaCatById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.schema_cat id = %d: %w", customerrors.ErrDelete, id, err)
	}

	return nil
}

// UndeleteSchemaCatById восстанавливает мягко удалённую строку dc.schema_cat.
// Существование удалённой записи проверяется заранее по той же причине,
// что и в DeleteSchemaCatById.
func (s *TablesService) UndeleteSchemaCatById(ctx context.Context, id int64) error {
	if _, err := s.GetDeletedSchemaCatById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrUndelete, err)
	}

	if err := s.TablesRepository.UndeleteSchemaCatById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.schema_cat id = %d: %w", customerrors.ErrUndelete, id, err)
	}

	return nil
}
