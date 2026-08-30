package game

import "math/rand"

func placeMines(graph Graph, size, totalMines int, safe CellRef, rng *rand.Rand) map[CellRef]bool {
	forbidden := map[CellRef]bool{safe: true}
	for _, neighbor := range graph[safe] {
		forbidden[neighbor] = true
	}

	candidates := make([]CellRef, 0, FaceCount*size*size)
	for face := 0; face < FaceCount; face++ {
		for row := 0; row < size; row++ {
			for col := 0; col < size; col++ {
				cell := CellRef{FaceID(face), row, col}
				if !forbidden[cell] {
					candidates = append(candidates, cell)
				}
			}
		}
	}

	if totalMines > len(candidates) {
		totalMines = len(candidates)
	}

	for i := 0; i < totalMines; i++ {
		j := i + rng.Intn(len(candidates)-i)
		candidates[i], candidates[j] = candidates[j], candidates[i]
	}

	mines := make(map[CellRef]bool, totalMines)
	for _, cell := range candidates[:totalMines] {
		mines[cell] = true
	}

	return mines
}
