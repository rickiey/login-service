package helper

import "errors"

var (
	ErrInternalServerError = errors.New("internal Server Error")
	ErrNotFound            = errors.New("your requested Item is not found")
	ErrBadParamInput       = errors.New("given Param is not valid")
	ErrAuthenticateError   = errors.New("invalid token")
	ErrUserForbidden       = errors.New("该账号已被禁用,请联系上级或管理员!")
	ErrMethodNotAllowed    = errors.New("invalid request method")
)
