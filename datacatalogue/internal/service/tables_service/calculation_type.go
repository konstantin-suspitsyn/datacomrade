package tablesservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	customerrors "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/utils/custom_errors"
)

// GetCalculationTypeById возвращает активную строку dc.calculation_type по id.
func (s *TablesService) GetCalculationTypeById(ctx context.Context, id int64) (tables_model.DcCalculationType, error) {
	row, err := s.TablesRepository.GetCalculationTypeById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcCalculationType{}, fmt.Errorf("dc.calculation_type id = %d: %w", id, customerrors.ErrNotFound)
		}

		return tables_model.DcCalculationType{}, fmt.Errorf("get dc.calculation_type id = %d: %w", id, err)
	}

	return row, nil
}

// GetCalculationTypes возвращает все активные строки dc.calculation_type.
func (s *TablesService) GetCalculationTypes(ctx context.Context) ([]tables_model.DcCalculationType, error) {
	rows, err := s.TablesRepository.GetCalculationTypes(ctx)

	if err != nil {
		return nil, fmt.Errorf("get dc.calculation_type: %w", err)
	}

	return rows, nil
}

// GetDeletedCalculationTypeById возвращает мягко удалённую строку dc.calculation_type по id.
func (s *TablesService) GetDeletedCalculationTypeById(ctx context.Context, id int64) (tables_model.DcCalculationType, error) {
	row, err := s.TablesRepository.GetDeletedCalculationTypeById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcCalculationType{}, fmt.Errorf("deleted dc.calculation_type id = %d: %w", id, customerrors.ErrNotFound)
		}

		return tables_model.DcCalculationType{}, fmt.Errorf("get deleted dc.calculation_type id = %d: %w", id, err)
	}

	return row, nil
}

// GetDeletedCalculationTypes возвращает все мягко удалённые строки dc.calculation_type.
func (s *TablesService) GetDeletedCalculationTypes(ctx context.Context) ([]tables_model.DcCalculationType, error) {
	rows, err := s.TablesRepository.GetDeletedCalculationTypes(ctx)

	if err != nil {
		return nil, fmt.Errorf("get deleted dc.calculation_type: %w", err)
	}

	return rows, nil
}

// CreateCalculationType вставляет строку dc.calculation_type и возвращает её целиком.
func (s *TablesService) CreateCalculationType(ctx context.Context, params tables_model.CreateCalculationTypeParams) (tables_model.DcCalculationType, error) {
	row, err := s.TablesRepository.CreateCalculationType(ctx, params)

	if err != nil {
		return tables_model.DcCalculationType{}, fmt.Errorf("%w: dc.calculation_type: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateCalculationTypeById обновляет активную строку dc.calculation_type и возвращает её целиком.
//
// Запрос фильтрует по is_deleted = false, поэтому попытка обновить удалённую
// или несуществующую запись даёт sql.ErrNoRows — переводим его в ErrNotFound,
// чтобы api-слой ответил NotFound, а не Internal.
func (s *TablesService) UpdateCalculationTypeById(ctx context.Context, params tables_model.UpdateCalculationTypeByIdParams) (tables_model.DcCalculationType, error) {
	row, err := s.TablesRepository.UpdateCalculationTypeById(ctx, params)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tables_model.DcCalculationType{}, fmt.Errorf("dc.calculation_type id = %d: %w", params.ID, customerrors.ErrNotFound)
		}

		return tables_model.DcCalculationType{}, fmt.Errorf("%w: dc.calculation_type id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return row, nil
}

// DeleteCalculationTypeById мягко удаляет строку dc.calculation_type.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *TablesService) DeleteCalculationTypeById(ctx context.Context, id int64) error {
	if _, err := s.GetCalculationTypeById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.TablesRepository.DeleteCalculationTypeById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.calculation_type id = %d: %w", customerrors.ErrDelete, id, err)
	}

	return nil
}

// UndeleteCalculationTypeById восстанавливает мягко удалённую строку dc.calculation_type.
// Существование удалённой записи проверяется заранее по той же причине,
// что и в DeleteCalculationTypeById.
func (s *TablesService) UndeleteCalculationTypeById(ctx context.Context, id int64) error {
	if _, err := s.GetDeletedCalculationTypeById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrUndelete, err)
	}

	if err := s.TablesRepository.UndeleteCalculationTypeById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.calculation_type id = %d: %w", customerrors.ErrUndelete, id, err)
	}

	return nil
}
