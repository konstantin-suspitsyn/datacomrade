package urlparams

import (
	"fmt"
	"net/http"
	"strconv"
)

func GetString(paramName string, r *http.Request) (string, error) {
	queryParams := r.URL.Query()

	stringParam := queryParams.Get(paramName)

	if stringParam == "" {
		return "", fmt.Errorf("ERROR. %w. Parameter name: %s", ErrNoParameterInUrl, paramName)
	}

	return stringParam, nil

}

func GetInt(paramName string, r *http.Request) (int64, error) {
	queryParams := r.URL.Query()

	stringParam := queryParams.Get(paramName)

	if stringParam == "" {
		return -1, fmt.Errorf("ERROR. %w. Parameter name: %s", ErrNoParameterInUrl, paramName)
	}
	n, err := strconv.ParseInt(stringParam, 10, 64)
	if err != nil {

		return -1, fmt.Errorf("ERROR converting string to int64. %w. Parameter name: %s", err, paramName)
	}

	return n, nil

}

