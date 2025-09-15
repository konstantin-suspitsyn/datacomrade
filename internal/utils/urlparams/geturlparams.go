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
