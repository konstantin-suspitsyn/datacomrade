// Package customerrors содержит сторожевые ошибки слоя сервисов.
//
// Ошибки оборачиваются через %w, поэтому вызывающий код проверяет их
// через errors.Is, а api-слой по ним выбирает gRPC-код ответа.
package customerrors

import "errors"

// ErrNotFound — запись с таким id не найдена (или уже удалена, если
// запрос шёл по активным записям).
var ErrNotFound = errors.New("entity not found")

// ErrDelete — мягкое удаление не удалось.
var ErrDelete = errors.New("cannot delete entity")

// ErrUndelete — восстановление удалённой записи не удалось.
var ErrUndelete = errors.New("cannot undelete entity")

// ErrCreate — вставка записи не удалась.
var ErrCreate = errors.New("cannot create entity")

// ErrUpdate — обновление записи не удалось.
var ErrUpdate = errors.New("cannot update entity")
