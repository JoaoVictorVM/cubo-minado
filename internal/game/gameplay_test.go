package game

import (
	"errors"
	"testing"
)

func forceMines(t *testing.T, mines ...CellRef) {
	t.Helper()
	if currentMatch == nil {
		t.Fatal("no active match to force a mine layout onto")
	}
	layout := make(map[CellRef]bool, len(mines))
	for _, cell := range mines {
		layout[cell] = true
	}
	currentMatch.mines = layout
	currentMatch.minesPlaced = true
}

func cellState(t *testing.T, board *BoardState, ref CellRef) Cell {
	t.Helper()
	return board.Faces[int(ref.Face)].Cells[ref.Row][ref.Col]
}

func countByState(board *BoardState, state CellState) int {
	total := 0
	for _, face := range board.Faces {
		for _, row := range face.Cells {
			for _, cell := range row {
				if cell.State == state {
					total++
				}
			}
		}
	}
	return total
}

func TestOpenCellFirstClickPlacesMinesSafely(t *testing.T) {
	for _, difficulty := range []Difficulty{DifficultyEasy, DifficultyMedium, DifficultyHard} {
		for _, first := range []CellRef{
			{FaceFront, 2, 2},
			{FaceTop, 0, 0},
			{FaceBack, 1, 3},
		} {
			startMatch(t, difficulty)

			if currentMatch.minesPlaced {
				t.Fatal("mines were placed before the first open")
			}

			board, err := OpenCell(int(first.Face), first.Row, first.Col)
			if err != nil {
				t.Fatalf("OpenCell returned error: %v", err)
			}

			if !currentMatch.minesPlaced {
				t.Fatal("first open did not place mines")
			}
			if len(currentMatch.mines) != currentMatch.totalMines {
				t.Errorf("placed %d mines, want %d", len(currentMatch.mines), currentMatch.totalMines)
			}
			if currentMatch.mines[first] {
				t.Errorf("%s: first opened cell %+v was a mine", difficulty, first)
			}
			for _, neighbor := range currentMatch.graph[first] {
				if currentMatch.mines[neighbor] {
					t.Errorf("%s: neighbor %+v of the first opened cell was a mine", difficulty, neighbor)
				}
			}

			opened := cellState(t, board, first)
			if opened.State != CellStateOpen {
				t.Errorf("first opened cell state = %q, want %q", opened.State, CellStateOpen)
			}
			if opened.IsMine {
				t.Error("first opened cell reported IsMine")
			}

			resetMatch(t)
		}
	}
}

func TestOpenCellReportsNeighborMineCount(t *testing.T) {
	startMatch(t, DifficultyEasy)

	target := CellRef{FaceFront, 2, 2}
	neighbors := currentMatch.graph[target]
	if len(neighbors) != 8 {
		t.Fatalf("expected an interior cell with 8 neighbors, got %d", len(neighbors))
	}

	forceMines(t, neighbors[0], neighbors[3], neighbors[7])

	board, err := OpenCell(int(target.Face), target.Row, target.Col)
	if err != nil {
		t.Fatalf("OpenCell returned error: %v", err)
	}

	opened := cellState(t, board, target)
	if opened.State != CellStateOpen {
		t.Errorf("state = %q, want %q", opened.State, CellStateOpen)
	}
	if opened.IsMine {
		t.Error("a non-mine cell reported IsMine")
	}
	if opened.AdjacentMines != 3 {
		t.Errorf("AdjacentMines = %d, want 3", opened.AdjacentMines)
	}
	if got := countByState(board, CellStateOpen); got != 1 {
		t.Errorf("%d cells are open, want exactly 1 (no cascade in Core Scope)", got)
	}
}

func TestOpenCellCountsMinesAcrossFaceBoundaries(t *testing.T) {
	startMatch(t, DifficultyEasy)

	target := CellRef{FaceFront, 0, 2}
	var crossFace []CellRef
	for _, neighbor := range currentMatch.graph[target] {
		if neighbor.Face != target.Face {
			crossFace = append(crossFace, neighbor)
		}
	}
	if len(crossFace) != 3 {
		t.Fatalf("expected 3 cross-face neighbors on an edge cell, got %d", len(crossFace))
	}

	forceMines(t, crossFace...)

	board, err := OpenCell(int(target.Face), target.Row, target.Col)
	if err != nil {
		t.Fatalf("OpenCell returned error: %v", err)
	}

	if got := cellState(t, board, target).AdjacentMines; got != 3 {
		t.Errorf("AdjacentMines = %d, want 3 from the adjacent face", got)
	}
}

