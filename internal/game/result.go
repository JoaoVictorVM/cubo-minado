package game

import "time"

func evaluateResult(match *matchState, opened CellRef) {
	if match.mines[opened] {
		revealMines(match)
		endMatch(match, GameResultDefeat)
		return
	}

	if allSafeCellsOpened(match) {
		endMatch(match, GameResultVictory)
	}
}

func endMatch(match *matchState, result GameResult) {
	endedAt := time.Now().Unix()
	match.board.Result = result
	match.board.EndedAt = &endedAt
}

func revealMines(match *matchState) {
	for mine := range match.mines {
		match.cellAt(mine).IsMine = true
	}
}

func allSafeCellsOpened(match *matchState) bool {
	size := match.board.BoardSize
	safeCells := FaceCount*size*size - len(match.mines)

	opened := 0
	for _, face := range match.board.Faces {
		for _, row := range face.Cells {
			for _, cell := range row {
				if cell.State == CellStateOpen {
					opened++
				}
			}
		}
	}

	return opened >= safeCells
}
