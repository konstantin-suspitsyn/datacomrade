package validator

// BoolTrue проверяет, что флаг взведён.
func (v *Validator) BoolTrue(field string, value bool) bool {
	if !value {
		v.AddError(field, MsgBoolTrue)
		return false
	}

	return true
}

// BoolFalse проверяет, что флаг снят.
func (v *Validator) BoolFalse(field string, value bool) bool {
	if value {
		v.AddError(field, MsgBoolFalse)
		return false
	}

	return true
}
