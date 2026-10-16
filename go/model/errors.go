package model

import "errors"

var (
	ErrInvalidFilter = errors.New("invalid filter parameter")
	ErrInvalidGame   = errors.New("invalid game")
	ErrGameConflict  = errors.New("game already exists")
)
