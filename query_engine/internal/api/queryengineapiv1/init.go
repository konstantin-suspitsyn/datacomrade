// Package queryengineapiv1 реализует gRPC-сервис query_engine.v1.QueryEngineService.
//
// Пока это заглушка: RPC не ходят в БД и не строят SQL, а возвращают ответ
// правильной структуры на фиксированных данных — чтобы shepherd и фронтенд
// могли интегрироваться с контрактом до появления настоящего конвейера
// (резолв алиасов через Metadata Service, сборка AST, рендер под диалект).
package queryengineapiv1

import (
	queryenginev1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/query_engine/v1"
)

// QueryEngineApiV1 — реализация query_engine.v1.QueryEngineService.
type QueryEngineApiV1 struct {
	queryenginev1.UnimplementedQueryEngineServiceServer
}

// New создаёт реализацию сервиса.
func New() *QueryEngineApiV1 {
	return &QueryEngineApiV1{}
}
