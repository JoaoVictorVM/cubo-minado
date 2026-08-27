package game

import (
	"errors"
	"testing"
)

func TestNewGameValidDifficulties(t *testing.T) {
	cases := map[Difficulty]int{
		DifficultyEasy:   5,
		DifficultyMedium: 7,
		DifficultyHard:   9,
	}

	for difficulty, wantSize := range cases {
		board, err := NewGame(difficulty)
		if err != nil {
			t.Fatalf("NewGame(%q) returned error: %v", difficulty, err)
		}
		if board.Difficulty != difficulty {
			t.Errorf("NewGame(%q).Difficulty = %q", difficulty, board.Difficulty)
		}
		if board.BoardSize != wantSize {
			t.Errorf("NewGame(%q).BoardSize = %d, want %d", difficulty, board.BoardSize, wantSize)
		}
		if len(board.Faces) != FaceCount {
			t.Fatalf("NewGame(%q) has %d faces, want %d", difficulty, len(board.Faces), FaceCount)
		}

		for f, face := range board.Faces {
			if len(face.Cells) != wantSize {
				t.Fatalf("face %d has %d rows, want %d", f, len(face.Cells), wantSize)
			}
			for row, cells := range face.Cells {
				if len(cells) != wantSize {
					t.Fatalf("face %d row %d has %d cols, want %d", f, row, len(cells), wantSize)
				}
				for col, cell := range cells {
					if cell.Face != f || cell.Row != row || cell.Col != col {
						t.Errorf("cell at (%d,%d,%d) carries coords (%d,%d,%d)",
							f, row, col, cell.Face, cell.Row, cell.Col)
					}
					if cell.State != CellStateClosed {
						t.Errorf("cell (%d,%d,%d).State = %q, want %q", f, row, col, cell.State, CellStateClosed)
					}
					if cell.AdjacentMines != 0 {
						t.Errorf("cell (%d,%d,%d).AdjacentMines = %d, want 0", f, row, col, cell.AdjacentMines)
					}
				}
			}
		}
	}
}

func TestNewGameInvalidDifficulty(t *testing.T) {
	board, err := NewGame(Difficulty("impossible"))
	if board != nil {
		t.Errorf("NewGame with invalid difficulty returned a board: %+v", board)
	}
	if !errors.Is(err, ErrInvalidDifficulty) {
		t.Errorf("NewGame with invalid difficulty returned %v, want %v", err, ErrInvalidDifficulty)
	}
}

func TestOpenCellReturnsNotImplemented(t *testing.T) {
	board, err := OpenCell(0, 0, 0)
	if board != nil {
		t.Errorf("OpenCell returned a board: %+v", board)
	}
	if !errors.Is(err, ErrNotImplemented) {
		t.Errorf("OpenCell returned %v, want %v", err, ErrNotImplemented)
	}
}

func TestFlagCellReturnsNotImplemented(t *testing.T) {
	board, err := FlagCell(0, 0, 0)
	if board != nil {
		t.Errorf("FlagCell returned a board: %+v", board)
	}
	if !errors.Is(err, ErrNotImplemented) {
		t.Errorf("FlagCell returned %v, want %v", err, ErrNotImplemented)
	}
}

func TestChordCellReturnsNotImplemented(t *testing.T) {
	board, err := ChordCell(0, 0, 0)
	if board != nil {
		t.Errorf("ChordCell returned a board: %+v", board)
	}
	if !errors.Is(err, ErrNotImplemented) {
		t.Errorf("ChordCell returned %v, want %v", err, ErrNotImplemented)
	}
}
