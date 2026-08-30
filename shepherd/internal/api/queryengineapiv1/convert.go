package queryengineapiv1

import (
	queryenginev1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/query_engine/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shepherd/internal/api/openapiv1"
)

// filterOperatorToProto переводит строковый оператор REST-контракта в enum
// gRPC. Множества значений совпадают один в один (см. query_engine.yaml и
// shared/proto/query_engine/v1); неизвестное значение — нулевой enum, его
// отбивает валидация до вызова.
var filterOperatorToProto = map[openapiv1.OlapFilterOperator]queryenginev1.FilterOperator{
	openapiv1.Equals:             queryenginev1.FilterOperator_FILTER_OPERATOR_EQUALS,
	openapiv1.NotEquals:          queryenginev1.FilterOperator_FILTER_OPERATOR_NOT_EQUALS,
	openapiv1.GreaterThan:        queryenginev1.FilterOperator_FILTER_OPERATOR_GREATER_THAN,
	openapiv1.GreaterThanOrEqual: queryenginev1.FilterOperator_FILTER_OPERATOR_GREATER_THAN_OR_EQUAL,
	openapiv1.LessThan:           queryenginev1.FilterOperator_FILTER_OPERATOR_LESS_THAN,
	openapiv1.LessThanOrEqual:    queryenginev1.FilterOperator_FILTER_OPERATOR_LESS_THAN_OR_EQUAL,
	openapiv1.Between:            queryenginev1.FilterOperator_FILTER_OPERATOR_BETWEEN,
	openapiv1.OneOf:              queryenginev1.FilterOperator_FILTER_OPERATOR_ONE_OF,
	openapiv1.NotOneOf:           queryenginev1.FilterOperator_FILTER_OPERATOR_NOT_ONE_OF,
	openapiv1.Contains:           queryenginev1.FilterOperator_FILTER_OPERATOR_CONTAINS,
	openapiv1.StartsWith:         queryenginev1.FilterOperator_FILTER_OPERATOR_STARTS_WITH,
	openapiv1.IsEmpty:            queryenginev1.FilterOperator_FILTER_OPERATOR_IS_EMPTY,
	openapiv1.IsNotEmpty:         queryenginev1.FilterOperator_FILTER_OPERATOR_IS_NOT_EMPTY,
}

var aggregationToProto = map[openapiv1.OlapAggregation]queryenginev1.Aggregation{
	openapiv1.Sum:           queryenginev1.Aggregation_AGGREGATION_SUM,
	openapiv1.Avg:           queryenginev1.Aggregation_AGGREGATION_AVG,
	openapiv1.Min:           queryenginev1.Aggregation_AGGREGATION_MIN,
	openapiv1.Max:           queryenginev1.Aggregation_AGGREGATION_MAX,
	openapiv1.Count:         queryenginev1.Aggregation_AGGREGATION_COUNT,
	openapiv1.CountDistinct: queryenginev1.Aggregation_AGGREGATION_COUNT_DISTINCT,
}

var measureKindToProto = map[openapiv1.OlapMeasureKind]queryenginev1.MeasureKind{
	openapiv1.Field:   queryenginev1.MeasureKind_MEASURE_KIND_FIELD,
	openapiv1.Formula: queryenginev1.MeasureKind_MEASURE_KIND_FORMULA,
}

var sortDirectionToProto = map[openapiv1.OlapSortDirection]queryenginev1.SortDirection{
	openapiv1.Asc:  queryenginev1.SortDirection_SORT_DIRECTION_ASC,
	openapiv1.Desc: queryenginev1.SortDirection_SORT_DIRECTION_DESC,
}

