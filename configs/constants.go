package configs

import "time"

var SessionLength = 3 * 24 * time.Hour
var TokenDuration = 3 * 24 * time.Hour

const QueryTimeoutShort = 5 * time.Second
const QueryTimeoutLong = 60 * time.Second
