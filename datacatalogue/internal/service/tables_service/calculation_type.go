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

// GetCalculationTypes возвращает строки dc.calculation_type.
func (s *TablesService) GetCalculationTypes(ctx context.Context, order string) ([]tables_model.DcCalculationType, error) {
	rows, err := s.TablesRepository.GetCalculationTypes(ctx, order)

	if err != nil {
		return nil, fmt.Errorf("GetCalculationTypes: %w", err)
	}

	return rows, nil
}
