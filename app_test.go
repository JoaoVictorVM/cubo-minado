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

	gameplayBindings := map[string]func(int, int, int) (*game.BoardState, error){
		"OpenCell":  app.OpenCell,
		"FlagCell":  app.FlagCell,
		"ChordCell": app.ChordCell,
	}
	for name, binding := range gameplayBindings {
		board, err := binding(0, 0, 0)
		if board != nil {
			t.Errorf("App.%s returned a board: %+v", name, board)
		}
		if !errors.Is(err, game.ErrNotImplemented) {
			t.Errorf("App.%s returned %v, want %v", name, err, game.ErrNotImplemented)
		}
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
