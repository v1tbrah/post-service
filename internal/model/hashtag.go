package model

import (
	"github.com/pkg/errors"
)

type Hashtag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Direction int32

const (
	First Direction = iota
	Next
	Prev
)

var ErrInvalidDirection = errors.New("invalid direction")
