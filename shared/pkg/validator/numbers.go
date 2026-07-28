package validator

import (
	"fmt"
	"math"
	"strings"
)

// Number объединяет все числовые типы.
//
// Универсальные проверки сделаны пакетными функциями, а не методами:
// методы в Go не могут иметь собственных параметров типа. Для типов, которые
// реально приходят из sqlc и protobuf (int64, int32, float64), ниже есть
// методы-обёртки.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// NumberRequired проверяет, что число не нулевое.
func NumberRequired[T Number](v *Validator, field string, value T) bool {
	if value == 0 {
		v.AddError(field, MsgNumberRequired)
		return false
	}

	return true
}

// NumberPositive проверяет, что число строго больше нуля.
func NumberPositive[T Number](v *Validator, field string, value T) bool {
	if value <= 0 {
		v.AddError(field, fmt.Sprintf(MsgNumberPositive, value))
		return false
	}

	return true
}

// NumberNegative проверяет, что число строго меньше нуля.
func NumberNegative[T Number](v *Validator, field string, value T) bool {
	if value >= 0 {
		v.AddError(field, fmt.Sprintf(MsgNumberNegative, value))
		return false
	}

	return true
}

// NumberMin проверяет нижнюю границу включительно.
func NumberMin[T Number](v *Validator, field string, value, min T) bool {
	if value < min {
		v.AddError(field, fmt.Sprintf(MsgNumberMin, min, value))
		return false
	}

	return true
}

// NumberMax проверяет верхнюю границу включительно.
func NumberMax[T Number](v *Validator, field string, value, max T) bool {
	if value > max {
		v.AddError(field, fmt.Sprintf(MsgNumberMax, max, value))
		return false
	}

	return true
}

// NumberBetween проверяет, что число попадает в диапазон [min, max] включительно.
func NumberBetween[T Number](v *Validator, field string, value, min, max T) bool {
	if value < min || value > max {
		v.AddError(field, fmt.Sprintf(MsgNumberBetween, min, max, value))
		return false
	}

	return true
}

// NumberIn проверяет, что число входит в список допустимых.
func NumberIn[T Number](v *Validator, field string, value T, allowed ...T) bool {
	for _, item := range allowed {
		if value == item {
			return true
		}
	}

	names := make([]string, 0, len(allowed))
	for _, item := range allowed {
		names = append(names, fmt.Sprintf("%v", item))
	}

	v.AddError(field, fmt.Sprintf(MsgNumberIn, value, strings.Join(names, ", ")))

	return false
}

// Int64Required проверяет, что int64 не нулевой.
func (v *Validator) Int64Required(field string, value int64) bool {
	return NumberRequired(v, field, value)
}

// Int64ID проверяет идентификатор (первичный или внешний ключ int8): должен быть > 0.
func (v *Validator) Int64ID(field string, value int64) bool {
	return NumberPositive(v, field, value)
}

// Int64Positive проверяет, что int64 строго больше нуля.
func (v *Validator) Int64Positive(field string, value int64) bool {
	return NumberPositive(v, field, value)
}

// Int64Min проверяет нижнюю границу int64 включительно.
func (v *Validator) Int64Min(field string, value, min int64) bool {
	return NumberMin(v, field, value, min)
}

// Int64Max проверяет верхнюю границу int64 включительно.
func (v *Validator) Int64Max(field string, value, max int64) bool {
	return NumberMax(v, field, value, max)
}

// Int64Between проверяет диапазон int64 включительно.
func (v *Validator) Int64Between(field string, value, min, max int64) bool {
	return NumberBetween(v, field, value, min, max)
}

// Int64In проверяет, что int64 входит в список допустимых.
func (v *Validator) Int64In(field string, value int64, allowed ...int64) bool {
	return NumberIn(v, field, value, allowed...)
}

// Int64FitsInt32 проверяет, что int64 из protobuf влезает в колонку int4.
func (v *Validator) Int64FitsInt32(field string, value int64) bool {
	if value < math.MinInt32 || value > math.MaxInt32 {
		v.AddError(field, fmt.Sprintf(MsgNumberInt32, value, int64(math.MinInt32), int64(math.MaxInt32)))
		return false
	}

	return true
}

// Int32Required проверяет, что int32 не нулевой.
func (v *Validator) Int32Required(field string, value int32) bool {
	return NumberRequired(v, field, value)
}

// Int32Positive проверяет, что int32 строго больше нуля.
func (v *Validator) Int32Positive(field string, value int32) bool {
	return NumberPositive(v, field, value)
}

// Int32Min проверяет нижнюю границу int32 включительно.
func (v *Validator) Int32Min(field string, value, min int32) bool {
	return NumberMin(v, field, value, min)
}

// Int32Max проверяет верхнюю границу int32 включительно.
func (v *Validator) Int32Max(field string, value, max int32) bool {
	return NumberMax(v, field, value, max)
}

// Int32Between проверяет диапазон int32 включительно.
func (v *Validator) Int32Between(field string, value, min, max int32) bool {
	return NumberBetween(v, field, value, min, max)
}

// Int32In проверяет, что int32 входит в список допустимых.
func (v *Validator) Int32In(field string, value int32, allowed ...int32) bool {
	return NumberIn(v, field, value, allowed...)
}

// Float64Positive проверяет, что float64 строго больше нуля.
func (v *Validator) Float64Positive(field string, value float64) bool {
	if !v.Float64Finite(field, value) {
		return false
	}

	return NumberPositive(v, field, value)
}

// Float64Min проверяет нижнюю границу float64 включительно.
func (v *Validator) Float64Min(field string, value, min float64) bool {
	if !v.Float64Finite(field, value) {
		return false
	}

	return NumberMin(v, field, value, min)
}

// Float64Max проверяет верхнюю границу float64 включительно.
func (v *Validator) Float64Max(field string, value, max float64) bool {
	if !v.Float64Finite(field, value) {
		return false
	}

	return NumberMax(v, field, value, max)
}

// Float64Between проверяет диапазон float64 включительно.
func (v *Validator) Float64Between(field string, value, min, max float64) bool {
	if !v.Float64Finite(field, value) {
		return false
	}

	return NumberBetween(v, field, value, min, max)
}

// Float64Finite проверяет, что значение не NaN и не бесконечность.
// Такие значения не сохраняются в числовые колонки БД.
func (v *Validator) Float64Finite(field string, value float64) bool {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		v.AddError(field, MsgNumberFinite)
		return false
	}

	return true
}
