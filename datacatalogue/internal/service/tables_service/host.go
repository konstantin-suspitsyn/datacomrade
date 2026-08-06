package tablesservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	customerrors "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/utils/custom_errors"
)

// GetHostById возвращает активную строку dc.host по id.
func (s *TablesService) GetHostById(ctx context.Context, id int64) (tables_model.DcHost, error) {
	row, err := s.TablesRepository.GetHostById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcHost{}, fmt.Errorf("dc.host id = %d: %w", id, customerrors.ErrNotFound)
		}

		return tables_model.DcHost{}, fmt.Errorf("get dc.host id = %d: %w", id, err)
	}

	return row, nil
}

// GetHosts возвращает страницу строк dc.host и её счётчики.
func (s *TablesService) GetHosts(ctx context.Context, params tables_model.GetHostsParams) ([]tables_model.DcHost, tables_model.CountGetHostsRow, error) {
	count, err := s.TablesRepository.CountGetHosts(ctx, params.PageLimit)
	if err != nil {
		return nil, tables_model.CountGetHostsRow{}, fmt.Errorf("count dc.host: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcHost{}, count, nil
	}

	rows, err := s.TablesRepository.GetHosts(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetHostsRow{}, fmt.Errorf("get dc.host page: %w", err)
	}

	return rows, count, nil
}

// GetHostsSearchName возвращает страницу строк dc.host и её счётчики.
func (s *TablesService) GetHostsSearchName(ctx context.Context, params tables_model.GetHostsSearchNameParams) ([]tables_model.DcHost, tables_model.CountGetHostsSearchNameRow, error) {
	countParams := tables_model.CountGetHostsSearchNameParams{
		PageLimit:  params.PageLimit,
		SearchName: params.SearchName,
	}
	count, err := s.TablesRepository.CountGetHostsSearchName(ctx, countParams)
	if err != nil {
		return nil, tables_model.CountGetHostsSearchNameRow{}, fmt.Errorf("count dc.host: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcHost{}, count, nil
	}

	rows, err := s.TablesRepository.GetHostsSearchName(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetHostsSearchNameRow{}, fmt.Errorf("get dc.host page: %w", err)
	}

	return rows, count, nil
}

// GetHostDeleted возвращает страницу строк dc.host и её счётчики.
func (s *TablesService) GetHostDeleted(ctx context.Context, params tables_model.GetHostDeletedParams) ([]tables_model.DcHost, tables_model.CountGetHostDeletedRow, error) {
	count, err := s.TablesRepository.CountGetHostDeleted(ctx, params.PageLimit)
	if err != nil {
		return nil, tables_model.CountGetHostDeletedRow{}, fmt.Errorf("count dc.host: %w", err)
	}

	if count.TotalItems == 0 {
		return []tables_model.DcHost{}, count, nil
	}

	rows, err := s.TablesRepository.GetHostDeleted(ctx, params)
	if err != nil {
		return nil, tables_model.CountGetHostDeletedRow{}, fmt.Errorf("get dc.host page: %w", err)
	}

	return rows, count, nil
}

// CreateHost вставляет строку dc.host и возвращает её целиком.
func (s *TablesService) CreateHost(ctx context.Context, params tables_model.CreateHostParams) (tables_model.DcHost, error) {
	row, err := s.TablesRepository.CreateHost(ctx, params)

	if err != nil {
		return tables_model.DcHost{}, fmt.Errorf("%w: dc.host: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateHostById обновляет строку dc.host.
//
// Запрос — :exec и не сообщает число затронутых строк, поэтому
// существование активной записи проверяется заранее.
func (s *TablesService) UpdateHostById(ctx context.Context, params tables_model.UpdateHostByIdParams) error {
	if _, err := s.GetHostById(ctx, params.ID); err != nil {
		return fmt.Errorf("%w: dc.host: %w", customerrors.ErrUpdate, err)
	}

	if err := s.TablesRepository.UpdateHostById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.host id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return nil
}

// DeleteHostById мягко удаляет строку dc.host.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *TablesService) DeleteHostById(ctx context.Context, params tables_model.DeleteHostByIdParams) error {
	if _, err := s.GetHostById(ctx, params.ID); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.TablesRepository.DeleteHostById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.host id = %d: %w", customerrors.ErrDelete, params.ID, err)
	}

	return nil
}

// UndeleteHostById восстанавливает мягко удалённую строку dc.host.
func (s *TablesService) UndeleteHostById(ctx context.Context, params tables_model.UndeleteHostByIdParams) error {
	if err := s.TablesRepository.UndeleteHostById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.host id = %d: %w", customerrors.ErrUndelete, params.ID, err)
	}

	return nil
}
