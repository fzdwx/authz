package authz

import "errors"

var (
	ErrStoreValueNotFound = errors.New("value not found")
	ErrNoToken            = errors.New("token not found in context")
)
