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

type GameResult string

const (
	GameResultNone    GameResult = ""
	GameResultVictory GameResult = "victory"
	GameResultDefeat  GameResult = "defeat"
)

const FaceCount = 6

type Cell struct {
	Face          int       `json:"face"`
	Row           int       `json:"row"`
	Col           int       `json:"col"`
	State         CellState `json:"state"`
	AdjacentMines int       `json:"adjacentMines"`
	IsMine        bool      `json:"isMine"`
}

type Face struct {
	Cells [][]Cell `json:"cells"`
}

type BoardState struct {
	Difficulty Difficulty `json:"difficulty"`
	BoardSize  int        `json:"boardSize"`
	Faces      []Face     `json:"faces"`
	Result     GameResult `json:"result"`
	StartedAt  *int64     `json:"startedAt,omitempty"`
	EndedAt    *int64     `json:"endedAt,omitempty"`
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

func MinesForDifficulty(difficulty Difficulty) (int, bool) {
	switch difficulty {
	case DifficultyEasy:
		return 15, true
	case DifficultyMedium:
		return 40, true
	case DifficultyHard:
		return 80, true
	default:
		return 0, false
	}
}
