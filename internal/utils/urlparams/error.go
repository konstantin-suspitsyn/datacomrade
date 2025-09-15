package urlparams

import "errors"

var ErrNoParameterInUrl = errors.New("No parameter in url query.")
var ErrNegativaParameterInUrl = errors.New("Negative parameter in url query")
