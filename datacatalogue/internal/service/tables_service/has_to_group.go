package tablesservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	customerrors "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/utils/custom_errors"
)

// GetHasToGroupById возвращает активную строку dc.has_to_group по id.
func (s *TablesService) GetHasToGroupById(ctx context.Context, id int64) (tables_model.DcHasToGroup, error) {
	row, err := s.TablesRepository.GetHasToGroupById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcHasToGroup{}, fmt.Errorf("dc.has_to_group id = %d: %w", id, customerrors.ErrNotFound)
		}

		return tables_model.DcHasToGroup{}, fmt.Errorf("get dc.has_to_group id = %d: %w", id, err)
	}

	return row, nil
}

// GetHasToGroups возвращает все активные строки dc.has_to_group.
func (s *TablesService) GetHasToGroups(ctx context.Context) ([]tables_model.DcHasToGroup, error) {
	rows, err := s.TablesRepository.GetHasToGroups(ctx)

	if err != nil {
		return nil, fmt.Errorf("get dc.has_to_group: %w", err)
	}

	return rows, nil
}

// GetDeletedHasToGroupById возвращает мягко удалённую строку dc.has_to_group по id.
func (s *TablesService) GetDeletedHasToGroupById(ctx context.Context, id int64) (tables_model.DcHasToGroup, error) {
	row, err := s.TablesRepository.GetDeletedHasToGroupById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcHasToGroup{}, fmt.Errorf("deleted dc.has_to_group id = %d: %w", id, customerrors.ErrNotFound)
		}

		return tables_model.DcHasToGroup{}, fmt.Errorf("get deleted dc.has_to_group id = %d: %w", id, err)
	}

	return row, nil
}

// GetDeletedHasToGroups возвращает все мягко удалённые строки dc.has_to_group.
func (s *TablesService) GetDeletedHasToGroups(ctx context.Context) ([]tables_model.DcHasToGroup, error) {
	rows, err := s.TablesRepository.GetDeletedHasToGroups(ctx)

	if err != nil {
		return nil, fmt.Errorf("get deleted dc.has_to_group: %w", err)
	}

	return rows, nil
}

// CreateHasToGroup вставляет строку dc.has_to_group и возвращает её целиком.
func (s *TablesService) CreateHasToGroup(ctx context.Context, params tables_model.CreateHasToGroupParams) (tables_model.DcHasToGroup, error) {
	row, err := s.TablesRepository.CreateHasToGroup(ctx, params)

	if err != nil {
		return tables_model.DcHasToGroup{}, fmt.Errorf("%w: dc.has_to_group: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateHasToGroupById обновляет активную строку dc.has_to_group и возвращает её целиком.
//
// Запрос фильтрует по is_deleted = false, поэтому попытка обновить удалённую
// или несуществующую запись даёт sql.ErrNoRows — переводим его в ErrNotFound,
// чтобы api-слой ответил NotFound, а не Internal.
func (s *TablesService) UpdateHasToGroupById(ctx context.Context, params tables_model.UpdateHasToGroupByIdParams) (tables_model.DcHasToGroup, error) {
	row, err := s.TablesRepository.UpdateHasToGroupById(ctx, params)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcHasToGroup{}, fmt.Errorf("dc.has_to_group id = %d: %w", params.ID, customerrors.ErrNotFound)
		}

		return tables_model.DcHasToGroup{}, fmt.Errorf("%w: dc.has_to_group id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return row, nil
}

// DeleteHasToGroupById мягко удаляет строку dc.has_to_group.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *TablesService) DeleteHasToGroupById(ctx context.Context, id int64) error {
	if _, err := s.GetHasToGroupById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.TablesRepository.DeleteHasToGroupById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.has_to_group id = %d: %w", customerrors.ErrDelete, id, err)
	}

	return nil
}

// UndeleteHasToGroupById восстанавливает мягко удалённую строку dc.has_to_group.
// Существование удалённой записи проверяется заранее по той же причине,
// что и в DeleteHasToGroupById.
func (s *TablesService) UndeleteHasToGroupById(ctx context.Context, id int64) error {
	if _, err := s.GetDeletedHasToGroupById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrUndelete, err)
	}

	if err := s.TablesRepository.UndeleteHasToGroupById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.has_to_group id = %d: %w", customerrors.ErrUndelete, id, err)
	}

	return nil
}
