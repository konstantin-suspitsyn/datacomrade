package configs

import "time"

var SessionLength = 3 * 24 * time.Hour
var TokenDuration = 3 * 24 * time.Hour

const QueryTimeoutShort = 5 * time.Second
const QueryTimeoutLong = 60 * time.Second

const INDEX_PAGE_LINK = "/"

const USERS_V1 = "/v1/users"

const LOGOUT_LINK = "/logout"
const LOGIN_LINK = "/login"
const REFRESH_JWT_LINK = "/refresh"
const ACTIVATE_LINK = "/activate"
const RefreshJWTCookieName = "refresh_token"

const DOMAIN_V1 = "/v1/domain"
const DOMAIN_LINK = DOMAIN_V1

const GET_DOMAIN = "/"

const ROLES_V1 = "/v1/roles"
const ROLES_LINK = ROLES_V1

const GET_ROLES = "/"
