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

// GetHasToGroupsSearchName возвращает страницу строк dc.has_to_group и её счётчики.
func (s *TablesService) GetHasToGroupsSearchName(ctx context.Context, params tables_model.GetHasToGroupsSearchNameParams) ([]tables_model.DcHasToGroup, tables_model.CountGetHasToGroupsSearchNameRow, error) {
	countParams := tables_model.CountGetHasToGroupsSearchNameParams{
		PageLimit:  params.PageLimit,
		SearchName: params.SearchName,
	}
	count, err := s.TablesRepository.CountGetHasToGroupsSearchName(ctx, countParams)
	if err != nil {
		return nil, tables_model.CountGetHasToGroupsSearchNameRow{}, fmt.Errorf("count dc.has_to_group: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcHasToGroup{}, count, nil
	}

	rows, err := s.TablesRepository.GetHasToGroupsSearchName(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetHasToGroupsSearchNameRow{}, fmt.Errorf("get dc.has_to_group page: %w", err)
	}

	return rows, count, nil
}

// GetHasToGroups возвращает страницу строк dc.has_to_group и её счётчики.
func (s *TablesService) GetHasToGroups(ctx context.Context, params tables_model.GetHasToGroupsParams) ([]tables_model.DcHasToGroup, tables_model.CountGetHasToGroupsRow, error) {
	count, err := s.TablesRepository.CountGetHasToGroups(ctx, params.PageLimit)
	if err != nil {
		return nil, tables_model.CountGetHasToGroupsRow{}, fmt.Errorf("count dc.has_to_group: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcHasToGroup{}, count, nil
	}

	rows, err := s.TablesRepository.GetHasToGroups(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetHasToGroupsRow{}, fmt.Errorf("get dc.has_to_group page: %w", err)
	}

	return rows, count, nil
}

// GetHasToGroupsDeleted возвращает страницу строк dc.has_to_group и её счётчики.
func (s *TablesService) GetHasToGroupsDeleted(ctx context.Context, params tables_model.GetHasToGroupsDeletedParams) ([]tables_model.DcHasToGroup, tables_model.CountGetHasToGroupsDeletedRow, error) {
	count, err := s.TablesRepository.CountGetHasToGroupsDeleted(ctx, params.PageLimit)
	if err != nil {
		return nil, tables_model.CountGetHasToGroupsDeletedRow{}, fmt.Errorf("count dc.has_to_group: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcHasToGroup{}, count, nil
	}

	rows, err := s.TablesRepository.GetHasToGroupsDeleted(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetHasToGroupsDeletedRow{}, fmt.Errorf("get dc.has_to_group page: %w", err)
	}

	return rows, count, nil
}

// CreateHasToGroup вставляет строку dc.has_to_group и возвращает её целиком.
func (s *TablesService) CreateHasToGroup(ctx context.Context, params tables_model.CreateHasToGroupParams) (tables_model.DcHasToGroup, error) {
	row, err := s.TablesRepository.CreateHasToGroup(ctx, params)

	if err != nil {
		return tables_model.DcHasToGroup{}, fmt.Errorf("%w: dc.has_to_group: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateHasToGroupById обновляет строку dc.has_to_group.
//
// Запрос — :exec и не сообщает число затронутых строк, поэтому
// существование активной записи проверяется заранее.
func (s *TablesService) UpdateHasToGroupById(ctx context.Context, params tables_model.UpdateHasToGroupByIdParams) error {
	if _, err := s.GetHasToGroupById(ctx, params.ID); err != nil {
		return fmt.Errorf("%w: dc.has_to_group: %w", customerrors.ErrUpdate, err)
	}

	if err := s.TablesRepository.UpdateHasToGroupById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.has_to_group id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return nil
}

// DeleteHasToGroupById мягко удаляет строку dc.has_to_group.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *TablesService) DeleteHasToGroupById(ctx context.Context, params tables_model.DeleteHasToGroupByIdParams) error {
	if _, err := s.GetHasToGroupById(ctx, params.ID); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.TablesRepository.DeleteHasToGroupById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.has_to_group id = %d: %w", customerrors.ErrDelete, params.ID, err)
	}

	return nil
}

// UndeleteHasToGroupById восстанавливает мягко удалённую строку dc.has_to_group.
func (s *TablesService) UndeleteHasToGroupById(ctx context.Context, params tables_model.UndeleteHasToGroupByIdParams) error {
	if err := s.TablesRepository.UndeleteHasToGroupById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.has_to_group id = %d: %w", customerrors.ErrUndelete, params.ID, err)
	}

	return nil
}
