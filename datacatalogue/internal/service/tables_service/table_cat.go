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

// GetTableCats возвращает все активные строки dc.table_cat.
func (s *TablesService) GetTableCats(ctx context.Context) ([]tables_model.DcTableCat, error) {
	rows, err := s.TablesRepository.GetTableCats(ctx)

	if err != nil {
		return nil, fmt.Errorf("get dc.table_cat: %w", err)
	}

	return rows, nil
}

// GetDeletedTableCatById возвращает мягко удалённую строку dc.table_cat по id.
func (s *TablesService) GetDeletedTableCatById(ctx context.Context, id int64) (tables_model.DcTableCat, error) {
	row, err := s.TablesRepository.GetDeletedTableCatById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcTableCat{}, fmt.Errorf("deleted dc.table_cat id = %d: %w", id, customerrors.ErrNotFound)
		}

		return tables_model.DcTableCat{}, fmt.Errorf("get deleted dc.table_cat id = %d: %w", id, err)
	}

	return row, nil
}

// GetDeletedTableCats возвращает все мягко удалённые строки dc.table_cat.
func (s *TablesService) GetDeletedTableCats(ctx context.Context) ([]tables_model.DcTableCat, error) {
	rows, err := s.TablesRepository.GetDeletedTableCats(ctx)

	if err != nil {
		return nil, fmt.Errorf("get deleted dc.table_cat: %w", err)
	}

	return rows, nil
}

// GetTableCatsBySchemaId возвращает активные строки dc.table_cat, отобранные по schema_id.
func (s *TablesService) GetTableCatsBySchemaId(ctx context.Context, schemaID int64) ([]tables_model.DcTableCat, error) {
	rows, err := s.TablesRepository.GetTableCatsBySchemaId(ctx, schemaID)

	if err != nil {
		return nil, fmt.Errorf("get dc.table_cat by schema_id = %d: %w", schemaID, err)
	}

	return rows, nil
}

// GetTableCatsByTableTypeId возвращает активные строки dc.table_cat, отобранные по table_type_id.
func (s *TablesService) GetTableCatsByTableTypeId(ctx context.Context, tableTypeID int64) ([]tables_model.DcTableCat, error) {
	rows, err := s.TablesRepository.GetTableCatsByTableTypeId(ctx, tableTypeID)

	if err != nil {
		return nil, fmt.Errorf("get dc.table_cat by table_type_id = %d: %w", tableTypeID, err)
	}

	return rows, nil
}

// GetTableCatsByDomainId возвращает активные строки dc.table_cat, отобранные по domain_id.
func (s *TablesService) GetTableCatsByDomainId(ctx context.Context, domainID int64) ([]tables_model.DcTableCat, error) {
	rows, err := s.TablesRepository.GetTableCatsByDomainId(ctx, domainID)

	if err != nil {
		return nil, fmt.Errorf("get dc.table_cat by domain_id = %d: %w", domainID, err)
	}

	return rows, nil
}

// CreateTableCat вставляет строку dc.table_cat и возвращает её целиком.
func (s *TablesService) CreateTableCat(ctx context.Context, params tables_model.CreateTableCatParams) (tables_model.DcTableCat, error) {
	row, err := s.TablesRepository.CreateTableCat(ctx, params)

	if err != nil {
		return tables_model.DcTableCat{}, fmt.Errorf("%w: dc.table_cat: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateTableCatById обновляет активную строку dc.table_cat и возвращает её целиком.
//
// Запрос фильтрует по is_deleted = false, поэтому попытка обновить удалённую
// или несуществующую запись даёт sql.ErrNoRows — переводим его в ErrNotFound,
// чтобы api-слой ответил NotFound, а не Internal.
func (s *TablesService) UpdateTableCatById(ctx context.Context, params tables_model.UpdateTableCatByIdParams) (tables_model.DcTableCat, error) {
	row, err := s.TablesRepository.UpdateTableCatById(ctx, params)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcTableCat{}, fmt.Errorf("dc.table_cat id = %d: %w", params.ID, customerrors.ErrNotFound)
		}

		return tables_model.DcTableCat{}, fmt.Errorf("%w: dc.table_cat id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return row, nil
}

// DeleteTableCatById мягко удаляет строку dc.table_cat.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *TablesService) DeleteTableCatById(ctx context.Context, id int64) error {
	if _, err := s.GetTableCatById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.TablesRepository.DeleteTableCatById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.table_cat id = %d: %w", customerrors.ErrDelete, id, err)
	}

	return nil
}

// UndeleteTableCatById восстанавливает мягко удалённую строку dc.table_cat.
// Существование удалённой записи проверяется заранее по той же причине,
// что и в DeleteTableCatById.
func (s *TablesService) UndeleteTableCatById(ctx context.Context, id int64) error {
	if _, err := s.GetDeletedTableCatById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrUndelete, err)
	}

	if err := s.TablesRepository.UndeleteTableCatById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.table_cat id = %d: %w", customerrors.ErrUndelete, id, err)
	}

	return nil
}
