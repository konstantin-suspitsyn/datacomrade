package tablesservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	customerrors "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/utils/custom_errors"
)

// GetTableCatById возвращает активную строку dc.table_cat по id.
func (s *TablesService) GetTableCatById(ctx context.Context, id int64) (tables_model.DcTableCat, error) {
	row, err := s.TablesRepository.GetTableCatById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcTableCat{}, fmt.Errorf("dc.table_cat id = %d: %w", id, customerrors.ErrNotFound)
		}

		return tables_model.DcTableCat{}, fmt.Errorf("get dc.table_cat id = %d: %w", id, err)
	}

	return row, nil
}

// GetTableCats возвращает страницу строк dc.table_cat и её счётчики.
func (s *TablesService) GetTableCats(ctx context.Context, params tables_model.GetTableCatsParams) ([]tables_model.DcTableCat, tables_model.CountGetTableCatsRow, error) {
	count, err := s.TablesRepository.CountGetTableCats(ctx, params.PageLimit)
	if err != nil {
		return nil, tables_model.CountGetTableCatsRow{}, fmt.Errorf("count dc.table_cat: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcTableCat{}, count, nil
	}

	rows, err := s.TablesRepository.GetTableCats(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetTableCatsRow{}, fmt.Errorf("get dc.table_cat page: %w", err)
	}

	return rows, count, nil
}

// GetTableCatsDeleted возвращает страницу строк dc.table_cat и её счётчики.
func (s *TablesService) GetTableCatsDeleted(ctx context.Context, params tables_model.GetTableCatsDeletedParams) ([]tables_model.DcTableCat, tables_model.CountGetTableCatsDeletedRow, error) {
	count, err := s.TablesRepository.CountGetTableCatsDeleted(ctx, params.PageLimit)
	if err != nil {
		return nil, tables_model.CountGetTableCatsDeletedRow{}, fmt.Errorf("count dc.table_cat: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcTableCat{}, count, nil
	}

	rows, err := s.TablesRepository.GetTableCatsDeleted(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetTableCatsDeletedRow{}, fmt.Errorf("get dc.table_cat page: %w", err)
	}

	return rows, count, nil
}

// CreateTableCat вставляет строку dc.table_cat и возвращает её целиком.
func (s *TablesService) CreateTableCat(ctx context.Context, params tables_model.CreateTableCatParams) (tables_model.DcTableCat, error) {
	row, err := s.TablesRepository.CreateTableCat(ctx, params)

	if err != nil {
		return tables_model.DcTableCat{}, fmt.Errorf("%w: dc.table_cat: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateTableCatById обновляет строку dc.table_cat.
//
// Запрос — :exec и не сообщает число затронутых строк, поэтому
// существование активной записи проверяется заранее.
func (s *TablesService) UpdateTableCatById(ctx context.Context, params tables_model.UpdateTableCatByIdParams) error {
	if _, err := s.GetTableCatById(ctx, params.ID); err != nil {
		return fmt.Errorf("%w: dc.table_cat: %w", customerrors.ErrUpdate, err)
	}

	if err := s.TablesRepository.UpdateTableCatById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.table_cat id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return nil
}

// DeleteTableCatById мягко удаляет строку dc.table_cat.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *TablesService) DeleteTableCatById(ctx context.Context, params tables_model.DeleteTableCatByIdParams) error {
	if _, err := s.GetTableCatById(ctx, params.ID); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.TablesRepository.DeleteTableCatById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.table_cat id = %d: %w", customerrors.ErrDelete, params.ID, err)
	}

	return nil
}

// UndeleteTableCatById восстанавливает мягко удалённую строку dc.table_cat.
func (s *TablesService) UndeleteTableCatById(ctx context.Context, params tables_model.UndeleteTableCatByIdParams) error {
	if err := s.TablesRepository.UndeleteTableCatById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.table_cat id = %d: %w", customerrors.ErrUndelete, params.ID, err)
	}

	return nil
}
