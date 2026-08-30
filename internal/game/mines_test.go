package game

import (
	"math/rand"
	"testing"
)

func TestPlaceMinesExcludesSafeZone(t *testing.T) {
	size := 5
	graph := mustBuildGraph(t, size)
	safeCells := []CellRef{
		{FaceFront, 2, 2},
		{FaceTop, 0, 0},
		{FaceBack, 4, 4},
		{FaceLeft, 0, 2},
	}

	for _, safe := range safeCells {
		for seed := int64(1); seed <= 20; seed++ {
			mines := placeMines(graph, size, 15, safe, rand.New(rand.NewSource(seed)))

			if mines[safe] {
				t.Errorf("seed %d: opened cell %+v was mined", seed, safe)
			}
			for _, neighbor := range graph[safe] {
				if mines[neighbor] {
					t.Errorf("seed %d: neighbor %+v of opened cell %+v was mined", seed, neighbor, safe)
				}
			}
		}
	}
}

func TestPlaceMinesCount(t *testing.T) {
	cases := []struct {
		size  int
		mines int
	}{
		{5, 15},
		{7, 40},
		{9, 80},
	}

	for _, tc := range cases {
		graph := mustBuildGraph(t, tc.size)
		safe := CellRef{FaceFront, tc.size / 2, tc.size / 2}

		for seed := int64(1); seed <= 10; seed++ {
			mines := placeMines(graph, tc.size, tc.mines, safe, rand.New(rand.NewSource(seed)))
			if len(mines) != tc.mines {
				t.Errorf("size %d seed %d: placed %d mines, want %d", tc.size, seed, len(mines), tc.mines)
			}
		}
	}
}

func TestPlaceMinesDeterministicWithSeed(t *testing.T) {
	size := 7
	graph := mustBuildGraph(t, size)
	safe := CellRef{FaceRight, 3, 3}

	first := placeMines(graph, size, 40, safe, rand.New(rand.NewSource(42)))
	second := placeMines(graph, size, 40, safe, rand.New(rand.NewSource(42)))

	if len(first) != len(second) {
		t.Fatalf("same seed produced %d and %d mines", len(first), len(second))
	}
	for cell := range first {
		if !second[cell] {
			t.Errorf("same seed produced different layouts: %+v missing from the second", cell)
		}
	}
}

func TestPlaceMinesStaysWithinBoard(t *testing.T) {
	size := 5
	graph := mustBuildGraph(t, size)
	safe := CellRef{FaceBottom, 1, 1}

	mines := placeMines(graph, size, 15, safe, rand.New(rand.NewSource(7)))

	for cell := range mines {
		if cell.Face < FaceTop || cell.Face > FaceRight {
			t.Errorf("mine on invalid face: %+v", cell)
		}
		if cell.Row < 0 || cell.Row >= size || cell.Col < 0 || cell.Col >= size {
			t.Errorf("mine outside the board: %+v", cell)
		}
	}
}

func TestPlaceMinesCappedByAvailableCells(t *testing.T) {
	size := 2
	graph := mustBuildGraph(t, size)
	safe := CellRef{FaceFront, 0, 0}

	available := FaceCount*size*size - 1 - len(graph[safe])
	mines := placeMines(graph, size, 1000, safe, rand.New(rand.NewSource(1)))

	if len(mines) != available {
		t.Errorf("placed %d mines, want the %d available cells", len(mines), available)
	}
	if mines[safe] {
		t.Error("opened cell was mined even when every other cell was needed")
	}
}
