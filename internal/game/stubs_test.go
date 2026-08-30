package game

import (
	"errors"
	"testing"
)

func TestChordCellStillReturnsNotImplemented(t *testing.T) {
	board, err := ChordCell(0, 0, 0)
	if board != nil {
		t.Errorf("ChordCell returned a board: %+v", board)
	}
	if !errors.Is(err, ErrNotImplemented) {
		t.Errorf("ChordCell returned %v, want %v", err, ErrNotImplemented)
	}
}
