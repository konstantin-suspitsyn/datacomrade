// Package converter содержит общие примитивы преобразования между типами
// gRPC (protobuf) и типами репозитория (sqlc).
//
// Преобразования по сущностям лежат в подпакетах, повторяющих деление
// на sqlc-пакеты и gRPC-сервисы:
//
//	converter/tables            — dc.alias, dc.host, dc.table_cat и остальные 15 таблиц
//	converter/user              — dc.user
//	converter/userdomainroles   — dc.domain_roles и остальные 6 таблиц
//
// Функции конвертации — чистые: они не ходят в базу, не логируют и не
// возвращают ошибку. Проверка входных данных — задача пакета validation,
// он вызывается в api-слое до конвертации.
package converter

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TimeToProto переводит time.Time из строки таблицы в protobuf Timestamp.
func TimeToProto(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}

// NullTimeToProto переводит nullable-колонку времени в protobuf Timestamp.
// SQL NULL становится nil: у message-полей в proto3 presence есть
// по умолчанию, поэтому отдельный флаг не нужен.
func NullTimeToProto(t sql.NullTime) *timestamppb.Timestamp {
	if !t.Valid {
		return nil
	}

	return timestamppb.New(t.Time)
}

// ProtoToNullTime переводит protobuf Timestamp в nullable-колонку времени.
// nil становится SQL NULL.
func ProtoToNullTime(t *timestamppb.Timestamp) sql.NullTime {
	if t == nil {
		return sql.NullTime{}
	}

	return sql.NullTime{Time: t.AsTime(), Valid: true}
}

// NullStringToProto переводит nullable-колонку строки в необязательное поле proto.
// SQL NULL становится nil.
func NullStringToProto(s sql.NullString) *string {
	if !s.Valid {
		return nil
	}

	value := s.String

	return &value
}

// ProtoToNullString переводит необязательное поле proto в nullable-колонку строки.
// nil становится SQL NULL.
func ProtoToNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}

	return sql.NullString{String: *s, Valid: true}
}

// Int64ToProto переводит колонку id в необязательное поле proto — внутренний
// идентификатор не гарантирован клиентам API, только external_id.
func Int64ToProto(id int64) *int64 {
	return &id
}

// UUIDToProto переводит колонку uuid в строковое поле proto:
// в protobuf нет типа для UUID, значение передаётся канонической записью.
func UUIDToProto(id uuid.UUID) string {
	return id.String()
}

// ProtoToUUID разбирает строковое поле proto в колонку uuid.
//
// Функции конвертации не возвращают ошибку, поэтому нераспознанное значение
// становится uuid.Nil. До конвертера такая строка не доходит: формат проверяет
// validator.StringUUID в api-слое, а uuid.Nil не совпадёт ни с одной записью
// в базе, если проверку когда-нибудь обойдут.
func ProtoToUUID(s string) uuid.UUID {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}

	return id
}
