package validator

// Optional проверяет необязательное поле, переданное указателем:
// nil проходит проверку, ненулевой указатель проверяется функцией check.
//
//	validator.Optional(v, "description", in.Description, func(v *validator.Validator, field, value string) bool {
//	    return v.StringMaxLen(field, value, 2000)
//	})
func Optional[T any](v *Validator, field string, value *T, check func(v *Validator, field string, value T) bool) bool {
	if value == nil {
		return true
	}

	return check(v, field, *value)
}

// OptionalStringVarchar проверяет необязательную колонку varchar(max):
// nil проходит, непустой указатель проверяется на длину.
func (v *Validator) OptionalStringVarchar(field string, value *string, max int) bool {
	if value == nil {
		return true
	}

	return v.StringMaxLen(field, *value, max)
}

// OptionalInt64ID проверяет необязательный внешний ключ:
// nil проходит, заполненный должен быть больше нуля.
func (v *Validator) OptionalInt64ID(field string, value *int64) bool {
	if value == nil {
		return true
	}

	return v.Int64ID(field, *value)
}
