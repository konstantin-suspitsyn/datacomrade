package tablesservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	customerrors "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/utils/custom_errors"
)

// GetAliasById возвращает активную строку dc.alias по id.
func (s *TablesService) GetAliasById(ctx context.Context, id int64) (tables_model.DcAlias, error) {
	row, err := s.TablesRepository.GetAliasById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcAlias{}, fmt.Errorf("dc.alias id = %d: %w", id, customerrors.ErrNotFound)
		}

		return tables_model.DcAlias{}, fmt.Errorf("get dc.alias id = %d: %w", id, err)
	}

	return row, nil
}

// GetAliases возвращает все активные строки dc.alias.
func (s *TablesService) GetAliases(ctx context.Context) ([]tables_model.DcAlias, error) {
	rows, err := s.TablesRepository.GetAliases(ctx)

	if err != nil {
		return nil, fmt.Errorf("get dc.alias: %w", err)
	}

	return rows, nil
}

// GetDeletedAliasById возвращает мягко удалённую строку dc.alias по id.
func (s *TablesService) GetDeletedAliasById(ctx context.Context, id int64) (tables_model.DcAlias, error) {
	row, err := s.TablesRepository.GetDeletedAliasById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcAlias{}, fmt.Errorf("deleted dc.alias id = %d: %w", id, customerrors.ErrNotFound)
		}

		return tables_model.DcAlias{}, fmt.Errorf("get deleted dc.alias id = %d: %w", id, err)
	}

	return row, nil
}

// GetDeletedAliases возвращает все мягко удалённые строки dc.alias.
func (s *TablesService) GetDeletedAliases(ctx context.Context) ([]tables_model.DcAlias, error) {
	rows, err := s.TablesRepository.GetDeletedAliases(ctx)

	if err != nil {
		return nil, fmt.Errorf("get deleted dc.alias: %w", err)
	}

	return rows, nil
}

// CreateAlias вставляет строку dc.alias и возвращает её целиком.
func (s *TablesService) CreateAlias(ctx context.Context, params tables_model.CreateAliasParams) (tables_model.DcAlias, error) {
	row, err := s.TablesRepository.CreateAlias(ctx, params)

	if err != nil {
		return tables_model.DcAlias{}, fmt.Errorf("%w: dc.alias: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateAliasById обновляет активную строку dc.alias и возвращает её целиком.
//
// Запрос фильтрует по is_deleted = false, поэтому попытка обновить удалённую
// или несуществующую запись даёт sql.ErrNoRows — переводим его в ErrNotFound,
// чтобы api-слой ответил NotFound, а не Internal.
func (s *TablesService) UpdateAliasById(ctx context.Context, params tables_model.UpdateAliasByIdParams) (tables_model.DcAlias, error) {
	row, err := s.TablesRepository.UpdateAliasById(ctx, params)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcAlias{}, fmt.Errorf("dc.alias id = %d: %w", params.ID, customerrors.ErrNotFound)
		}

		return tables_model.DcAlias{}, fmt.Errorf("%w: dc.alias id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return row, nil
}

// DeleteAliasById мягко удаляет строку dc.alias.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *TablesService) DeleteAliasById(ctx context.Context, id int64) error {
	if _, err := s.GetAliasById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.TablesRepository.DeleteAliasById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.alias id = %d: %w", customerrors.ErrDelete, id, err)
	}

	return nil
}

// UndeleteAliasById восстанавливает мягко удалённую строку dc.alias.
// Существование удалённой записи проверяется заранее по той же причине,
// что и в DeleteAliasById.
func (s *TablesService) UndeleteAliasById(ctx context.Context, id int64) error {
	if _, err := s.GetDeletedAliasById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrUndelete, err)
	}

	if err := s.TablesRepository.UndeleteAliasById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.alias id = %d: %w", customerrors.ErrUndelete, id, err)
	}

	return nil
}
