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

// GetDomainCats возвращает все активные строки dc.domain_cat.
func (s *TablesService) GetDomainCats(ctx context.Context) ([]tables_model.DcDomainCat, error) {
	rows, err := s.TablesRepository.GetDomainCats(ctx)

	if err != nil {
		return nil, fmt.Errorf("get dc.domain_cat: %w", err)
	}

	return rows, nil
}

// GetDeletedDomainCatById возвращает мягко удалённую строку dc.domain_cat по id.
func (s *TablesService) GetDeletedDomainCatById(ctx context.Context, id int64) (tables_model.DcDomainCat, error) {
	row, err := s.TablesRepository.GetDeletedDomainCatById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcDomainCat{}, fmt.Errorf("deleted dc.domain_cat id = %d: %w", id, customerrors.ErrNotFound)
		}

		return tables_model.DcDomainCat{}, fmt.Errorf("get deleted dc.domain_cat id = %d: %w", id, err)
	}

	return row, nil
}

// GetDeletedDomainCats возвращает все мягко удалённые строки dc.domain_cat.
func (s *TablesService) GetDeletedDomainCats(ctx context.Context) ([]tables_model.DcDomainCat, error) {
	rows, err := s.TablesRepository.GetDeletedDomainCats(ctx)

	if err != nil {
		return nil, fmt.Errorf("get deleted dc.domain_cat: %w", err)
	}

	return rows, nil
}

// CreateDomainCat вставляет строку dc.domain_cat и возвращает её целиком.
func (s *TablesService) CreateDomainCat(ctx context.Context, params tables_model.CreateDomainCatParams) (tables_model.DcDomainCat, error) {
	row, err := s.TablesRepository.CreateDomainCat(ctx, params)

	if err != nil {
		return tables_model.DcDomainCat{}, fmt.Errorf("%w: dc.domain_cat: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateDomainCatById обновляет активную строку dc.domain_cat и возвращает её целиком.
//
// Запрос фильтрует по is_deleted = false, поэтому попытка обновить удалённую
// или несуществующую запись даёт sql.ErrNoRows — переводим его в ErrNotFound,
// чтобы api-слой ответил NotFound, а не Internal.
func (s *TablesService) UpdateDomainCatById(ctx context.Context, params tables_model.UpdateDomainCatByIdParams) (tables_model.DcDomainCat, error) {
	row, err := s.TablesRepository.UpdateDomainCatById(ctx, params)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcDomainCat{}, fmt.Errorf("dc.domain_cat id = %d: %w", params.ID, customerrors.ErrNotFound)
		}

		return tables_model.DcDomainCat{}, fmt.Errorf("%w: dc.domain_cat id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return row, nil
}

// DeleteDomainCatById мягко удаляет строку dc.domain_cat.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *TablesService) DeleteDomainCatById(ctx context.Context, id int64) error {
	if _, err := s.GetDomainCatById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.TablesRepository.DeleteDomainCatById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.domain_cat id = %d: %w", customerrors.ErrDelete, id, err)
	}

	return nil
}

// UndeleteDomainCatById восстанавливает мягко удалённую строку dc.domain_cat.
// Существование удалённой записи проверяется заранее по той же причине,
// что и в DeleteDomainCatById.
func (s *TablesService) UndeleteDomainCatById(ctx context.Context, id int64) error {
	if _, err := s.GetDeletedDomainCatById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrUndelete, err)
	}

	if err := s.TablesRepository.UndeleteDomainCatById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.domain_cat id = %d: %w", customerrors.ErrUndelete, id, err)
	}

	return nil
}
