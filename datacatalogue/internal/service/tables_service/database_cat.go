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

// GetDatabaseCats возвращает страницу строк dc.database_cat и её счётчики.
func (s *TablesService) GetDatabaseCats(ctx context.Context, params tables_model.GetDatabaseCatsParams) ([]tables_model.DcDatabaseCat, tables_model.CountGetDatabaseCatsRow, error) {
	count, err := s.TablesRepository.CountGetDatabaseCats(ctx, params.PageLimit)
	if err != nil {
		return nil, tables_model.CountGetDatabaseCatsRow{}, fmt.Errorf("count dc.database_cat: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcDatabaseCat{}, count, nil
	}

	rows, err := s.TablesRepository.GetDatabaseCats(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetDatabaseCatsRow{}, fmt.Errorf("get dc.database_cat page: %w", err)
	}

	return rows, count, nil
}

// GetDatabaseCatsDeleted возвращает страницу строк dc.database_cat и её счётчики.
func (s *TablesService) GetDatabaseCatsDeleted(ctx context.Context, params tables_model.GetDatabaseCatsDeletedParams) ([]tables_model.DcDatabaseCat, tables_model.CountGetDatabaseCatsDeletedRow, error) {
	count, err := s.TablesRepository.CountGetDatabaseCatsDeleted(ctx, params.PageLimit)
	if err != nil {
		return nil, tables_model.CountGetDatabaseCatsDeletedRow{}, fmt.Errorf("count dc.database_cat: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcDatabaseCat{}, count, nil
	}

	rows, err := s.TablesRepository.GetDatabaseCatsDeleted(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetDatabaseCatsDeletedRow{}, fmt.Errorf("get dc.database_cat page: %w", err)
	}

	return rows, count, nil
}

// GetDatabaseCatsByHostIdDatabaseTypeId возвращает страницу строк dc.database_cat и её счётчики.
func (s *TablesService) GetDatabaseCatsByHostIdDatabaseTypeId(ctx context.Context, params tables_model.GetDatabaseCatsByHostIdDatabaseTypeIdParams) ([]tables_model.DcDatabaseCat, tables_model.CountGetDatabaseCatsByHostIdDatabaseTypeIdRow, error) {
	countParams := tables_model.CountGetDatabaseCatsByHostIdDatabaseTypeIdParams{
		PageLimit:      params.PageLimit,
		HostID:         params.HostID,
		DatabaseTypeID: params.DatabaseTypeID,
	}
	count, err := s.TablesRepository.CountGetDatabaseCatsByHostIdDatabaseTypeId(ctx, countParams)
	if err != nil {
		return nil, tables_model.CountGetDatabaseCatsByHostIdDatabaseTypeIdRow{}, fmt.Errorf("count dc.database_cat: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcDatabaseCat{}, count, nil
	}

	rows, err := s.TablesRepository.GetDatabaseCatsByHostIdDatabaseTypeId(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetDatabaseCatsByHostIdDatabaseTypeIdRow{}, fmt.Errorf("get dc.database_cat page: %w", err)
	}

	return rows, count, nil
}

// GetDatabaseCatsSearchName возвращает страницу строк dc.database_cat и её счётчики.
func (s *TablesService) GetDatabaseCatsSearchName(ctx context.Context, params tables_model.GetDatabaseCatsSearchNameParams) ([]tables_model.DcDatabaseCat, tables_model.CountGetDatabaseCatsSearchNameRow, error) {
	countParams := tables_model.CountGetDatabaseCatsSearchNameParams{
		PageLimit:  params.PageLimit,
		SearchName: params.SearchName,
	}
	count, err := s.TablesRepository.CountGetDatabaseCatsSearchName(ctx, countParams)
	if err != nil {
		return nil, tables_model.CountGetDatabaseCatsSearchNameRow{}, fmt.Errorf("count dc.database_cat: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcDatabaseCat{}, count, nil
	}

	rows, err := s.TablesRepository.GetDatabaseCatsSearchName(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetDatabaseCatsSearchNameRow{}, fmt.Errorf("get dc.database_cat page: %w", err)
	}

	return rows, count, nil
}

// CreateDatabaseCat вставляет строку dc.database_cat и возвращает её целиком.
func (s *TablesService) CreateDatabaseCat(ctx context.Context, params tables_model.CreateDatabaseCatParams) (tables_model.DcDatabaseCat, error) {
	row, err := s.TablesRepository.CreateDatabaseCat(ctx, params)

	if err != nil {
		return tables_model.DcDatabaseCat{}, fmt.Errorf("%w: dc.database_cat: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateDatabaseCatById обновляет строку dc.database_cat.
//
// Запрос — :exec и не сообщает число затронутых строк, поэтому
// существование активной записи проверяется заранее.
func (s *TablesService) UpdateDatabaseCatById(ctx context.Context, params tables_model.UpdateDatabaseCatByIdParams) error {
	if _, err := s.GetDatabaseCatById(ctx, params.ID); err != nil {
		return fmt.Errorf("%w: dc.database_cat: %w", customerrors.ErrUpdate, err)
	}

	if err := s.TablesRepository.UpdateDatabaseCatById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.database_cat id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return nil
}

// DeleteDatabaseCatById мягко удаляет строку dc.database_cat.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *TablesService) DeleteDatabaseCatById(ctx context.Context, params tables_model.DeleteDatabaseCatByIdParams) error {
	if _, err := s.GetDatabaseCatById(ctx, params.ID); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.TablesRepository.DeleteDatabaseCatById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.database_cat id = %d: %w", customerrors.ErrDelete, params.ID, err)
	}

	return nil
}

// UndeleteDatabaseCatById восстанавливает мягко удалённую строку dc.database_cat.
func (s *TablesService) UndeleteDatabaseCatById(ctx context.Context, params tables_model.UndeleteDatabaseCatByIdParams) error {
	if err := s.TablesRepository.UndeleteDatabaseCatById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.database_cat id = %d: %w", customerrors.ErrUndelete, params.ID, err)
	}

	return nil
}
