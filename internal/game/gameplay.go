package game

import (
	"math/rand"
	"time"
)

func newRandSource() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

func (m *matchState) cellAt(ref CellRef) *Cell {
	return &m.board.Faces[int(ref.Face)].Cells[ref.Row][ref.Col]
}

func (m *matchState) countAdjacentMines(ref CellRef) int {
	count := 0
	for _, neighbor := range m.graph[ref] {
		if m.mines[neighbor] {
			count++
		}
	}
	return count
}

func resolveAction(face, row, col int) (*matchState, CellRef, error) {
	if currentMatch == nil {
		return nil, CellRef{}, ErrNoActiveMatch
	}

	size := currentMatch.board.BoardSize
	if face < 0 || face >= FaceCount || row < 0 || row >= size || col < 0 || col >= size {
		return nil, CellRef{}, ErrInvalidCell
	}

	return currentMatch, CellRef{FaceID(face), row, col}, nil
}

func OpenCell(face, row, col int) (*BoardState, error) {
	matchMu.Lock()
	defer matchMu.Unlock()

	match, ref, err := resolveAction(face, row, col)
	if err != nil {
		return nil, err
	}

	cell := match.cellAt(ref)
	if cell.State != CellStateClosed {
		return match.board, nil
	}

	if !match.minesPlaced {
		match.mines = placeMines(match.graph, match.board.BoardSize, match.totalMines, ref, newRandSource())
		match.minesPlaced = true
	}

	cell.State = CellStateOpen
	if match.mines[ref] {
		cell.IsMine = true
		cell.AdjacentMines = 0
	} else {
		cell.IsMine = false
		cell.AdjacentMines = match.countAdjacentMines(ref)
	}

	return match.board, nil
}

func FlagCell(face, row, col int) (*BoardState, error) {
	matchMu.Lock()
	defer matchMu.Unlock()

	match, ref, err := resolveAction(face, row, col)
	if err != nil {
		return nil, err
	}

	cell := match.cellAt(ref)

	switch cell.State {
	case CellStateFlagged:
		cell.State = CellStateClosed
		match.flagsPlaced--
	case CellStateClosed:
		if match.flagsPlaced >= match.totalMines {
			return match.board, nil
		}
		cell.State = CellStateFlagged
		match.flagsPlaced++
	}

	return match.board, nil
}
