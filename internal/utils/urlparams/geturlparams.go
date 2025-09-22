package urlparams

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/konstantin-suspitsyn/datacomrade/configs"
)

type Pager struct {
	CurrentPage int64
	PageSize    int //Items per page
	Sort        string
}

// GetString retrieves a string value from the query parameters of an HTTP request.
//
// If the parameter is missing (i.e., has no value), it returns an error wrapped with ErrNoParameterInUrl,
// including the parameter name in the error message for clarity.
func GetString(paramName string, r *http.Request) (string, error) {
	queryParams := r.URL.Query()

	stringParam := queryParams.Get(paramName)

	if stringParam == "" {
		return "", fmt.Errorf("ERROR. %w. Parameter name: %s", ErrNoParameterInUrl, paramName)
	}

	return stringParam, nil

}

// GetInt retrieves an integer value from the query parameters of an HTTP request.
//
// If the parameter is missing (i.e., has no value), it returns -1 and an error wrapped with ErrNoParameterInUrl,
// including the parameter name in the error message for clarity.
// If the parameter cannot be parsed as a valid 64-bit integer, it returns -1 and the underlying error,
// wrapped with a message indicating the conversion failure and the parameter name.
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

// GetPager parses pagination-related query parameters from an HTTP request and constructs a Pager struct.
//
// Behavior:
//   - Uses `configs.PAGE_PARAM` to retrieve the current page number. Defaults to 1 if the parameter is missing.
//   - Uses `configs.ITEMS_PER_PAGE_PARAM` to retrieve the number of items per page. Defaults to 0 if the parameter is missing.
//   - Returns an error if either parameter is negative.
//   - Wraps underlying errors from GetString/GetInt with context about the parameter name.
func GetPager(r *http.Request) (*Pager, error) {
	pageNo, err := GetInt(configs.PAGE_PARAM, r)
	if err != nil {
		switch {
		case errors.Is(err, ErrNoParameterInUrl):
			pageNo = 1
		default:
			return nil, fmt.Errorf("Getting parameter form request. %w", err)
		}
	}

	if pageNo < 0 {
		return nil, fmt.Errorf("%w. Parameter: %s", ErrNegativaParameterInUrl, configs.PAGE_PARAM)
	}

	itemsPerPage, err := GetInt(configs.ITEMS_PER_PAGE_PARAM, r)
	if err != nil {
		switch {
		case errors.Is(err, ErrNoParameterInUrl):
			itemsPerPage = 0
		default:
			return nil, fmt.Errorf("Getting parameter form request. %w", err)
		}
	}

	if itemsPerPage < 0 {
		return nil, fmt.Errorf("%w. Parameter: %s", ErrNegativaParameterInUrl, configs.ITEMS_PER_PAGE_PARAM)
	}

	return &Pager{
		CurrentPage: pageNo,
		PageSize:    int(itemsPerPage),
	}, nil
}