func TestOpenCellOnMineRevealsMineNoCount(t *testing.T) {
	startMatch(t, DifficultyEasy)

	target := CellRef{FaceFront, 2, 2}
	neighbors := currentMatch.graph[target]
	forceMines(t, target, neighbors[0], neighbors[1])

	board, err := OpenCell(int(target.Face), target.Row, target.Col)
	if err != nil {
		t.Fatalf("OpenCell returned error: %v", err)
	}

	opened := cellState(t, board, target)
	if opened.State != CellStateOpen {
		t.Errorf("state = %q, want %q", opened.State, CellStateOpen)
	}
	if !opened.IsMine {
		t.Error("opening a mine did not report IsMine")
	}
	if opened.AdjacentMines != 0 {
		t.Errorf("AdjacentMines = %d on a mine cell, want 0", opened.AdjacentMines)
	}
}

func TestOpenCellDoesNotLeakUnopenedMines(t *testing.T) {
	startMatch(t, DifficultyEasy)

	board, err := OpenCell(2, 2, 2)
	if err != nil {
		t.Fatalf("OpenCell returned error: %v", err)
	}

	for _, face := range board.Faces {
		for _, row := range face.Cells {
			for _, cell := range row {
				if cell.State != CellStateOpen && cell.IsMine {
					t.Errorf("unopened cell (%d,%d,%d) reports IsMine", cell.Face, cell.Row, cell.Col)
				}
			}
		}
	}
}

func TestOpenCellAlreadyOpenIsNoOp(t *testing.T) {
	startMatch(t, DifficultyEasy)

	target := CellRef{FaceFront, 2, 2}
	forceMines(t, currentMatch.graph[target][0])

	first, err := OpenCell(int(target.Face), target.Row, target.Col)
	if err != nil {
		t.Fatalf("first OpenCell returned error: %v", err)
	}
	before := cellState(t, first, target)
	openedBefore := countByState(first, CellStateOpen)

	second, err := OpenCell(int(target.Face), target.Row, target.Col)
	if err != nil {
		t.Fatalf("second OpenCell returned error: %v", err)
	}
	if second == nil {
		t.Fatal("second OpenCell returned a nil board")
	}
	if cellState(t, second, target) != before {
		t.Error("reopening an open cell changed it")
	}
	if got := countByState(second, CellStateOpen); got != openedBefore {
		t.Errorf("%d cells open after reopening, want %d", got, openedBefore)
	}
}

func TestOpenCellOnFlaggedIsNoOp(t *testing.T) {
	startMatch(t, DifficultyEasy)

	target := CellRef{FaceFront, 2, 2}
	if _, err := FlagCell(int(target.Face), target.Row, target.Col); err != nil {
		t.Fatalf("FlagCell returned error: %v", err)
	}

	board, err := OpenCell(int(target.Face), target.Row, target.Col)
	if err != nil {
		t.Fatalf("OpenCell returned error: %v", err)
	}

	if got := cellState(t, board, target).State; got != CellStateFlagged {
		t.Errorf("state = %q, want %q", got, CellStateFlagged)
	}
	if got := countByState(board, CellStateOpen); got != 0 {
		t.Errorf("%d cells open, want 0", got)
	}
	if currentMatch.minesPlaced {
		t.Error("a rejected open triggered mine placement")
	}
}

func TestOpenCellInvalidCellRejected(t *testing.T) {
	startMatch(t, DifficultyEasy)

	cases := [][3]int{
		{-1, 0, 0},
		{6, 0, 0},
		{0, -1, 0},
		{0, 5, 0},
		{0, 0, -1},
		{0, 0, 5},
	}

	for _, tc := range cases {
		board, err := OpenCell(tc[0], tc[1], tc[2])
		if board != nil {
			t.Errorf("OpenCell%v returned a board", tc)
		}
		if !errors.Is(err, ErrInvalidCell) {
			t.Errorf("OpenCell%v returned %v, want %v", tc, err, ErrInvalidCell)
		}
	}
}

func TestOpenCellNoActiveMatchRejected(t *testing.T) {
	resetMatch(t)

	board, err := OpenCell(0, 0, 0)
	if board != nil {
		t.Errorf("OpenCell returned a board with no active match: %+v", board)
	}
	if !errors.Is(err, ErrNoActiveMatch) {
		t.Errorf("OpenCell returned %v, want %v", err, ErrNoActiveMatch)
	}
}

