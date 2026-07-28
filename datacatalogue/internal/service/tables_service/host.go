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

// GetHosts возвращает все активные строки dc.host.
func (s *TablesService) GetHosts(ctx context.Context) ([]tables_model.DcHost, error) {
	rows, err := s.TablesRepository.GetHosts(ctx)

	if err != nil {
		return nil, fmt.Errorf("get dc.host: %w", err)
	}

	return rows, nil
}

// GetDeletedHostById возвращает мягко удалённую строку dc.host по id.
func (s *TablesService) GetDeletedHostById(ctx context.Context, id int64) (tables_model.DcHost, error) {
	row, err := s.TablesRepository.GetDeletedHostById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcHost{}, fmt.Errorf("deleted dc.host id = %d: %w", id, customerrors.ErrNotFound)
		}

		return tables_model.DcHost{}, fmt.Errorf("get deleted dc.host id = %d: %w", id, err)
	}

	return row, nil
}

// GetDeletedHosts возвращает все мягко удалённые строки dc.host.
func (s *TablesService) GetDeletedHosts(ctx context.Context) ([]tables_model.DcHost, error) {
	rows, err := s.TablesRepository.GetDeletedHosts(ctx)

	if err != nil {
		return nil, fmt.Errorf("get deleted dc.host: %w", err)
	}

	return rows, nil
}

// CreateHost вставляет строку dc.host и возвращает её целиком.
func (s *TablesService) CreateHost(ctx context.Context, params tables_model.CreateHostParams) (tables_model.DcHost, error) {
	row, err := s.TablesRepository.CreateHost(ctx, params)

	if err != nil {
		return tables_model.DcHost{}, fmt.Errorf("%w: dc.host: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateHostById обновляет активную строку dc.host и возвращает её целиком.
//
// Запрос фильтрует по is_deleted = false, поэтому попытка обновить удалённую
// или несуществующую запись даёт sql.ErrNoRows — переводим его в ErrNotFound,
// чтобы api-слой ответил NotFound, а не Internal.
func (s *TablesService) UpdateHostById(ctx context.Context, params tables_model.UpdateHostByIdParams) (tables_model.DcHost, error) {
	row, err := s.TablesRepository.UpdateHostById(ctx, params)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcHost{}, fmt.Errorf("dc.host id = %d: %w", params.ID, customerrors.ErrNotFound)
		}

		return tables_model.DcHost{}, fmt.Errorf("%w: dc.host id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return row, nil
}

// DeleteHostById мягко удаляет строку dc.host.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *TablesService) DeleteHostById(ctx context.Context, id int64) error {
	if _, err := s.GetHostById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.TablesRepository.DeleteHostById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.host id = %d: %w", customerrors.ErrDelete, id, err)
	}

	return nil
}

// UndeleteHostById восстанавливает мягко удалённую строку dc.host.
// Существование удалённой записи проверяется заранее по той же причине,
// что и в DeleteHostById.
func (s *TablesService) UndeleteHostById(ctx context.Context, id int64) error {
	if _, err := s.GetDeletedHostById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrUndelete, err)
	}

	if err := s.TablesRepository.UndeleteHostById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.host id = %d: %w", customerrors.ErrUndelete, id, err)
	}

	return nil
}
