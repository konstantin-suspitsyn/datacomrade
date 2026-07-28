package validator

import "fmt"

// SliceRequired проверяет, что список не пустой.
func SliceRequired[T any](v *Validator, field string, values []T) bool {
	if len(values) == 0 {
		v.AddError(field, MsgSliceRequired)
		return false
	}

	return true
}

// SliceMinLen проверяет минимальное количество элементов.
func SliceMinLen[T any](v *Validator, field string, values []T, min int) bool {
	if len(values) < min {
		v.AddError(field, fmt.Sprintf(MsgSliceMinLen, min, len(values)))
		return false
	}

	return true
}

// SliceMaxLen проверяет максимальное количество элементов.
func SliceMaxLen[T any](v *Validator, field string, values []T, max int) bool {
	if len(values) > max {
		v.AddError(field, fmt.Sprintf(MsgSliceMaxLen, max, len(values)))
		return false
	}

	return true
}

// SliceLenBetween проверяет, что количество элементов попадает в диапазон
// [min, max] включительно.
func SliceLenBetween[T any](v *Validator, field string, values []T, min, max int) bool {
	return SliceMinLen(v, field, values, min) && SliceMaxLen(v, field, values, max)
}

// SliceUnique проверяет, что элементы списка не повторяются.
func SliceUnique[T comparable](v *Validator, field string, values []T) bool {
	seen := make(map[T]struct{}, len(values))

	for _, value := range values {
		if _, ok := seen[value]; ok {
			v.AddError(field, fmt.Sprintf(MsgSliceUnique, value))
			return false
		}

		seen[value] = struct{}{}
	}

	return true
}

// SliceEach прогоняет проверку по каждому элементу списка. К имени поля
// добавляется индекс элемента: "columns" -> "columns.0".
func SliceEach[T any](v *Validator, field string, values []T, check func(v *Validator, field string, value T) bool) bool {
	ok := true

	for i, value := range values {
		if !check(v, fmt.Sprintf("%s.%d", field, i), value) {
			ok = false
		}
	}

	return ok
}
