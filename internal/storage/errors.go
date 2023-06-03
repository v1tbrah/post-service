package storage

import "errors"

var (
	ErrPostNotFoundByID     = errors.New("post not found by id")
	ErrHashtagAlreadyExists = errors.New("hashtag already exists")
)
