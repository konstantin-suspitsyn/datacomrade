package paginationmodel

import "errors"

var ErrNegativaPageSize = errors.New("Eror. Negative page size")
var ErrCurrentPageIsBiggerThanTotal = errors.New("Eror. Current page is bigger than total")

