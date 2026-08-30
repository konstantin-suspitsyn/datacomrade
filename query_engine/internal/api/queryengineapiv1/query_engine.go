package queryengineapiv1

import (
	"context"
	"encoding/json"
	"time"

	queryenginev1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/query_engine/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// stubSQL и stubParams — фиксированные «финальный SQL и его параметры».
// Настоящий query_engine соберёт их из формы и целевого диалекта; до тех пор
// это заглушка, одинаковая для всех вызовов.
var (
	stubSQL    = "SELECT 1 WHERE $1 = $1"
	stubParams = []string{"stub"}
)

// CreateQuery не сохраняет ничего в queries.queries, не строит SQL и не
// публикует в Kafka: собирает ответ правильной структуры на фиксированных
// данных. request заполняется JSON'ом присланной конфигурации (чтобы фронтенд
// видел ровно то, что отправил), query/params — фиксированной заглушкой, id и
// даты — постоянные. Ветку DomainChoice заглушка не отдаёт никогда.
func (a *QueryEngineApiV1) CreateQuery(_ context.Context, req *queryenginev1.CreateQueryRequest) (*queryenginev1.CreateQueryResponse, error) {
	if req.GetQuery() == nil {
		return nil, status.Error(codes.InvalidArgument, "query is required")
	}

	requestJSON, err := protojson.Marshal(req.GetQuery())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "marshal query to json: %v", err)
	}

	return &queryenginev1.CreateQueryResponse{
		Result: &queryenginev1.CreateQueryResponse_Saved{
			Saved: stubQuery(1, req.GetDatabaseId(), req.GetQuery().GetDomainId(), req.GetCreatedBy(), requestJSON),
		},
	}, nil
}

// GetQueryById возвращает фиксированную строку queries.queries с запрошенным id.
func (a *QueryEngineApiV1) GetQueryById(_ context.Context, req *queryenginev1.GetQueryByIdRequest) (*queryenginev1.GetQueryByIdResponse, error) {
	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id must be positive")
	}

	// Фиксированная конфигурация — та же, что отдаёт фронтенду пример из
	// документации формы OLAP.
	stubRequest, err := json.Marshal(stubOlapRequest())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "marshal stub request: %v", err)
	}

	return &queryenginev1.GetQueryByIdResponse{
		Query: stubQuery(req.GetId(), 1, 7, 1, stubRequest),
	}, nil
}

// stubQuery собирает строку queries.queries из фиксированных полей плюс
// переданных id/database_id/domain_id/created_by/request.
func stubQuery(id, databaseID, domainID, createdBy int64, request []byte) *queryenginev1.Query {
	ts := timestamppb.New(time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC))

	return &queryenginev1.Query{
		Id:         id,
		Query:      stubSQL,
		Params:     stubParams,
		Request:    string(request),
		DatabaseId: databaseID,
		DomainId:   domainID,
		CreatedAt:  ts,
		UpdatedAt:  ts,
		CreatedBy:  createdBy,
	}
}

// stubOlapRequest — фиксированная конфигурация сводной для GetQueryById.
func stubOlapRequest() map[string]any {
	return map[string]any{
		"version":   1,
		"domain_id": 7,
		"filters": []map[string]any{
			{"alias_id": 12, "operator": "between", "values": []string{"2026-01-01", "2026-03-31"}},
			{"alias_id": 5, "operator": "one_of", "values": []string{"RU", "KZ", "BY"}},
			{"alias_id": 8, "operator": "is_not_empty", "values": []string{}},
		},
		"rows": []map[string]any{
			{"alias_id": 5},
			{"alias_id": 9},
		},
		"measures": []map[string]any{
			{"kind": "field", "alias_id": 7, "aggregation": "sum", "label": "Выручка"},
			{"kind": "field", "alias_id": 3, "aggregation": "count_distinct", "label": "Заказы"},
			{"kind": "formula", "label": "Средний чек", "formula": "[Выручка] / [Заказы]"},
		},
	}
}
