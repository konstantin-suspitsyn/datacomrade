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

// GetColumnCats возвращает все активные строки dc.column_cat.
func (s *TablesService) GetColumnCats(ctx context.Context) ([]tables_model.DcColumnCat, error) {
	rows, err := s.TablesRepository.GetColumnCats(ctx)

	if err != nil {
		return nil, fmt.Errorf("get dc.column_cat: %w", err)
	}

	return rows, nil
}

// GetDeletedColumnCatById возвращает мягко удалённую строку dc.column_cat по id.
func (s *TablesService) GetDeletedColumnCatById(ctx context.Context, id int64) (tables_model.DcColumnCat, error) {
	row, err := s.TablesRepository.GetDeletedColumnCatById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcColumnCat{}, fmt.Errorf("deleted dc.column_cat id = %d: %w", id, customerrors.ErrNotFound)
		}

		return tables_model.DcColumnCat{}, fmt.Errorf("get deleted dc.column_cat id = %d: %w", id, err)
	}

	return row, nil
}

// GetDeletedColumnCats возвращает все мягко удалённые строки dc.column_cat.
func (s *TablesService) GetDeletedColumnCats(ctx context.Context) ([]tables_model.DcColumnCat, error) {
	rows, err := s.TablesRepository.GetDeletedColumnCats(ctx)

	if err != nil {
		return nil, fmt.Errorf("get deleted dc.column_cat: %w", err)
	}

	return rows, nil
}

// GetColumnCatsByTableId возвращает активные строки dc.column_cat, отобранные по table_id.
func (s *TablesService) GetColumnCatsByTableId(ctx context.Context, tableID int64) ([]tables_model.DcColumnCat, error) {
	rows, err := s.TablesRepository.GetColumnCatsByTableId(ctx, tableID)

	if err != nil {
		return nil, fmt.Errorf("get dc.column_cat by table_id = %d: %w", tableID, err)
	}

	return rows, nil
}

// GetColumnCatsByAliasId возвращает активные строки dc.column_cat, отобранные по alias_id.
func (s *TablesService) GetColumnCatsByAliasId(ctx context.Context, aliasID int64) ([]tables_model.DcColumnCat, error) {
	rows, err := s.TablesRepository.GetColumnCatsByAliasId(ctx, aliasID)

	if err != nil {
		return nil, fmt.Errorf("get dc.column_cat by alias_id = %d: %w", aliasID, err)
	}

	return rows, nil
}

// CreateColumnCat вставляет строку dc.column_cat и возвращает её целиком.
func (s *TablesService) CreateColumnCat(ctx context.Context, params tables_model.CreateColumnCatParams) (tables_model.DcColumnCat, error) {
	row, err := s.TablesRepository.CreateColumnCat(ctx, params)

	if err != nil {
		return tables_model.DcColumnCat{}, fmt.Errorf("%w: dc.column_cat: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateColumnCatById обновляет активную строку dc.column_cat и возвращает её целиком.
//
// Запрос фильтрует по is_deleted = false, поэтому попытка обновить удалённую
// или несуществующую запись даёт sql.ErrNoRows — переводим его в ErrNotFound,
// чтобы api-слой ответил NotFound, а не Internal.
func (s *TablesService) UpdateColumnCatById(ctx context.Context, params tables_model.UpdateColumnCatByIdParams) (tables_model.DcColumnCat, error) {
	row, err := s.TablesRepository.UpdateColumnCatById(ctx, params)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcColumnCat{}, fmt.Errorf("dc.column_cat id = %d: %w", params.ID, customerrors.ErrNotFound)
		}

		return tables_model.DcColumnCat{}, fmt.Errorf("%w: dc.column_cat id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return row, nil
}

// DeleteColumnCatById мягко удаляет строку dc.column_cat.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *TablesService) DeleteColumnCatById(ctx context.Context, id int64) error {
	if _, err := s.GetColumnCatById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.TablesRepository.DeleteColumnCatById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.column_cat id = %d: %w", customerrors.ErrDelete, id, err)
	}

	return nil
}

// UndeleteColumnCatById восстанавливает мягко удалённую строку dc.column_cat.
// Существование удалённой записи проверяется заранее по той же причине,
// что и в DeleteColumnCatById.
func (s *TablesService) UndeleteColumnCatById(ctx context.Context, id int64) error {
	if _, err := s.GetDeletedColumnCatById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrUndelete, err)
	}

	if err := s.TablesRepository.UndeleteColumnCatById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.column_cat id = %d: %w", customerrors.ErrUndelete, id, err)
	}

	return nil
}
