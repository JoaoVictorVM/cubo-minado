package game

import (
	"errors"
	"testing"
)

var gameBoardSizes = []int{5, 7, 9}

func mustBuildGraph(t *testing.T, size int) Graph {
	t.Helper()
	graph, err := BuildAdjacencyGraph(size)
	if err != nil {
		t.Fatalf("BuildAdjacencyGraph(%d) returned error: %v", size, err)
	}
	return graph
}

func contains(cells []CellRef, target CellRef) bool {
	for _, cell := range cells {
		if cell == target {
			return true
		}
	}
	return false
}

func countOnFace(cells []CellRef, face FaceID) int {
	total := 0
	for _, cell := range cells {
		if cell.Face == face {
			total++
		}
	}
	return total
}

func TestEdgeTableCoversEveryFaceSideOnce(t *testing.T) {
	if len(edgeTable) != 12 {
		t.Fatalf("edge table has %d entries, want 12", len(edgeTable))
	}

	used := make(map[faceSide]int)
	for _, link := range edgeTable {
		if link.A == link.B {
			t.Errorf("edge %+v links a face side to itself", link)
		}
		used[link.A]++
		used[link.B]++
	}

	allFaces := []FaceID{FaceTop, FaceBottom, FaceFront, FaceBack, FaceLeft, FaceRight}
	allSides := []Side{SideTop, SideBottom, SideLeft, SideRight}

	for _, face := range allFaces {
		for _, side := range allSides {
			switch got := used[faceSide{face, side}]; got {
			case 1:
			case 0:
				t.Errorf("face %d side %d is not mapped by any edge", face, side)
			default:
				t.Errorf("face %d side %d is mapped by %d edges, want 1", face, side, got)
			}
		}
	}

	if len(used) != 24 {
		t.Errorf("edge table references %d distinct face sides, want 24", len(used))
	}
}

func TestEdgeTableLookupIsBidirectional(t *testing.T) {
	for _, link := range edgeTable {
		partner, reversed, ok := partnerSide(link.A)
		if !ok || partner != link.B || reversed != link.Reversed {
			t.Errorf("lookup from A of %+v gave (%+v, %t, %t)", link, partner, reversed, ok)
		}

		partner, reversed, ok = partnerSide(link.B)
		if !ok || partner != link.A || reversed != link.Reversed {
			t.Errorf("lookup from B of %+v gave (%+v, %t, %t)", link, partner, reversed, ok)
		}
	}
}

func TestFaceIndicesMatchCanonicalOrder(t *testing.T) {
	want := map[FaceID]int{
		FaceTop:    0,
		FaceBottom: 1,
		FaceFront:  2,
		FaceBack:   3,
		FaceLeft:   4,
		FaceRight:  5,
	}
	for face, index := range want {
		if int(face) != index {
			t.Errorf("face constant %d has index %d, want %d", face, int(face), index)
		}
	}
}

func TestBuildAdjacencyGraphCellCount(t *testing.T) {
	for _, size := range gameBoardSizes {
		graph := mustBuildGraph(t, size)
		want := FaceCount * size * size
		if len(graph) != want {
			t.Errorf("graph for size %d has %d cells, want %d", size, len(graph), want)
		}

		for face := 0; face < FaceCount; face++ {
			for row := 0; row < size; row++ {
				for col := 0; col < size; col++ {
					cell := CellRef{FaceID(face), row, col}
					if _, ok := graph[cell]; !ok {
						t.Fatalf("graph for size %d is missing cell %+v", size, cell)
					}
				}
			}
		}
	}
}

func TestBuildAdjacencyGraphRejectsInvalidSize(t *testing.T) {
	for _, size := range []int{1, 0, -1} {
		graph, err := BuildAdjacencyGraph(size)
		if graph != nil {
			t.Errorf("BuildAdjacencyGraph(%d) returned a graph", size)
		}
		if !errors.Is(err, ErrInvalidBoardSize) {
			t.Errorf("BuildAdjacencyGraph(%d) returned %v, want %v", size, err, ErrInvalidBoardSize)
		}
	}
}

func TestInteriorCellHasEightSameFaceNeighbors(t *testing.T) {
	size := 9
	graph := mustBuildGraph(t, size)

	for face := 0; face < FaceCount; face++ {
		for row := 1; row < size-1; row++ {
			for col := 1; col < size-1; col++ {
				cell := CellRef{FaceID(face), row, col}
				neighbors := graph[cell]
				if len(neighbors) != 8 {
					t.Fatalf("interior cell %+v has %d neighbors, want 8", cell, len(neighbors))
				}
				if countOnFace(neighbors, cell.Face) != 8 {
					t.Errorf("interior cell %+v has neighbors off its own face: %+v", cell, neighbors)
				}
			}
		}
	}
}

