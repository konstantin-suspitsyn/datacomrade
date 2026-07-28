package validator

import (
	"fmt"
	"time"
)

// TimeRequired проверяет, что дата заполнена (не нулевая).
func (v *Validator) TimeRequired(field string, value time.Time) bool {
	if value.IsZero() {
		v.AddError(field, MsgTimeRequired)
		return false
	}

	return true
}

// TimeAfter проверяет, что дата строго позже границы.
func (v *Validator) TimeAfter(field string, value, after time.Time) bool {
	if !value.After(after) {
		v.AddError(field, fmt.Sprintf(MsgTimeAfter, after.Format(TimeLayout), value.Format(TimeLayout)))
		return false
	}

	return true
}

// TimeBefore проверяет, что дата строго раньше границы.
func (v *Validator) TimeBefore(field string, value, before time.Time) bool {
	if !value.Before(before) {
		v.AddError(field, fmt.Sprintf(MsgTimeBefore, before.Format(TimeLayout), value.Format(TimeLayout)))
		return false
	}

	return true
}

// TimeBetween проверяет, что дата попадает в диапазон [from, to] включительно.
func (v *Validator) TimeBetween(field string, value, from, to time.Time) bool {
	if value.Before(from) || value.After(to) {
		v.AddError(field, fmt.Sprintf(MsgTimeBetween, from.Format(TimeLayout), to.Format(TimeLayout), value.Format(TimeLayout)))
		return false
	}

	return true
}

// TimeNotInFuture проверяет, что дата не позже текущего момента.
func (v *Validator) TimeNotInFuture(field string, value time.Time) bool {
	if value.After(time.Now()) {
		v.AddError(field, MsgTimeFuture)
		return false
	}

	return true
}

// TimeNotInPast проверяет, что дата не раньше текущего момента.
func (v *Validator) TimeNotInPast(field string, value time.Time) bool {
	if value.Before(time.Now()) {
		v.AddError(field, MsgTimePast)
		return false
	}

	return true
}
