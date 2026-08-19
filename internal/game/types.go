package game

type Difficulty string

const (
	DifficultyEasy   Difficulty = "easy"
	DifficultyMedium Difficulty = "medium"
	DifficultyHard   Difficulty = "hard"
)

type CellState string

const (
	CellStateClosed  CellState = "closed"
	CellStateOpen    CellState = "open"
	CellStateFlagged CellState = "flagged"
)

const FaceCount = 6

type Cell struct {
	Face          int       `json:"face"`
	Row           int       `json:"row"`
	Col           int       `json:"col"`
	State         CellState `json:"state"`
	AdjacentMines int       `json:"adjacentMines"`
}

type Face struct {
	Cells [][]Cell `json:"cells"`
}

type BoardState struct {
	Difficulty Difficulty `json:"difficulty"`
	BoardSize  int        `json:"boardSize"`
	Faces      []Face     `json:"faces"`
}

type BestTimes struct {
	Easy   *int `json:"easy"`
	Medium *int `json:"medium"`
	Hard   *int `json:"hard"`
}

func BoardSizeForDifficulty(difficulty Difficulty) (int, bool) {
	switch difficulty {
	case DifficultyEasy:
		return 5, true
	case DifficultyMedium:
		return 7, true
	case DifficultyHard:
		return 9, true
	default:
		return 0, false
	}
}