func TestFlagCellTogglesClosedCell(t *testing.T) {
	startMatch(t, DifficultyEasy)

	target := CellRef{FaceFront, 1, 1}

	board, err := FlagCell(int(target.Face), target.Row, target.Col)
	if err != nil {
		t.Fatalf("FlagCell returned error: %v", err)
	}
	if got := cellState(t, board, target).State; got != CellStateFlagged {
		t.Errorf("state = %q, want %q", got, CellStateFlagged)
	}
	if currentMatch.flagsPlaced != 1 {
		t.Errorf("flagsPlaced = %d, want 1", currentMatch.flagsPlaced)
	}

	board, err = FlagCell(int(target.Face), target.Row, target.Col)
	if err != nil {
		t.Fatalf("second FlagCell returned error: %v", err)
	}
	if got := cellState(t, board, target).State; got != CellStateClosed {
		t.Errorf("state = %q, want %q", got, CellStateClosed)
	}
	if currentMatch.flagsPlaced != 0 {
		t.Errorf("flagsPlaced = %d, want 0", currentMatch.flagsPlaced)
	}
}

func TestFlagCellRejectedOnOpenCell(t *testing.T) {
	startMatch(t, DifficultyEasy)

	target := CellRef{FaceFront, 2, 2}
	if _, err := OpenCell(int(target.Face), target.Row, target.Col); err != nil {
		t.Fatalf("OpenCell returned error: %v", err)
	}

	board, err := FlagCell(int(target.Face), target.Row, target.Col)
	if err != nil {
		t.Fatalf("FlagCell returned error: %v", err)
	}
	if got := cellState(t, board, target).State; got != CellStateOpen {
		t.Errorf("state = %q, want %q", got, CellStateOpen)
	}
	if currentMatch.flagsPlaced != 0 {
		t.Errorf("flagsPlaced = %d, want 0", currentMatch.flagsPlaced)
	}
}

func TestFlagCellCappedAtTotalMines(t *testing.T) {
	startMatch(t, DifficultyEasy)

	total := currentMatch.totalMines
	size := currentMatch.board.BoardSize

	placed := 0
	var extra CellRef
	for face := 0; face < FaceCount && placed <= total; face++ {
		for row := 0; row < size && placed <= total; row++ {
			for col := 0; col < size && placed <= total; col++ {
				if placed == total {
					extra = CellRef{FaceID(face), row, col}
					placed++
					continue
				}
				if _, err := FlagCell(face, row, col); err != nil {
					t.Fatalf("FlagCell returned error: %v", err)
				}
				placed++
			}
		}
	}

	if currentMatch.flagsPlaced != total {
		t.Fatalf("flagsPlaced = %d after filling the cap, want %d", currentMatch.flagsPlaced, total)
	}

	board, err := FlagCell(int(extra.Face), extra.Row, extra.Col)
	if err != nil {
		t.Fatalf("over-cap FlagCell returned error: %v", err)
	}
	if got := cellState(t, board, extra).State; got != CellStateClosed {
		t.Errorf("over-cap flag set state to %q, want %q", got, CellStateClosed)
	}
	if currentMatch.flagsPlaced != total {
		t.Errorf("flagsPlaced = %d after an over-cap attempt, want %d", currentMatch.flagsPlaced, total)
	}

	if _, err := FlagCell(0, 0, 0); err != nil {
		t.Fatalf("unflag returned error: %v", err)
	}
	if currentMatch.flagsPlaced != total-1 {
		t.Fatalf("flagsPlaced = %d after unflagging, want %d", currentMatch.flagsPlaced, total-1)
	}

	board, err = FlagCell(int(extra.Face), extra.Row, extra.Col)
	if err != nil {
		t.Fatalf("FlagCell after freeing a slot returned error: %v", err)
	}
	if got := cellState(t, board, extra).State; got != CellStateFlagged {
		t.Errorf("state = %q after freeing a flag slot, want %q", got, CellStateFlagged)
	}
}

func TestFlagCellInvalidCellRejected(t *testing.T) {
	startMatch(t, DifficultyEasy)

	for _, tc := range [][3]int{{6, 0, 0}, {0, 5, 0}, {0, 0, -1}} {
		board, err := FlagCell(tc[0], tc[1], tc[2])
		if board != nil {
			t.Errorf("FlagCell%v returned a board", tc)
		}
		if !errors.Is(err, ErrInvalidCell) {
			t.Errorf("FlagCell%v returned %v, want %v", tc, err, ErrInvalidCell)
		}
	}
}

func TestFlagCellNoActiveMatchRejected(t *testing.T) {
	resetMatch(t)

	board, err := FlagCell(0, 0, 0)
	if board != nil {
		t.Errorf("FlagCell returned a board with no active match: %+v", board)
	}
	if !errors.Is(err, ErrNoActiveMatch) {
		t.Errorf("FlagCell returned %v, want %v", err, ErrNoActiveMatch)
	}
}

func TestFlagCellDoesNotPlaceMines(t *testing.T) {
	startMatch(t, DifficultyEasy)

	if _, err := FlagCell(0, 0, 0); err != nil {
		t.Fatalf("FlagCell returned error: %v", err)
	}
	if currentMatch.minesPlaced {
		t.Error("flagging placed mines; only the first open should")
	}
}
