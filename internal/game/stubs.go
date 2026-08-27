package game

func NewGame(difficulty Difficulty) (*BoardState, error) {
	size, ok := BoardSizeForDifficulty(difficulty)
	if !ok {
		return nil, ErrInvalidDifficulty
	}

	faces := make([]Face, FaceCount)
	for f := range faces {
		cells := make([][]Cell, size)
		for row := range cells {
			cells[row] = make([]Cell, size)
			for col := range cells[row] {
				cells[row][col] = Cell{
					Face:          f,
					Row:           row,
					Col:           col,
					State:         CellStateClosed,
					AdjacentMines: 0,
				}
			}
		}
		faces[f] = Face{Cells: cells}
	}

	return &BoardState{
		Difficulty: difficulty,
		BoardSize:  size,
		Faces:      faces,
	}, nil
}

func OpenCell(face, row, col int) (*BoardState, error) {
	return nil, ErrNotImplemented
}

func FlagCell(face, row, col int) (*BoardState, error) {
	return nil, ErrNotImplemented
}

func ChordCell(face, row, col int) (*BoardState, error) {
	return nil, ErrNotImplemented
}
