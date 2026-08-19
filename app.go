package main

import (
	"context"

	"cubo-minado/internal/game"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown(ctx context.Context) {
	a.ctx = nil
}

func (a *App) NewGame(difficulty string) (*game.BoardState, error) {
	return game.NewGame(game.Difficulty(difficulty))
}

func (a *App) OpenCell(face, row, col int) (*game.BoardState, error) {
	return game.OpenCell(face, row, col)
}

func (a *App) FlagCell(face, row, col int) (*game.BoardState, error) {
	return game.FlagCell(face, row, col)
}

func (a *App) ChordCell(face, row, col int) (*game.BoardState, error) {
	return game.ChordCell(face, row, col)
}

func (a *App) GetBestTimes() (*game.BestTimes, error) {
	return game.GetBestTimes()
}
