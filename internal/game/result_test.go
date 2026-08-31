package game

import (
	"errors"
	"testing"
)

func openEverySafeCellExcept(t *testing.T, skip CellRef) {
	t.Helper()
	size := currentMatch.board.BoardSize

	for face := 0; face < FaceCount; face++ {
		for row := 0; row < size; row++ {
			for col := 0; col < size; col++ {
				ref := CellRef{FaceID(face), row, col}
				if ref == skip || currentMatch.mines[ref] {
					continue
				}
				board, err := OpenCell(face, row, col)
				if err != nil {
					t.Fatalf("OpenCell%+v returned error: %v", ref, err)
				}
				if board.Result != GameResultNone {
					t.Fatalf("match ended early at %+v with result %q", ref, board.Result)
				}
			}
		}
	}
}

func playToVictory(t *testing.T) *BoardState {
	t.Helper()

	first := CellRef{FaceFront, 2, 2}
	if _, err := OpenCell(int(first.Face), first.Row, first.Col); err != nil {
		t.Fatalf("first OpenCell returned error: %v", err)
	}

	size := currentMatch.board.BoardSize
	for face := 0; face < FaceCount; face++ {
		for row := 0; row < size; row++ {
			for col := 0; col < size; col++ {
				if currentMatch.mines[CellRef{FaceID(face), row, col}] {
					continue
				}
				if _, err := OpenCell(face, row, col); err != nil {
					t.Fatalf("OpenCell(%d,%d,%d) returned error: %v", face, row, col, err)
				}
			}
		}
	}

	board := currentMatch.board
	if board.Result != GameResultVictory {
		t.Fatalf("opening every safe cell gave result %q, want %q", board.Result, GameResultVictory)
	}
	return board
}

func snapshotStates(board *BoardState) map[CellRef]CellState {
	states := make(map[CellRef]CellState)
	for f, face := range board.Faces {
		for r, row := range face.Cells {
			for c, cell := range row {
				states[CellRef{FaceID(f), r, c}] = cell.State
			}
		}
	}
	return states
}

func TestOpenCellOngoingMatchLeavesResultEmpty(t *testing.T) {
	startMatch(t, DifficultyEasy)

	target := CellRef{FaceFront, 2, 2}
	forceMines(t, CellRef{FaceBack, 0, 0}, CellRef{FaceBack, 0, 1})

	board, err := OpenCell(int(target.Face), target.Row, target.Col)
	if err != nil {
		t.Fatalf("OpenCell returned error: %v", err)
	}

	if board.Result != GameResultNone {
		t.Errorf("Result = %q on an ongoing match, want %q", board.Result, GameResultNone)
	}
	if board.EndedAt != nil {
		t.Errorf("EndedAt = %v on an ongoing match, want nil", *board.EndedAt)
	}
}

func TestOpenCellOnMineTriggersDefeat(t *testing.T) {
	startMatch(t, DifficultyEasy)

	trigger := CellRef{FaceFront, 2, 2}
	forceMines(t, trigger, CellRef{FaceBack, 0, 0})

	board, err := OpenCell(int(trigger.Face), trigger.Row, trigger.Col)
	if err != nil {
		t.Fatalf("OpenCell returned error: %v", err)
	}

	if board.Result != GameResultDefeat {
		t.Errorf("Result = %q, want %q", board.Result, GameResultDefeat)
	}
	if board.EndedAt == nil {
		t.Error("EndedAt is nil after a defeat")
	}
}

func TestOpenCellDefeatRevealsAllMines(t *testing.T) {
	startMatch(t, DifficultyEasy)

	trigger := CellRef{FaceFront, 2, 2}
	closedMine := CellRef{FaceBack, 1, 1}
	flaggedMine := CellRef{FaceTop, 3, 3}
	forceMines(t, trigger, closedMine, flaggedMine)

	if _, err := FlagCell(int(flaggedMine.Face), flaggedMine.Row, flaggedMine.Col); err != nil {
		t.Fatalf("FlagCell returned error: %v", err)
	}

	board, err := OpenCell(int(trigger.Face), trigger.Row, trigger.Col)
	if err != nil {
		t.Fatalf("OpenCell returned error: %v", err)
	}

	for _, mine := range []CellRef{trigger, closedMine, flaggedMine} {
		if !cellState(t, board, mine).IsMine {
			t.Errorf("mine %+v was not revealed", mine)
		}
	}

	for f, face := range board.Faces {
		for r, row := range face.Cells {
			for c, cell := range row {
				ref := CellRef{FaceID(f), r, c}
				if !currentMatch.mines[ref] && cell.IsMine {
					t.Errorf("non-mine cell %+v reports IsMine", ref)
				}
			}
		}
	}
}

