package validator

import (
	"regexp"
	"slices"
)

type Validator struct {
	Errors map[string]string
}

var (
	EmailRX = regexp.MustCompile("(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\\])")
)

// Creates new Validator with no Errors
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Checks if Validator has no errors
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// Add error to validator
func (v *Validator) AddError(key, message string) {

	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Adds error only when !ok
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func (v *Validator) In(value string, list ...string) bool {

	if slices.Contains(list, value) {
		return true
	}
	return false

}

func (v *Validator) Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func (v *Validator) Unique(values []string) bool {
	uniqueValues := make(map[string]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}
	return len(uniqueValues) == len(values)
}
