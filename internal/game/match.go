package game

import "sync"

type matchState struct {
	board       *BoardState
	graph       Graph
	mines       map[CellRef]bool
	minesPlaced bool
	totalMines  int
	flagsPlaced int
}

func (m *matchState) hasEnded() bool {
	return m.board.Result != GameResultNone
}

var (
	matchMu      sync.Mutex
	currentMatch *matchState
)

func NewGame(difficulty Difficulty) (*BoardState, error) {
	size, ok := BoardSizeForDifficulty(difficulty)
	if !ok {
		return nil, ErrInvalidDifficulty
	}

	totalMines, ok := MinesForDifficulty(difficulty)
	if !ok {
		return nil, ErrInvalidDifficulty
	}

	graph, err := BuildAdjacencyGraph(size)
	if err != nil {
		return nil, err
	}

	board := emptyBoard(difficulty, size)

	matchMu.Lock()
	defer matchMu.Unlock()

	currentMatch = &matchState{
		board:      board,
		graph:      graph,
		mines:      make(map[CellRef]bool, totalMines),
		totalMines: totalMines,
	}

	return board, nil
}

func emptyBoard(difficulty Difficulty, size int) *BoardState {
	faces := make([]Face, FaceCount)
	for f := range faces {
		cells := make([][]Cell, size)
		for row := range cells {
			cells[row] = make([]Cell, size)
			for col := range cells[row] {
				cells[row][col] = Cell{
					Face:  f,
					Row:   row,
					Col:   col,
					State: CellStateClosed,
				}
			}
		}
		faces[f] = Face{Cells: cells}
	}

	return &BoardState{
		Difficulty: difficulty,
		BoardSize:  size,
		Faces:      faces,
	}
}
