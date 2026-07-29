package validator

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// EmailRX — регулярное выражение для проверки email (HTML5 spec).
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// IdentifierRX — имя таблицы, схемы или колонки: латиница, цифры и подчёркивание,
// первый символ не цифра.
var IdentifierRX = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// UUIDRX — UUID в каноническом виде 8-4-4-4-12, регистр букв любой.
// Сокращённые формы (без дефисов, в фигурных скобках, с префиксом urn:uuid:)
// намеренно не принимаются: колонка uuid хранит одно значение, а вариантов
// его записи быть не должно.
var UUIDRX = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// StringRequired проверяет, что строка не пустая и состоит не из одних пробелов.
func (v *Validator) StringRequired(field, value string) bool {
	if strings.TrimSpace(value) == "" {
		v.AddError(field, MsgRequired)
		return false
	}

	return true
}

// StringMaxLen проверяет верхнюю границу длины строки.
// Длина считается в символах (рунах), как в varchar(n) в PostgreSQL,
// а не в байтах.
func (v *Validator) StringMaxLen(field, value string, max int) bool {
	length := utf8.RuneCountInString(value)

	if length > max {
		v.AddError(field, fmt.Sprintf(MsgMaxLen, max, length))
		return false
	}

	return true
}

// StringMinLen проверяет нижнюю границу длины строки в символах.
func (v *Validator) StringMinLen(field, value string, min int) bool {
	length := utf8.RuneCountInString(value)

	if length < min {
		v.AddError(field, fmt.Sprintf(MsgMinLen, min, length))
		return false
	}

	return true
}

// StringLenBetween проверяет, что длина строки в символах попадает в диапазон
// [min, max] включительно.
func (v *Validator) StringLenBetween(field, value string, min, max int) bool {
	return v.StringMinLen(field, value, min) && v.StringMaxLen(field, value, max)
}

// StringExactLen проверяет точную длину строки в символах.
func (v *Validator) StringExactLen(field, value string, length int) bool {
	actual := utf8.RuneCountInString(value)

	if actual != length {
		v.AddError(field, fmt.Sprintf(MsgExactLen, length, actual))
		return false
	}

	return true
}

// StringVarchar — проверка для колонки varchar(max) NOT NULL:
// строка обязательна и не длиннее max символов.
func (v *Validator) StringVarchar(field, value string, max int) bool {
	if !v.StringRequired(field, value) {
		return false
	}

	return v.StringMaxLen(field, value, max)
}

// StringOptionalVarchar — проверка для колонки varchar(max), которую разрешено
// оставить пустой: пустая строка проходит, непустая проверяется на длину.
func (v *Validator) StringOptionalVarchar(field, value string, max int) bool {
	if value == "" {
		return true
	}

	return v.StringMaxLen(field, value, max)
}

// StringIn проверяет, что значение входит в список допустимых.
func (v *Validator) StringIn(field, value string, allowed ...string) bool {
	for _, item := range allowed {
		if value == item {
			return true
		}
	}

	v.AddError(field, fmt.Sprintf(MsgStringIn, value, strings.Join(allowed, ", ")))

	return false
}

// StringMatches проверяет строку регулярным выражением.
// hint — понятное человеку описание ожидаемого формата.
func (v *Validator) StringMatches(field, value string, rx *regexp.Regexp, hint string) bool {
	if rx == nil || !rx.MatchString(value) {
		v.AddError(field, fmt.Sprintf(MsgStringPattern, value, hint))
		return false
	}

	return true
}

// StringEmail проверяет, что строка похожа на email.
func (v *Validator) StringEmail(field, value string) bool {
	if !EmailRX.MatchString(value) {
		v.AddError(field, fmt.Sprintf(MsgEmail, value))
		return false
	}

	return true
}

// StringIdentifier проверяет, что строка пригодна как имя объекта БД
// (схема, таблица, колонка).
func (v *Validator) StringIdentifier(field, value string) bool {
	if !IdentifierRX.MatchString(value) {
		v.AddError(field, fmt.Sprintf(MsgIdentifier, value))
		return false
	}

	return true
}

// StringUUID — проверка для колонки uuid NOT NULL: значение обязательно
// и записано в каноническом виде. Проверка идёт до конвертера, поэтому
// разбор строки в uuid.UUID дальше по стеку уже не может не удаться.
func (v *Validator) StringUUID(field, value string) bool {
	if !v.StringRequired(field, value) {
		return false
	}

	if !UUIDRX.MatchString(value) {
		v.AddError(field, fmt.Sprintf(MsgUUID, value))
		return false
	}

	return true
}

// StringNoSpaces проверяет, что в строке нет пробельных символов.
func (v *Validator) StringNoSpaces(field, value string) bool {
	if strings.ContainsAny(value, " \t\n\r\v\f") {
		v.AddError(field, MsgNoSpaces)
		return false
	}

	return true
}

// StringTrimmed проверяет, что строка не начинается и не заканчивается пробелом.
func (v *Validator) StringTrimmed(field, value string) bool {
	if value != strings.TrimSpace(value) {
		v.AddError(field, MsgTrimmed)
		return false
	}

	return true
}
