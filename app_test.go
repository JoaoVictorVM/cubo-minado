package main

import (
	"errors"
	"reflect"
	"testing"

	"cubo-minado/internal/game"
)

func useTempConfigDir(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("AppData", dir)
	t.Setenv("XDG_CONFIG_HOME", dir)
}

func TestAppBindingsDelegateToGame(t *testing.T) {
	useTempConfigDir(t)

	app := NewApp()

	for _, difficulty := range []string{"easy", "medium", "hard"} {
		got, err := app.NewGame(difficulty)
		if err != nil {
			t.Fatalf("App.NewGame(%q) returned error: %v", difficulty, err)
		}
		want, _ := game.NewGame(game.Difficulty(difficulty))
		if !reflect.DeepEqual(got, want) {
			t.Errorf("App.NewGame(%q) diverges from game.NewGame", difficulty)
		}
	}

	if _, err := app.NewGame("impossible"); !errors.Is(err, game.ErrInvalidDifficulty) {
		t.Errorf("App.NewGame with invalid difficulty returned %v, want %v", err, game.ErrInvalidDifficulty)
	}

	board, err := app.ChordCell(0, 0, 0)
	if board != nil {
		t.Errorf("App.ChordCell returned a board: %+v", board)
	}
	if !errors.Is(err, game.ErrNotImplemented) {
		t.Errorf("App.ChordCell returned %v, want %v", err, game.ErrNotImplemented)
	}

	gotTimes, err := app.GetBestTimes()
	if err != nil {
		t.Fatalf("App.GetBestTimes returned error: %v", err)
	}
	wantTimes, _ := game.GetBestTimes()
	if !reflect.DeepEqual(gotTimes, wantTimes) {
		t.Errorf("App.GetBestTimes = %+v, want %+v", gotTimes, wantTimes)
	}
}

func TestAppSubmitTimeDelegatesToGame(t *testing.T) {
	useTempConfigDir(t)

	app := NewApp()

	got, err := app.SubmitTime("medium", 87)
	if err != nil {
		t.Fatalf("App.SubmitTime returned error: %v", err)
	}
	want, _ := game.SubmitTime(game.DifficultyMedium, 87)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("App.SubmitTime = %+v, want %+v", got, want)
	}

	if _, err := app.SubmitTime("impossible", 87); !errors.Is(err, game.ErrInvalidDifficulty) {
		t.Errorf("App.SubmitTime with invalid difficulty returned %v, want %v", err, game.ErrInvalidDifficulty)
	}

	if _, err := app.SubmitTime("medium", -1); !errors.Is(err, game.ErrInvalidTime) {
		t.Errorf("App.SubmitTime with negative seconds returned %v, want %v", err, game.ErrInvalidTime)
	}
}

func TestAppGameplayBindingsDelegateToGame(t *testing.T) {
	useTempConfigDir(t)

	app := NewApp()

	if _, err := app.NewGame("easy"); err != nil {
		t.Fatalf("App.NewGame returned error: %v", err)
	}

	opened, err := app.OpenCell(2, 2, 2)
	if err != nil {
		t.Fatalf("App.OpenCell returned error: %v", err)
	}
	if got := opened.Faces[2].Cells[2][2].State; got != game.CellStateOpen {
		t.Errorf("App.OpenCell left the cell as %q, want %q", got, game.CellStateOpen)
	}

	flagged, err := app.FlagCell(0, 0, 0)
	if err != nil {
		t.Fatalf("App.FlagCell returned error: %v", err)
	}
	if got := flagged.Faces[0].Cells[0][0].State; got != game.CellStateFlagged {
		t.Errorf("App.FlagCell left the cell as %q, want %q", got, game.CellStateFlagged)
	}

	if _, err := app.OpenCell(9, 0, 0); !errors.Is(err, game.ErrInvalidCell) {
		t.Errorf("App.OpenCell with a bad face returned %v, want %v", err, game.ErrInvalidCell)
	}
	if _, err := app.FlagCell(0, 99, 0); !errors.Is(err, game.ErrInvalidCell) {
		t.Errorf("App.FlagCell with a bad row returned %v, want %v", err, game.ErrInvalidCell)
	}
}
