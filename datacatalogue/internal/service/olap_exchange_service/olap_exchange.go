package olapexchangeservice

import (
	"context"
	"fmt"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/olapexchange"
)

// GetOlapAliases оборачивает sqlc-запрос olap_exchange.GetOlapAliases.
func (s *OlapExchangeService) GetOlapAliases(ctx context.Context, params olapexchange.GetOlapAliasesParams) ([]olapexchange.GetOlapAliasesRow, error) {
	rows, err := s.OlapExchangeRepository.GetOlapAliases(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("GetOlapAliases: %w", err)
	}

	return rows, nil
}
