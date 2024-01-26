package action

import "errors"

var (
	ErrAllRegionsNotAllowed = errors.New("all regions is not allowed")
	ErrRegionsRequired      = errors.New("regions is required")
	ErrTTLRequired          = errors.New("ttl is required")
	ErrTokenRequired        = errors.New("token is required")
)