// toOlapQuery собирает proto-конфигурацию сводной из тела REST-запроса.
// database_id остаётся соседним полем CreateQueryRequest, domain_id уезжает
// внутрь OlapQuery (и, значит, в queries.queries.request).
func toOlapQuery(body openapiv1.OlapQueryCreate) *queryenginev1.OlapQuery {
	filters := make([]*queryenginev1.OlapFilter, 0, len(body.Filters))
	for _, f := range body.Filters {
		values := f.Values
		if values == nil {
			values = []string{}
		}
		filters = append(filters, &queryenginev1.OlapFilter{
			AliasId:  f.AliasId,
			Operator: filterOperatorToProto[f.Operator],
			Values:   values,
		})
	}

	rows := make([]*queryenginev1.OlapRow, 0, len(body.Rows))
	for _, r := range body.Rows {
		rows = append(rows, &queryenginev1.OlapRow{AliasId: r.AliasId})
	}

	var columns []*queryenginev1.OlapColumn
	if body.Columns != nil {
		columns = make([]*queryenginev1.OlapColumn, 0, len(*body.Columns))
		for _, c := range *body.Columns {
			columns = append(columns, &queryenginev1.OlapColumn{AliasId: c.AliasId})
		}
	}

	measures := make([]*queryenginev1.OlapMeasure, 0, len(body.Measures))
	for _, m := range body.Measures {
		measure := &queryenginev1.OlapMeasure{
			Kind:  measureKindToProto[m.Kind],
			Label: m.Label,
		}
		if m.AliasId != nil {
			measure.AliasId = *m.AliasId
		}
		if m.Aggregation != nil {
			measure.Aggregation = aggregationToProto[*m.Aggregation]
		}
		if m.Formula != nil {
			measure.Formula = *m.Formula
		}
		measures = append(measures, measure)
	}

	var orderBy []*queryenginev1.OlapSort
	if body.OrderBy != nil {
		orderBy = make([]*queryenginev1.OlapSort, 0, len(*body.OrderBy))
		for _, s := range *body.OrderBy {
			sort := &queryenginev1.OlapSort{Direction: sortDirectionToProto[s.Direction]}
			if s.AliasId != nil {
				sort.AliasId = *s.AliasId
			}
			if s.MeasureLabel != nil {
				sort.MeasureLabel = *s.MeasureLabel
			}
			orderBy = append(orderBy, sort)
		}
	}

	query := &queryenginev1.OlapQuery{
		Version:  body.Version,
		Filters:  filters,
		Rows:     rows,
		Columns:  columns,
		Measures: measures,
		OrderBy:  orderBy,
	}
	if body.DomainId != nil {
		query.DomainId = body.DomainId
	}
	if body.Limit != nil {
		query.Limit = body.Limit
	}

	return query
}

// toOlapQueryRecord переводит строку queries.queries из gRPC-ответа в тело
// REST-ответа.
func toOlapQueryRecord(q *queryenginev1.Query) openapiv1.OlapQueryRecord {
	params := q.GetParams()
	if params == nil {
		params = []string{}
	}

	return openapiv1.OlapQueryRecord{
		Id:         q.GetId(),
		Query:      q.GetQuery(),
		Params:     params,
		Request:    q.GetRequest(),
		DatabaseId: q.GetDatabaseId(),
		DomainId:   q.GetDomainId(),
		CreatedAt:  q.GetCreatedAt().AsTime(),
		UpdatedAt:  q.GetUpdatedAt().AsTime(),
		CreatedBy:  q.GetCreatedBy(),
	}
}

// toDomainChoice переводит список доменов на выбор из gRPC-ответа в тело
// REST-ответа 409.
func toDomainChoice(dc *queryenginev1.DomainChoice) openapiv1.OlapDomainChoice {
	domains := make([]openapiv1.OlapDomainOption, 0, len(dc.GetDomains()))
	for _, d := range dc.GetDomains() {
		domains = append(domains, openapiv1.OlapDomainOption{Id: d.GetId(), Name: d.GetName()})
	}
	return openapiv1.OlapDomainChoice{Domains: domains}
}

// unknownFilterOperators возвращает операторы фильтров, которых нет в контракте.
func unknownFilterOperators(body openapiv1.OlapQueryCreate) []string {
	var bad []string
	for _, f := range body.Filters {
		if _, ok := filterOperatorToProto[f.Operator]; !ok {
			bad = append(bad, string(f.Operator))
		}
	}
	return bad
}

// invalidMeasures проверяет соответствие полей показателя его виду.
func invalidMeasures(body openapiv1.OlapQueryCreate) bool {
	for _, m := range body.Measures {
		switch m.Kind {
		case openapiv1.Field:
			if m.AliasId == nil || m.Aggregation == nil {
				return true
			}
			if _, ok := aggregationToProto[*m.Aggregation]; !ok {
				return true
			}
		case openapiv1.Formula:
			if m.Formula == nil || *m.Formula == "" {
				return true
			}
		default:
			return true
		}
	}
	return false
}
