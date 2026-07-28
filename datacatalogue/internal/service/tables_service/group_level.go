package tablesservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	customerrors "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/utils/custom_errors"
)

// GetGroupLevelById возвращает активную строку dc.group_levels по id.
func (s *TablesService) GetGroupLevelById(ctx context.Context, id int64) (tables_model.DcGroupLevel, error) {
	row, err := s.TablesRepository.GetGroupLevelById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcGroupLevel{}, fmt.Errorf("dc.group_levels id = %d: %w", id, customerrors.ErrNotFound)
		}

		return tables_model.DcGroupLevel{}, fmt.Errorf("get dc.group_levels id = %d: %w", id, err)
	}

	return row, nil
}

// GetGroupLevels возвращает все активные строки dc.group_levels.
func (s *TablesService) GetGroupLevels(ctx context.Context) ([]tables_model.DcGroupLevel, error) {
	rows, err := s.TablesRepository.GetGroupLevels(ctx)

	if err != nil {
		return nil, fmt.Errorf("get dc.group_levels: %w", err)
	}

	return rows, nil
}

// GetDeletedGroupLevelById возвращает мягко удалённую строку dc.group_levels по id.
func (s *TablesService) GetDeletedGroupLevelById(ctx context.Context, id int64) (tables_model.DcGroupLevel, error) {
	row, err := s.TablesRepository.GetDeletedGroupLevelById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcGroupLevel{}, fmt.Errorf("deleted dc.group_levels id = %d: %w", id, customerrors.ErrNotFound)
		}

		return tables_model.DcGroupLevel{}, fmt.Errorf("get deleted dc.group_levels id = %d: %w", id, err)
	}

	return row, nil
}

// GetDeletedGroupLevels возвращает все мягко удалённые строки dc.group_levels.
func (s *TablesService) GetDeletedGroupLevels(ctx context.Context) ([]tables_model.DcGroupLevel, error) {
	rows, err := s.TablesRepository.GetDeletedGroupLevels(ctx)

	if err != nil {
		return nil, fmt.Errorf("get deleted dc.group_levels: %w", err)
	}

	return rows, nil
}

// CreateGroupLevel вставляет строку dc.group_levels и возвращает её целиком.
func (s *TablesService) CreateGroupLevel(ctx context.Context, params tables_model.CreateGroupLevelParams) (tables_model.DcGroupLevel, error) {
	row, err := s.TablesRepository.CreateGroupLevel(ctx, params)

	if err != nil {
		return tables_model.DcGroupLevel{}, fmt.Errorf("%w: dc.group_levels: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateGroupLevelById обновляет активную строку dc.group_levels и возвращает её целиком.
//
// Запрос фильтрует по is_deleted = false, поэтому попытка обновить удалённую
// или несуществующую запись даёт sql.ErrNoRows — переводим его в ErrNotFound,
// чтобы api-слой ответил NotFound, а не Internal.
func (s *TablesService) UpdateGroupLevelById(ctx context.Context, params tables_model.UpdateGroupLevelByIdParams) (tables_model.DcGroupLevel, error) {
	row, err := s.TablesRepository.UpdateGroupLevelById(ctx, params)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcGroupLevel{}, fmt.Errorf("dc.group_levels id = %d: %w", params.ID, customerrors.ErrNotFound)
		}

		return tables_model.DcGroupLevel{}, fmt.Errorf("%w: dc.group_levels id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return row, nil
}

// DeleteGroupLevelById мягко удаляет строку dc.group_levels.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *TablesService) DeleteGroupLevelById(ctx context.Context, id int64) error {
	if _, err := s.GetGroupLevelById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.TablesRepository.DeleteGroupLevelById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.group_levels id = %d: %w", customerrors.ErrDelete, id, err)
	}

	return nil
}

// UndeleteGroupLevelById восстанавливает мягко удалённую строку dc.group_levels.
// Существование удалённой записи проверяется заранее по той же причине,
// что и в DeleteGroupLevelById.
func (s *TablesService) UndeleteGroupLevelById(ctx context.Context, id int64) error {
	if _, err := s.GetDeletedGroupLevelById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrUndelete, err)
	}

	if err := s.TablesRepository.UndeleteGroupLevelById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.group_levels id = %d: %w", customerrors.ErrUndelete, id, err)
	}

	return nil
}