func TestOpenCellDefeatKeepsTriggeringCellDistinguishable(t *testing.T) {
	startMatch(t, DifficultyEasy)

	trigger := CellRef{FaceFront, 2, 2}
	closedMine := CellRef{FaceBack, 1, 1}
	flaggedMine := CellRef{FaceTop, 3, 3}
	forceMines(t, trigger, closedMine, flaggedMine)

	if _, err := FlagCell(int(flaggedMine.Face), flaggedMine.Row, flaggedMine.Col); err != nil {
		t.Fatalf("FlagCell returned error: %v", err)
	}

	board, err := OpenCell(int(trigger.Face), trigger.Row, trigger.Col)
	if err != nil {
		t.Fatalf("OpenCell returned error: %v", err)
	}

	openMines := 0
	for mine := range currentMatch.mines {
		if cellState(t, board, mine).State == CellStateOpen {
			openMines++
		}
	}
	if openMines != 1 {
		t.Errorf("%d mine cells are open, want exactly 1 (the triggering one)", openMines)
	}

	if got := cellState(t, board, trigger).State; got != CellStateOpen {
		t.Errorf("triggering mine state = %q, want %q", got, CellStateOpen)
	}
	if got := cellState(t, board, closedMine).State; got != CellStateClosed {
		t.Errorf("untouched mine state = %q, want %q", got, CellStateClosed)
	}
	if got := cellState(t, board, flaggedMine).State; got != CellStateFlagged {
		t.Errorf("flagged mine state = %q, want %q", got, CellStateFlagged)
	}
}

func TestOpenCellLastSafeCellTriggersVictory(t *testing.T) {
	startMatch(t, DifficultyEasy)

	first := CellRef{FaceFront, 2, 2}
	if _, err := OpenCell(int(first.Face), first.Row, first.Col); err != nil {
		t.Fatalf("first OpenCell returned error: %v", err)
	}

	last := CellRef{}
	found := false
	size := currentMatch.board.BoardSize
	for face := 0; face < FaceCount && !found; face++ {
		for row := 0; row < size && !found; row++ {
			for col := 0; col < size && !found; col++ {
				ref := CellRef{FaceID(face), row, col}
				if ref != first && !currentMatch.mines[ref] {
					last = ref
					found = true
				}
			}
		}
	}
	if !found {
		t.Fatal("could not pick a final safe cell")
	}

	openEverySafeCellExcept(t, last)

	board, err := OpenCell(int(last.Face), last.Row, last.Col)
	if err != nil {
		t.Fatalf("final OpenCell returned error: %v", err)
	}

	if board.Result != GameResultVictory {
		t.Errorf("Result = %q, want %q", board.Result, GameResultVictory)
	}
	if board.EndedAt == nil {
		t.Error("EndedAt is nil after a victory")
	}
}

func TestVictoryLeavesMinesUnopened(t *testing.T) {
	startMatch(t, DifficultyEasy)

	board := playToVictory(t)

	for mine := range currentMatch.mines {
		if got := cellState(t, board, mine).State; got == CellStateOpen {
			t.Errorf("mine %+v is open after a victory", mine)
		}
	}
}

func TestOpenCellRejectedAfterDefeat(t *testing.T) {
	startMatch(t, DifficultyEasy)

	trigger := CellRef{FaceFront, 2, 2}
	other := CellRef{FaceBack, 4, 4}
	forceMines(t, trigger, CellRef{FaceBack, 0, 0})

	if _, err := OpenCell(int(trigger.Face), trigger.Row, trigger.Col); err != nil {
		t.Fatalf("OpenCell returned error: %v", err)
	}
	before := snapshotStates(currentMatch.board)

	board, err := OpenCell(int(other.Face), other.Row, other.Col)
	if err != nil {
		t.Fatalf("post-defeat OpenCell returned error: %v", err)
	}
	if board == nil {
		t.Fatal("post-defeat OpenCell returned a nil board")
	}
	if board.Result != GameResultDefeat {
		t.Errorf("Result = %q after a rejected open, want %q", board.Result, GameResultDefeat)
	}
	if got := cellState(t, board, other).State; got != CellStateClosed {
		t.Errorf("rejected target state = %q, want %q", got, CellStateClosed)
	}

	for ref, state := range snapshotStates(board) {
		if before[ref] != state {
			t.Errorf("cell %+v changed after a rejected open: %q to %q", ref, before[ref], state)
		}
	}
}

func TestOpenCellRejectedAfterVictory(t *testing.T) {
	startMatch(t, DifficultyEasy)

	before := snapshotStates(playToVictory(t))

	var mine CellRef
	for m := range currentMatch.mines {
		mine = m
		break
	}

	board, err := OpenCell(int(mine.Face), mine.Row, mine.Col)
	if err != nil {
		t.Fatalf("post-victory OpenCell returned error: %v", err)
	}
	if board.Result != GameResultVictory {
		t.Errorf("Result = %q after a rejected open, want %q", board.Result, GameResultVictory)
	}
	for ref, state := range snapshotStates(board) {
		if before[ref] != state {
			t.Errorf("cell %+v changed after a rejected open: %q to %q", ref, before[ref], state)
		}
	}
}

