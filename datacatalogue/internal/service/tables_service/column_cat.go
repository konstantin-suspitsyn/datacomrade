package tablesservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	customerrors "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/utils/custom_errors"
)

// GetColumnCatById возвращает активную строку dc.column_cat по id.
func (s *TablesService) GetColumnCatById(ctx context.Context, id int64) (tables_model.DcColumnCat, error) {
	row, err := s.TablesRepository.GetColumnCatById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcColumnCat{}, fmt.Errorf("dc.column_cat id = %d: %w", id, customerrors.ErrNotFound)
		}

		return tables_model.DcColumnCat{}, fmt.Errorf("get dc.column_cat id = %d: %w", id, err)
	}

	return row, nil
}

// GetColumnCats возвращает страницу строк dc.column_cat и её счётчики.
func (s *TablesService) GetColumnCats(ctx context.Context, params tables_model.GetColumnCatsParams) ([]tables_model.DcColumnCat, tables_model.CountGetColumnCatsRow, error) {
	count, err := s.TablesRepository.CountGetColumnCats(ctx, params.PageLimit)
	if err != nil {
		return nil, tables_model.CountGetColumnCatsRow{}, fmt.Errorf("count dc.column_cat: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcColumnCat{}, count, nil
	}

	rows, err := s.TablesRepository.GetColumnCats(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetColumnCatsRow{}, fmt.Errorf("get dc.column_cat page: %w", err)
	}

	return rows, count, nil
}

// GetColumnCatsDeleted возвращает страницу строк dc.column_cat и её счётчики.
func (s *TablesService) GetColumnCatsDeleted(ctx context.Context, params tables_model.GetColumnCatsDeletedParams) ([]tables_model.DcColumnCat, tables_model.CountGetColumnCatsDeletedRow, error) {
	count, err := s.TablesRepository.CountGetColumnCatsDeleted(ctx, params.PageLimit)
	if err != nil {
		return nil, tables_model.CountGetColumnCatsDeletedRow{}, fmt.Errorf("count dc.column_cat: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcColumnCat{}, count, nil
	}

	rows, err := s.TablesRepository.GetColumnCatsDeleted(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetColumnCatsDeletedRow{}, fmt.Errorf("get dc.column_cat page: %w", err)
	}

	return rows, count, nil
}

// GetColumnCatsOptionalFilters возвращает страницу строк dc.column_cat и её счётчики.
func (s *TablesService) GetColumnCatsOptionalFilters(ctx context.Context, params tables_model.GetColumnCatsOptionalFiltersParams) ([]tables_model.DcColumnCat, tables_model.CountGetColumnCatsOptionalFiltersRow, error) {
	countParams := tables_model.CountGetColumnCatsOptionalFiltersParams{
		PageLimit:         params.PageLimit,
		TableID:           params.TableID,
		AliasID:           params.AliasID,
		ColumnTypeID:      params.ColumnTypeID,
		CalculationTypeID: params.CalculationTypeID,
		ShowInUi:          params.ShowInUi,
	}
	count, err := s.TablesRepository.CountGetColumnCatsOptionalFilters(ctx, countParams)
	if err != nil {
		return nil, tables_model.CountGetColumnCatsOptionalFiltersRow{}, fmt.Errorf("count dc.column_cat: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcColumnCat{}, count, nil
	}

	rows, err := s.TablesRepository.GetColumnCatsOptionalFilters(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetColumnCatsOptionalFiltersRow{}, fmt.Errorf("get dc.column_cat page: %w", err)
	}

	return rows, count, nil
}

// GetColumnCatsSearchName возвращает страницу строк dc.column_cat и её счётчики.
func (s *TablesService) GetColumnCatsSearchName(ctx context.Context, params tables_model.GetColumnCatsSearchNameParams) ([]tables_model.DcColumnCat, tables_model.CountGetColumnCatsSearchNameRow, error) {
	countParams := tables_model.CountGetColumnCatsSearchNameParams{
		PageLimit:  params.PageLimit,
		SearchName: params.SearchName,
	}
	count, err := s.TablesRepository.CountGetColumnCatsSearchName(ctx, countParams)
	if err != nil {
		return nil, tables_model.CountGetColumnCatsSearchNameRow{}, fmt.Errorf("count dc.column_cat: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcColumnCat{}, count, nil
	}

	rows, err := s.TablesRepository.GetColumnCatsSearchName(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetColumnCatsSearchNameRow{}, fmt.Errorf("get dc.column_cat page: %w", err)
	}

	return rows, count, nil
}

// CreateColumnCat вставляет строку dc.column_cat и возвращает её целиком.
func (s *TablesService) CreateColumnCat(ctx context.Context, params tables_model.CreateColumnCatParams) (tables_model.DcColumnCat, error) {
	row, err := s.TablesRepository.CreateColumnCat(ctx, params)

	if err != nil {
		return tables_model.DcColumnCat{}, fmt.Errorf("%w: dc.column_cat: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateColumnCatById обновляет строку dc.column_cat.
//
// Запрос — :exec и не сообщает число затронутых строк, поэтому
// существование активной записи проверяется заранее.
func (s *TablesService) UpdateColumnCatById(ctx context.Context, params tables_model.UpdateColumnCatByIdParams) error {
	if _, err := s.GetColumnCatById(ctx, params.ID); err != nil {
		return fmt.Errorf("%w: dc.column_cat: %w", customerrors.ErrUpdate, err)
	}

	if err := s.TablesRepository.UpdateColumnCatById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.column_cat id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return nil
}

// DeleteColumnCatById мягко удаляет строку dc.column_cat.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *TablesService) DeleteColumnCatById(ctx context.Context, params tables_model.DeleteColumnCatByIdParams) error {
	if _, err := s.GetColumnCatById(ctx, params.ID); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.TablesRepository.DeleteColumnCatById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.column_cat id = %d: %w", customerrors.ErrDelete, params.ID, err)
	}

	return nil
}

// UndeleteColumnCatById восстанавливает мягко удалённую строку dc.column_cat.
func (s *TablesService) UndeleteColumnCatById(ctx context.Context, params tables_model.UndeleteColumnCatByIdParams) error {
	if err := s.TablesRepository.UndeleteColumnCatById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.column_cat id = %d: %w", customerrors.ErrUndelete, params.ID, err)
	}

	return nil
}
