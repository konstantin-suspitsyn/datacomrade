// Package validator предоставляет объект для проверки типов и размеров полей
// при создании и обновлении сущностей.
//
// Объект накапливает ошибки по каждому полю и на выходе отдаёт две вещи:
//
//  1. map[string][]string — название поля и понятный список того, что сломалось;
//  2. флаг наличия/отсутствия ошибок (Valid / HasErrors).
//
// Каждая проверка — отдельная функция. Проверки сгруппированы по типу поля:
// строки (strings.go), числа (numbers.go), время (time.go),
// булевы значения (bool.go), слайсы (slices.go), необязательные поля (optional.go).
//
// Пример:
//
//	v := validator.New()
//	v.StringVarchar("name", in.Name, 128)         // varchar(128) NOT NULL
//	v.StringVarchar("description", in.Desc, 2000) // varchar(2000) NOT NULL
//	v.Int64ID("schema_id", in.SchemaId)
//	v.Int64ID("domain_id", in.DomainId)
//
//	if !v.Valid() {
//	    return nil, status.Error(codes.InvalidArgument, v.Err().Eror())
//	}
//
// Каждая проверка возвращает bool, поэтому проверки можно связывать:
//
//	if v.StringRequired("name", in.Name) {
//	    v.StringMaxLen("name", in.Name, 128)
//	}
package validator

import (
	"fmt"
	"sort"
	"strings"
)

// Validator накапливает ошибки валидации по полям сущности.
//
// Готовый к работе объект создаётся через New. Нулевое значение
// (var v validator.Validator) тоже пригодно к использованию — внутренняя
// мапа создаётся лениво.
//
// Validator не предназначен для конкурентного использования: один объект
// живёт в рамках одной проверки одной сущности.
type Validator struct {
	errors map[string][]string
}

// New создаёт пустой валидатор.
func New() *Validator {
	return &Validator{errors: make(map[string][]string)}
}

// AddError добавляет ошибку к полю. Одинаковые сообщения для одного поля
// не дублируются.
func (v *Validator) AddError(field, message string) {
	if v.errors == nil {
		v.errors = make(map[string][]string)
	}

	for _, existing := range v.errors[field] {
		if existing == message {
			return
		}
	}

	v.errors[field] = append(v.errors[field], message)
}

// Check добавляет ошибку к полю, если условие ok не выполнено.
// Нужен для проверок, которых нет среди готовых функций.
func (v *Validator) Check(ok bool, field, message string) bool {
	if !ok {
		v.AddError(field, message)
		return false
	}

	return true
}

// Valid сообщает, что ошибок нет.
func (v *Validator) Valid() bool {
	return len(v.errors) == 0
}

// HasErrors сообщает, что ошибки есть. Обратная сторона Valid.
func (v *Validator) HasErrors() bool {
	return len(v.errors) > 0
}

// Errors возвращает копию мапы «поле -> список ошибок».
// Изменение результата не влияет на валидатор.
func (v *Validator) Errors() map[string][]string {
	out := make(map[string][]string, len(v.errors))

	for field, messages := range v.errors {
		out[field] = append([]string(nil), messages...)
	}

	return out
}

// FieldErrors возвращает список ошибок конкретного поля.
// Для поля без ошибок возвращается nil.
func (v *Validator) FieldErrors(field string) []string {
	if len(v.errors[field]) == 0 {
		return nil
	}

	return append([]string(nil), v.errors[field]...)
}

// FieldValid сообщает, что у конкретного поля нет ошибок.
func (v *Validator) FieldValid(field string) bool {
	return len(v.errors[field]) == 0
}

// Fields возвращает отсортированный список полей, у которых есть ошибки.
func (v *Validator) Fields() []string {
	fields := make([]string, 0, len(v.errors))

	for field := range v.errors {
		fields = append(fields, field)
	}

	sort.Strings(fields)

	return fields
}

// Count возвращает общее количество ошибок по всем полям.
func (v *Validator) Count() int {
	count := 0

	for _, messages := range v.errors {
		count += len(messages)
	}

	return count
}

// Merge переносит ошибки другого валидатора в текущий.
// Нужен для вложенных сущностей: непустой prefix добавляется к имени поля
// через точку (например, "columns.0" + "name" -> "columns.0.name").
func (v *Validator) Merge(prefix string, other *Validator) {
	if other == nil {
		return
	}

	for field, messages := range other.errors {
		name := field
		if prefix != "" {
			name = prefix + "." + field
		}

		for _, message := range messages {
			v.AddError(name, message)
		}
	}
}

// Reset очищает валидатор, чтобы переиспользовать объект.
func (v *Validator) Reset() {
	v.errors = make(map[string][]string)
}

// Err возвращает *ValidationError, если ошибки есть, и nil, если их нет.
// Удобно отдавать наверх из сервисного слоя:
//
//	if err := v.Err(); err != nil {
//	    return err
//	}
func (v *Validator) Err() error {
	if v.Valid() {
		return nil
	}

	return &ValidationError{Errors: v.Errors()}
}

// ValidationError — ошибка валидации со списком проблем по каждому полю.
type ValidationError struct {
	Errors map[string][]string
}

// Error собирает все сообщения в одну детерминированную строку.
func (e *ValidationError) Error() string {
	if e == nil || len(e.Errors) == 0 {
		return MsgValidationFailed
	}

	fields := make([]string, 0, len(e.Errors))
	for field := range e.Errors {
		fields = append(fields, field)
	}

	sort.Strings(fields)

	parts := make([]string, 0, len(fields))
	for _, field := range fields {
		parts = append(parts, fmt.Sprintf("%s: %s", field, strings.Join(e.Errors[field], ", ")))
	}

	return MsgValidationFailed + ": " + strings.Join(parts, "; ")
}
