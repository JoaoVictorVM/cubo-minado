package game

import (
	"errors"
	"testing"
)

func resetMatch(t *testing.T) {
	t.Helper()
	matchMu.Lock()
	currentMatch = nil
	matchMu.Unlock()
}

func startMatch(t *testing.T, difficulty Difficulty) *BoardState {
	t.Helper()
	board, err := NewGame(difficulty)
	if err != nil {
		t.Fatalf("NewGame(%q) returned error: %v", difficulty, err)
	}
	t.Cleanup(func() { resetMatch(t) })
	return board
}

func TestNewGameValidDifficulties(t *testing.T) {
	cases := map[Difficulty]int{
		DifficultyEasy:   5,
		DifficultyMedium: 7,
		DifficultyHard:   9,
	}

	for difficulty, wantSize := range cases {
		board := startMatch(t, difficulty)

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
					if cell.IsMine {
						t.Errorf("cell (%d,%d,%d).IsMine is true on a fresh board", f, row, col)
					}
				}
			}
		}
	}
}

func TestNewGameInvalidDifficulty(t *testing.T) {
	resetMatch(t)

	board, err := NewGame(Difficulty("impossible"))
	if board != nil {
		t.Errorf("NewGame with invalid difficulty returned a board: %+v", board)
	}
	if !errors.Is(err, ErrInvalidDifficulty) {
		t.Errorf("NewGame with invalid difficulty returned %v, want %v", err, ErrInvalidDifficulty)
	}
	if currentMatch != nil {
		t.Error("NewGame with invalid difficulty installed a match")
	}
}

func TestNewGameInstallsMatchState(t *testing.T) {
	cases := map[Difficulty]struct {
		size  int
		mines int
	}{
		DifficultyEasy:   {5, 15},
		DifficultyMedium: {7, 40},
		DifficultyHard:   {9, 80},
	}

	for difficulty, want := range cases {
		board := startMatch(t, difficulty)

		if currentMatch == nil {
			t.Fatalf("NewGame(%q) did not install a match", difficulty)
		}
		if currentMatch.board != board {
			t.Errorf("NewGame(%q) installed a different board than it returned", difficulty)
		}
		if got := len(currentMatch.graph); got != FaceCount*want.size*want.size {
			t.Errorf("NewGame(%q) graph has %d cells, want %d", difficulty, got, FaceCount*want.size*want.size)
		}
		if currentMatch.totalMines != want.mines {
			t.Errorf("NewGame(%q).totalMines = %d, want %d", difficulty, currentMatch.totalMines, want.mines)
		}
		if currentMatch.minesPlaced {
			t.Errorf("NewGame(%q) placed mines before the first open", difficulty)
		}
		if len(currentMatch.mines) != 0 {
			t.Errorf("NewGame(%q) installed %d mines before the first open", difficulty, len(currentMatch.mines))
		}
		if currentMatch.flagsPlaced != 0 {
			t.Errorf("NewGame(%q).flagsPlaced = %d, want 0", difficulty, currentMatch.flagsPlaced)
		}
	}
}

func TestNewGameReplacesPriorMatch(t *testing.T) {
	first := startMatch(t, DifficultyEasy)

	if _, err := OpenCell(0, 2, 2); err != nil {
		t.Fatalf("OpenCell returned error: %v", err)
	}
	if _, err := FlagCell(1, 0, 0); err != nil {
		t.Fatalf("FlagCell returned error: %v", err)
	}
	if !currentMatch.minesPlaced || currentMatch.flagsPlaced != 1 {
		t.Fatal("first match did not reach the expected in-progress state")
	}

	second := startMatch(t, DifficultyHard)

	if second == first {
		t.Error("NewGame returned the previous match's board")
	}
	if currentMatch.board != second {
		t.Error("currentMatch still points at the previous board")
	}
	if currentMatch.board.BoardSize != 9 {
		t.Errorf("replacement match board size = %d, want 9", currentMatch.board.BoardSize)
	}
	if currentMatch.minesPlaced || len(currentMatch.mines) != 0 {
		t.Error("replacement match kept the previous mine layout")
	}
	if currentMatch.flagsPlaced != 0 {
		t.Errorf("replacement match kept flagsPlaced = %d, want 0", currentMatch.flagsPlaced)
	}
	if currentMatch.totalMines != 80 {
		t.Errorf("replacement match totalMines = %d, want 80", currentMatch.totalMines)
	}
}
