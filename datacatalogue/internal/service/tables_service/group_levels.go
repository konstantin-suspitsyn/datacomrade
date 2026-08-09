package tablesservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	customerrors "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/utils/custom_errors"
)

// GetGroupLevelsById возвращает активную строку dc.group_levels по id.
func (s *TablesService) GetGroupLevelsById(ctx context.Context, id int64) (tables_model.DcGroupLevel, error) {
	row, err := s.TablesRepository.GetGroupLevelsById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcGroupLevel{}, fmt.Errorf("dc.group_levels id = %d: %w", id, customerrors.ErrNotFound)
		}

		return tables_model.DcGroupLevel{}, fmt.Errorf("get dc.group_levels id = %d: %w", id, err)
	}

	return row, nil
}

// GetGroupLevelsSearchName возвращает страницу строк dc.group_levels и её счётчики.
func (s *TablesService) GetGroupLevelsSearchName(ctx context.Context, params tables_model.GetGroupLevelsSearchNameParams) ([]tables_model.DcGroupLevel, tables_model.CountGetGroupLevelsSearchNameRow, error) {
	countParams := tables_model.CountGetGroupLevelsSearchNameParams{
		PageLimit:  params.PageLimit,
		SearchName: params.SearchName,
	}
	count, err := s.TablesRepository.CountGetGroupLevelsSearchName(ctx, countParams)
	if err != nil {
		return nil, tables_model.CountGetGroupLevelsSearchNameRow{}, fmt.Errorf("count dc.group_levels: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcGroupLevel{}, count, nil
	}

	rows, err := s.TablesRepository.GetGroupLevelsSearchName(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetGroupLevelsSearchNameRow{}, fmt.Errorf("get dc.group_levels page: %w", err)
	}

	return rows, count, nil
}

// GetGroupLevels возвращает страницу строк dc.group_levels и её счётчики.
func (s *TablesService) GetGroupLevels(ctx context.Context, params tables_model.GetGroupLevelsParams) ([]tables_model.DcGroupLevel, tables_model.CountGetGroupLevelsRow, error) {
	count, err := s.TablesRepository.CountGetGroupLevels(ctx, params.PageLimit)
	if err != nil {
		return nil, tables_model.CountGetGroupLevelsRow{}, fmt.Errorf("count dc.group_levels: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcGroupLevel{}, count, nil
	}

	rows, err := s.TablesRepository.GetGroupLevels(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetGroupLevelsRow{}, fmt.Errorf("get dc.group_levels page: %w", err)
	}

	return rows, count, nil
}

// GetGroupLevelsFiltered возвращает страницу строк dc.group_levels и её счётчики.
func (s *TablesService) GetGroupLevelsFiltered(ctx context.Context, params tables_model.GetGroupLevelsFilteredParams) ([]tables_model.DcGroupLevel, tables_model.CountGetGroupLevelsFilteredRow, error) {
	countParams := tables_model.CountGetGroupLevelsFilteredParams{
		PageLimit:      params.PageLimit,
		ColumnID:       params.ColumnID,
		ParentColumnID: params.ParentColumnID,
	}
	count, err := s.TablesRepository.CountGetGroupLevelsFiltered(ctx, countParams)
	if err != nil {
		return nil, tables_model.CountGetGroupLevelsFilteredRow{}, fmt.Errorf("count dc.group_levels: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcGroupLevel{}, count, nil
	}

	rows, err := s.TablesRepository.GetGroupLevelsFiltered(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetGroupLevelsFilteredRow{}, fmt.Errorf("get dc.group_levels page: %w", err)
	}

	return rows, count, nil
}

// GetGroupLevelsDeleted возвращает строки dc.group_levels.
func (s *TablesService) GetGroupLevelsDeleted(ctx context.Context, order string) ([]tables_model.DcGroupLevel, error) {
	rows, err := s.TablesRepository.GetGroupLevelsDeleted(ctx, order)

	if err != nil {
		return nil, fmt.Errorf("GetGroupLevelsDeleted: %w", err)
	}

	return rows, nil
}

// CreateGroupLevels вставляет строку dc.group_levels и возвращает её целиком.
func (s *TablesService) CreateGroupLevels(ctx context.Context, params tables_model.CreateGroupLevelsParams) (tables_model.DcGroupLevel, error) {
	row, err := s.TablesRepository.CreateGroupLevels(ctx, params)

	if err != nil {
		return tables_model.DcGroupLevel{}, fmt.Errorf("%w: dc.group_levels: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateGroupLevelsById обновляет строку dc.group_levels.
//
// Запрос — :exec и не сообщает число затронутых строк, поэтому
// существование активной записи проверяется заранее.
func (s *TablesService) UpdateGroupLevelsById(ctx context.Context, params tables_model.UpdateGroupLevelsByIdParams) error {
	if _, err := s.GetGroupLevelsById(ctx, params.ID); err != nil {
		return fmt.Errorf("%w: dc.group_levels: %w", customerrors.ErrUpdate, err)
	}

	if err := s.TablesRepository.UpdateGroupLevelsById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.group_levels id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return nil
}

// DeleteGroupLevelsById мягко удаляет строку dc.group_levels.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *TablesService) DeleteGroupLevelsById(ctx context.Context, params tables_model.DeleteGroupLevelsByIdParams) error {
	if _, err := s.GetGroupLevelsById(ctx, params.ID); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.TablesRepository.DeleteGroupLevelsById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.group_levels id = %d: %w", customerrors.ErrDelete, params.ID, err)
	}

	return nil
}

// UndeleteGroupLevelsById восстанавливает мягко удалённую строку dc.group_levels.
func (s *TablesService) UndeleteGroupLevelsById(ctx context.Context, params tables_model.UndeleteGroupLevelsByIdParams) error {
	if err := s.TablesRepository.UndeleteGroupLevelsById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.group_levels id = %d: %w", customerrors.ErrUndelete, params.ID, err)
	}

	return nil
}