func TestNonCornerEdgeCellHasEightNeighbors(t *testing.T) {
	size := 5
	graph := mustBuildGraph(t, size)
	middle := size / 2

	cases := []struct {
		cell    CellRef
		crossed Side
	}{
		{CellRef{FaceFront, 0, middle}, SideTop},
		{CellRef{FaceFront, size - 1, middle}, SideBottom},
		{CellRef{FaceFront, middle, 0}, SideLeft},
		{CellRef{FaceFront, middle, size - 1}, SideRight},
		{CellRef{FaceTop, 0, middle}, SideTop},
		{CellRef{FaceBottom, size - 1, middle}, SideBottom},
		{CellRef{FaceLeft, middle, 0}, SideLeft},
		{CellRef{FaceBack, middle, size - 1}, SideRight},
	}

	for _, tc := range cases {
		neighbors := graph[tc.cell]
		if len(neighbors) != 8 {
			t.Errorf("edge cell %+v has %d neighbors, want 8", tc.cell, len(neighbors))
			continue
		}

		partner, _, ok := partnerSide(faceSide{tc.cell.Face, tc.crossed})
		if !ok {
			t.Fatalf("no partner for face %d side %d", tc.cell.Face, tc.crossed)
		}

		if got := countOnFace(neighbors, tc.cell.Face); got != 5 {
			t.Errorf("edge cell %+v has %d same-face neighbors, want 5", tc.cell, got)
		}
		if got := countOnFace(neighbors, partner.Face); got != 3 {
			t.Errorf("edge cell %+v has %d neighbors on face %d, want 3", tc.cell, got, partner.Face)
		}
	}
}

func cornerCells(size int) []CellRef {
	last := size - 1
	corners := make([]CellRef, 0, 24)
	for face := 0; face < FaceCount; face++ {
		for _, rc := range [][2]int{{0, 0}, {0, last}, {last, 0}, {last, last}} {
			corners = append(corners, CellRef{FaceID(face), rc[0], rc[1]})
		}
	}
	return corners
}

func TestCornerCellHasSevenNeighbors(t *testing.T) {
	for _, size := range gameBoardSizes {
		graph := mustBuildGraph(t, size)
		corners := cornerCells(size)

		if len(corners) != 24 {
			t.Fatalf("expected 24 corner cells, built %d", len(corners))
		}

		for _, corner := range corners {
			neighbors := graph[corner]
			if len(neighbors) != 7 {
				t.Errorf("size %d: corner %+v has %d neighbors, want 7", size, corner, len(neighbors))
				continue
			}

			if got := countOnFace(neighbors, corner.Face); got != 3 {
				t.Errorf("size %d: corner %+v has %d same-face neighbors, want 3", size, corner, got)
			}

			perFace := make(map[FaceID]int)
			for _, neighbor := range neighbors {
				if neighbor.Face != corner.Face {
					perFace[neighbor.Face]++
				}
			}
			if len(perFace) != 2 {
				t.Errorf("size %d: corner %+v touches %d other faces, want 2", size, corner, len(perFace))
			}
			for face, count := range perFace {
				if count != 2 {
					t.Errorf("size %d: corner %+v has %d neighbors on face %d, want 2", size, corner, count, face)
				}
			}
		}
	}
}

func vertexTriples(size int) [][3]CellRef {
	last := size - 1
	return [][3]CellRef{
		{{FaceFront, 0, 0}, {FaceTop, last, 0}, {FaceLeft, 0, last}},
		{{FaceFront, 0, last}, {FaceTop, last, last}, {FaceRight, 0, 0}},
		{{FaceBack, 0, last}, {FaceTop, 0, 0}, {FaceLeft, 0, 0}},
		{{FaceBack, 0, 0}, {FaceTop, 0, last}, {FaceRight, 0, last}},
		{{FaceFront, last, 0}, {FaceBottom, 0, 0}, {FaceLeft, last, last}},
		{{FaceFront, last, last}, {FaceBottom, 0, last}, {FaceRight, last, 0}},
		{{FaceBack, last, last}, {FaceBottom, last, 0}, {FaceLeft, last, 0}},
		{{FaceBack, last, 0}, {FaceBottom, last, last}, {FaceRight, last, last}},
	}
}

func TestAllEightVerticesResolveConsistently(t *testing.T) {
	for _, size := range gameBoardSizes {
		graph := mustBuildGraph(t, size)
		triples := vertexTriples(size)

		if len(triples) != 8 {
			t.Fatalf("expected 8 vertices, described %d", len(triples))
		}

		seen := make(map[CellRef]bool)
		for i, triple := range triples {
			for _, cell := range triple {
				if seen[cell] {
					t.Errorf("cell %+v appears in more than one vertex group", cell)
				}
				seen[cell] = true
			}

			for _, a := range triple {
				for _, b := range triple {
					if a == b {
						continue
					}
					if !contains(graph[a], b) {
						t.Errorf("size %d vertex %d: %+v does not list %+v as a neighbor", size, i, a, b)
					}
				}
			}
		}

		if len(seen) != 24 {
			t.Errorf("the 8 vertex groups cover %d distinct cells, want 24", len(seen))
		}
	}
}