func TestFlagCellRejectedAfterDefeat(t *testing.T) {
	startMatch(t, DifficultyEasy)

	trigger := CellRef{FaceFront, 2, 2}
	target := CellRef{FaceBack, 4, 4}
	forceMines(t, trigger, CellRef{FaceBack, 0, 0})

	if _, err := OpenCell(int(trigger.Face), trigger.Row, trigger.Col); err != nil {
		t.Fatalf("OpenCell returned error: %v", err)
	}
	flagsBefore := currentMatch.flagsPlaced

	board, err := FlagCell(int(target.Face), target.Row, target.Col)
	if err != nil {
		t.Fatalf("post-defeat FlagCell returned error: %v", err)
	}
	if got := cellState(t, board, target).State; got != CellStateClosed {
		t.Errorf("state = %q after a rejected flag, want %q", got, CellStateClosed)
	}
	if currentMatch.flagsPlaced != flagsBefore {
		t.Errorf("flagsPlaced = %d after a rejected flag, want %d", currentMatch.flagsPlaced, flagsBefore)
	}
}

func TestFlagCellRejectedAfterVictory(t *testing.T) {
	startMatch(t, DifficultyEasy)

	playToVictory(t)

	var mine CellRef
	for m := range currentMatch.mines {
		mine = m
		break
	}
	flagsBefore := currentMatch.flagsPlaced

	board, err := FlagCell(int(mine.Face), mine.Row, mine.Col)
	if err != nil {
		t.Fatalf("post-victory FlagCell returned error: %v", err)
	}
	if got := cellState(t, board, mine).State; got != CellStateClosed {
		t.Errorf("state = %q after a rejected flag, want %q", got, CellStateClosed)
	}
	if currentMatch.flagsPlaced != flagsBefore {
		t.Errorf("flagsPlaced = %d after a rejected flag, want %d", currentMatch.flagsPlaced, flagsBefore)
	}
}

func TestFlagCellStillWorksBeforeResult(t *testing.T) {
	startMatch(t, DifficultyEasy)

	target := CellRef{FaceFront, 1, 1}

	board, err := FlagCell(int(target.Face), target.Row, target.Col)
	if err != nil {
		t.Fatalf("FlagCell returned error: %v", err)
	}
	if got := cellState(t, board, target).State; got != CellStateFlagged {
		t.Errorf("state = %q, want %q", got, CellStateFlagged)
	}
	if board.Result != GameResultNone {
		t.Errorf("flagging set Result to %q", board.Result)
	}

	board, err = FlagCell(int(target.Face), target.Row, target.Col)
	if err != nil {
		t.Fatalf("second FlagCell returned error: %v", err)
	}
	if got := cellState(t, board, target).State; got != CellStateClosed {
		t.Errorf("state = %q, want %q", got, CellStateClosed)
	}
}

func TestNewGameClearsPriorResult(t *testing.T) {
	startMatch(t, DifficultyEasy)

	trigger := CellRef{FaceFront, 2, 2}
	forceMines(t, trigger)
	if _, err := OpenCell(int(trigger.Face), trigger.Row, trigger.Col); err != nil {
		t.Fatalf("OpenCell returned error: %v", err)
	}
	if currentMatch.board.Result != GameResultDefeat {
		t.Fatalf("expected a defeat, got %q", currentMatch.board.Result)
	}

	board := startMatch(t, DifficultyEasy)

	if board.Result != GameResultNone {
		t.Errorf("new match Result = %q, want %q", board.Result, GameResultNone)
	}
	if board.EndedAt != nil {
		t.Error("new match carries an EndedAt")
	}
	if currentMatch.hasEnded() {
		t.Error("new match reports itself as ended")
	}

	if _, err := OpenCell(0, 0, 0); err != nil {
		t.Fatalf("OpenCell on the new match returned error: %v", err)
	}
	if got := cellState(t, currentMatch.board, CellRef{FaceTop, 0, 0}).State; got != CellStateOpen {
		t.Errorf("new match does not accept actions: state = %q", got)
	}
}

func TestChordCellStillNotImplementedAfterF07(t *testing.T) {
	startMatch(t, DifficultyEasy)

	trigger := CellRef{FaceFront, 2, 2}
	forceMines(t, trigger)
	if _, err := OpenCell(int(trigger.Face), trigger.Row, trigger.Col); err != nil {
		t.Fatalf("OpenCell returned error: %v", err)
	}

	board, err := ChordCell(0, 0, 0)
	if board != nil {
		t.Errorf("ChordCell returned a board: %+v", board)
	}
	if !errors.Is(err, ErrNotImplemented) {
		t.Errorf("ChordCell returned %v, want %v", err, ErrNotImplemented)
	}
}
