package game

import "testing"

func TestBoardSizeForDifficulty(t *testing.T) {
	cases := []struct {
		difficulty Difficulty
		want       int
		wantOK     bool
	}{
		{DifficultyEasy, 5, true},
		{DifficultyMedium, 7, true},
		{DifficultyHard, 9, true},
		{Difficulty("impossible"), 0, false},
	}

	for _, tc := range cases {
		got, ok := BoardSizeForDifficulty(tc.difficulty)
		if got != tc.want || ok != tc.wantOK {
			t.Errorf("BoardSizeForDifficulty(%q) = (%d, %t), want (%d, %t)",
				tc.difficulty, got, ok, tc.want, tc.wantOK)
		}
	}
}