func TestAdjacencySymmetryExhaustive(t *testing.T) {
	for _, size := range gameBoardSizes {
		graph := mustBuildGraph(t, size)
		for cell, neighbors := range graph {
			for _, neighbor := range neighbors {
				back, ok := graph[neighbor]
				if !ok {
					t.Fatalf("size %d: %+v lists %+v, which is not a cell in the graph", size, cell, neighbor)
				}
				if !contains(back, cell) {
					t.Errorf("size %d: %+v lists %+v, but not the other way around", size, cell, neighbor)
				}
			}
		}
	}
}

func TestNoDuplicateNeighbors(t *testing.T) {
	for _, size := range gameBoardSizes {
		graph := mustBuildGraph(t, size)
		for cell, neighbors := range graph {
			seen := make(map[CellRef]bool, len(neighbors))
			for _, neighbor := range neighbors {
				if neighbor == cell {
					t.Errorf("size %d: %+v lists itself as a neighbor", size, cell)
				}
				if seen[neighbor] {
					t.Errorf("size %d: %+v lists %+v more than once", size, cell, neighbor)
				}
				seen[neighbor] = true
			}
		}
	}
}

func TestEveryCellHasSevenOrEightNeighbors(t *testing.T) {
	for _, size := range gameBoardSizes {
		graph := mustBuildGraph(t, size)
		counts := make(map[int]int)
		for _, neighbors := range graph {
			counts[len(neighbors)]++
		}

		for count := range counts {
			if count != 7 && count != 8 {
				t.Errorf("size %d: found cells with %d neighbors, want only 7 or 8", size, count)
			}
		}
		if counts[7] != 24 {
			t.Errorf("size %d: %d cells have 7 neighbors, want exactly 24", size, counts[7])
		}
		if want := FaceCount*size*size - 24; counts[8] != want {
			t.Errorf("size %d: %d cells have 8 neighbors, want %d", size, counts[8], want)
		}
	}
}

func TestNeighborsAreAlwaysWithinBoard(t *testing.T) {
	for _, size := range gameBoardSizes {
		graph := mustBuildGraph(t, size)
		for cell, neighbors := range graph {
			for _, neighbor := range neighbors {
				if neighbor.Face < FaceTop || neighbor.Face > FaceRight {
					t.Fatalf("size %d: %+v lists neighbor on invalid face %d", size, cell, neighbor.Face)
				}
				if neighbor.Row < 0 || neighbor.Row >= size || neighbor.Col < 0 || neighbor.Col >= size {
					t.Fatalf("size %d: %+v lists out-of-board neighbor %+v", size, cell, neighbor)
				}
			}
		}
	}
}

func TestGraphIsDeterministic(t *testing.T) {
	for _, size := range gameBoardSizes {
		first := mustBuildGraph(t, size)
		second := mustBuildGraph(t, size)

		for cell, neighbors := range first {
			other := second[cell]
			if len(neighbors) != len(other) {
				t.Fatalf("size %d: %+v has %d neighbors on rebuild, want %d", size, cell, len(other), len(neighbors))
			}
			for i := range neighbors {
				if neighbors[i] != other[i] {
					t.Errorf("size %d: %+v neighbor %d differs on rebuild: %+v vs %+v",
						size, cell, i, neighbors[i], other[i])
				}
			}
		}
	}
}

func TestBuildAdjacencyGraphSmallestValidSize(t *testing.T) {
	graph := mustBuildGraph(t, 2)
	if len(graph) != FaceCount*4 {
		t.Fatalf("graph for size 2 has %d cells, want %d", len(graph), FaceCount*4)
	}

	for cell, neighbors := range graph {
		if len(neighbors) != 7 {
			t.Errorf("size 2: every cell is a corner, but %+v has %d neighbors, want 7", cell, len(neighbors))
		}
		seen := make(map[CellRef]bool, len(neighbors))
		for _, neighbor := range neighbors {
			if seen[neighbor] {
				t.Errorf("size 2: %+v lists %+v more than once", cell, neighbor)
			}
			seen[neighbor] = true
		}
	}
}

func TestPartnerSideRejectsUnknownSide(t *testing.T) {
	partner, reversed, ok := partnerSide(faceSide{FaceFront, Side(99)})
	if ok {
		t.Errorf("partnerSide accepted an unknown side, returning %+v", partner)
	}
	if reversed {
		t.Error("partnerSide reported a reversal for an unknown side")
	}
	if partner != (faceSide{}) {
		t.Errorf("partnerSide returned %+v for an unknown side, want the zero value", partner)
	}
}
