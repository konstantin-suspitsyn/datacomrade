package tablesservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	customerrors "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/utils/custom_errors"
)

// GetDatabaseCatById возвращает активную строку dc.database_cat по id.
func (s *TablesService) GetDatabaseCatById(ctx context.Context, id int64) (tables_model.DcDatabaseCat, error) {
	row, err := s.TablesRepository.GetDatabaseCatById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcDatabaseCat{}, fmt.Errorf("dc.database_cat id = %d: %w", id, customerrors.ErrNotFound)
		}

		return tables_model.DcDatabaseCat{}, fmt.Errorf("get dc.database_cat id = %d: %w", id, err)
	}

	return row, nil
}

// GetDatabaseCats возвращает все активные строки dc.database_cat.
func (s *TablesService) GetDatabaseCats(ctx context.Context) ([]tables_model.DcDatabaseCat, error) {
	rows, err := s.TablesRepository.GetDatabaseCats(ctx)

	if err != nil {
		return nil, fmt.Errorf("get dc.database_cat: %w", err)
	}

	return rows, nil
}

// GetDeletedDatabaseCatById возвращает мягко удалённую строку dc.database_cat по id.
func (s *TablesService) GetDeletedDatabaseCatById(ctx context.Context, id int64) (tables_model.DcDatabaseCat, error) {
	row, err := s.TablesRepository.GetDeletedDatabaseCatById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcDatabaseCat{}, fmt.Errorf("deleted dc.database_cat id = %d: %w", id, customerrors.ErrNotFound)
		}

		return tables_model.DcDatabaseCat{}, fmt.Errorf("get deleted dc.database_cat id = %d: %w", id, err)
	}

	return row, nil
}

// GetDeletedDatabaseCats возвращает все мягко удалённые строки dc.database_cat.
func (s *TablesService) GetDeletedDatabaseCats(ctx context.Context) ([]tables_model.DcDatabaseCat, error) {
	rows, err := s.TablesRepository.GetDeletedDatabaseCats(ctx)

	if err != nil {
		return nil, fmt.Errorf("get deleted dc.database_cat: %w", err)
	}

	return rows, nil
}

// GetDatabaseCatsByHostId возвращает активные строки dc.database_cat, отобранные по host_id.
func (s *TablesService) GetDatabaseCatsByHostId(ctx context.Context, hostID int64) ([]tables_model.DcDatabaseCat, error) {
	rows, err := s.TablesRepository.GetDatabaseCatsByHostId(ctx, hostID)

	if err != nil {
		return nil, fmt.Errorf("get dc.database_cat by host_id = %d: %w", hostID, err)
	}

	return rows, nil
}

// GetDatabaseCatsByDatabaseTypeId возвращает активные строки dc.database_cat, отобранные по database_type_id.
func (s *TablesService) GetDatabaseCatsByDatabaseTypeId(ctx context.Context, databaseTypeID int64) ([]tables_model.DcDatabaseCat, error) {
	rows, err := s.TablesRepository.GetDatabaseCatsByDatabaseTypeId(ctx, databaseTypeID)

	if err != nil {
		return nil, fmt.Errorf("get dc.database_cat by database_type_id = %d: %w", databaseTypeID, err)
	}

	return rows, nil
}

// CreateDatabaseCat вставляет строку dc.database_cat и возвращает её целиком.
func (s *TablesService) CreateDatabaseCat(ctx context.Context, params tables_model.CreateDatabaseCatParams) (tables_model.DcDatabaseCat, error) {
	row, err := s.TablesRepository.CreateDatabaseCat(ctx, params)

	if err != nil {
		return tables_model.DcDatabaseCat{}, fmt.Errorf("%w: dc.database_cat: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateDatabaseCatById обновляет активную строку dc.database_cat и возвращает её целиком.
//
// Запрос фильтрует по is_deleted = false, поэтому попытка обновить удалённую
// или несуществующую запись даёт sql.ErrNoRows — переводим его в ErrNotFound,
// чтобы api-слой ответил NotFound, а не Internal.
func (s *TablesService) UpdateDatabaseCatById(ctx context.Context, params tables_model.UpdateDatabaseCatByIdParams) (tables_model.DcDatabaseCat, error) {
	row, err := s.TablesRepository.UpdateDatabaseCatById(ctx, params)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcDatabaseCat{}, fmt.Errorf("dc.database_cat id = %d: %w", params.ID, customerrors.ErrNotFound)
		}

		return tables_model.DcDatabaseCat{}, fmt.Errorf("%w: dc.database_cat id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return row, nil
}

// DeleteDatabaseCatById мягко удаляет строку dc.database_cat.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *TablesService) DeleteDatabaseCatById(ctx context.Context, id int64) error {
	if _, err := s.GetDatabaseCatById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.TablesRepository.DeleteDatabaseCatById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.database_cat id = %d: %w", customerrors.ErrDelete, id, err)
	}

	return nil
}

// UndeleteDatabaseCatById восстанавливает мягко удалённую строку dc.database_cat.
// Существование удалённой записи проверяется заранее по той же причине,
// что и в DeleteDatabaseCatById.
func (s *TablesService) UndeleteDatabaseCatById(ctx context.Context, id int64) error {
	if _, err := s.GetDeletedDatabaseCatById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrUndelete, err)
	}

	if err := s.TablesRepository.UndeleteDatabaseCatById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.database_cat id = %d: %w", customerrors.ErrUndelete, id, err)
	}

	return nil
}
