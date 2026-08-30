package game

import "errors"

var (
	ErrInvalidDifficulty = errors.New("ERR_INVALID_DIFFICULTY")
	ErrNotImplemented    = errors.New("ERR_NOT_IMPLEMENTED")
	ErrInvalidTime       = errors.New("ERR_INVALID_TIME")
	ErrInvalidBoardSize  = errors.New("ERR_INVALID_BOARD_SIZE")
	ErrNoActiveMatch     = errors.New("ERR_NO_ACTIVE_MATCH")
	ErrInvalidCell       = errors.New("ERR_INVALID_CELL")
)
