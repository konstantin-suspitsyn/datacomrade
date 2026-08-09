package tablesservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	customerrors "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/utils/custom_errors"
)

// GetDomainCatById возвращает активную строку dc.domain_cat по id.
func (s *TablesService) GetDomainCatById(ctx context.Context, id int64) (tables_model.DcDomainCat, error) {
	row, err := s.TablesRepository.GetDomainCatById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcDomainCat{}, fmt.Errorf("dc.domain_cat id = %d: %w", id, customerrors.ErrNotFound)
		}

		return tables_model.DcDomainCat{}, fmt.Errorf("get dc.domain_cat id = %d: %w", id, err)
	}

	return row, nil
}

// GetDomainCatsSearchName возвращает страницу строк dc.domain_cat и её счётчики.
func (s *TablesService) GetDomainCatsSearchName(ctx context.Context, params tables_model.GetDomainCatsSearchNameParams) ([]tables_model.DcDomainCat, tables_model.CountGetDomainCatsSearchNameRow, error) {
	countParams := tables_model.CountGetDomainCatsSearchNameParams{
		PageLimit:  params.PageLimit,
		SearchName: params.SearchName,
	}
	count, err := s.TablesRepository.CountGetDomainCatsSearchName(ctx, countParams)
	if err != nil {
		return nil, tables_model.CountGetDomainCatsSearchNameRow{}, fmt.Errorf("count dc.domain_cat: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcDomainCat{}, count, nil
	}

	rows, err := s.TablesRepository.GetDomainCatsSearchName(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetDomainCatsSearchNameRow{}, fmt.Errorf("get dc.domain_cat page: %w", err)
	}

	return rows, count, nil
}

// GetDomainCats возвращает страницу строк dc.domain_cat и её счётчики.
func (s *TablesService) GetDomainCats(ctx context.Context, params tables_model.GetDomainCatsParams) ([]tables_model.DcDomainCat, tables_model.CountGetDomainCatsRow, error) {
	count, err := s.TablesRepository.CountGetDomainCats(ctx, params.PageLimit)
	if err != nil {
		return nil, tables_model.CountGetDomainCatsRow{}, fmt.Errorf("count dc.domain_cat: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcDomainCat{}, count, nil
	}

	rows, err := s.TablesRepository.GetDomainCats(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetDomainCatsRow{}, fmt.Errorf("get dc.domain_cat page: %w", err)
	}

	return rows, count, nil
}

// GetDomainCatsDeleted возвращает страницу строк dc.domain_cat и её счётчики.
func (s *TablesService) GetDomainCatsDeleted(ctx context.Context, params tables_model.GetDomainCatsDeletedParams) ([]tables_model.DcDomainCat, tables_model.CountGetDomainCatsDeletedRow, error) {
	count, err := s.TablesRepository.CountGetDomainCatsDeleted(ctx, params.PageLimit)
	if err != nil {
		return nil, tables_model.CountGetDomainCatsDeletedRow{}, fmt.Errorf("count dc.domain_cat: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcDomainCat{}, count, nil
	}

	rows, err := s.TablesRepository.GetDomainCatsDeleted(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetDomainCatsDeletedRow{}, fmt.Errorf("get dc.domain_cat page: %w", err)
	}

	return rows, count, nil
}

// CreateDomainCat вставляет строку dc.domain_cat и возвращает её целиком.
func (s *TablesService) CreateDomainCat(ctx context.Context, params tables_model.CreateDomainCatParams) (tables_model.DcDomainCat, error) {
	row, err := s.TablesRepository.CreateDomainCat(ctx, params)

	if err != nil {
		return tables_model.DcDomainCat{}, fmt.Errorf("%w: dc.domain_cat: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateDomainCatById обновляет строку dc.domain_cat.
//
// Запрос — :exec и не сообщает число затронутых строк, поэтому
// существование активной записи проверяется заранее.
func (s *TablesService) UpdateDomainCatById(ctx context.Context, params tables_model.UpdateDomainCatByIdParams) error {
	if _, err := s.GetDomainCatById(ctx, params.ID); err != nil {
		return fmt.Errorf("%w: dc.domain_cat: %w", customerrors.ErrUpdate, err)
	}

	if err := s.TablesRepository.UpdateDomainCatById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.domain_cat id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return nil
}

// DeleteDomainCatById мягко удаляет строку dc.domain_cat.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *TablesService) DeleteDomainCatById(ctx context.Context, params tables_model.DeleteDomainCatByIdParams) error {
	if _, err := s.GetDomainCatById(ctx, params.ID); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.TablesRepository.DeleteDomainCatById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.domain_cat id = %d: %w", customerrors.ErrDelete, params.ID, err)
	}

	return nil
}

// UndeleteDomainCatById восстанавливает мягко удалённую строку dc.domain_cat.
func (s *TablesService) UndeleteDomainCatById(ctx context.Context, params tables_model.UndeleteDomainCatByIdParams) error {
	if err := s.TablesRepository.UndeleteDomainCatById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.domain_cat id = %d: %w", customerrors.ErrUndelete, params.ID, err)
	}

	return nil
}
