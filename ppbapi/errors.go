package ppbapi

import "errors"

var ErrEmptyRequest = errors.New("empty request")

var (
	ErrEmptyName = errors.New("empty name")
)

var (
	ErrPostNotFoundByID    = errors.New("post not found by id")
	ErrHashtagNotFoundByID = errors.New("hashtag not found by id")
)
