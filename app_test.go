package main

import (
	"errors"
	"reflect"
	"testing"

	"cubo-minado/internal/game"
)

func TestAppBindingsDelegateToGame(t *testing.T) {
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
