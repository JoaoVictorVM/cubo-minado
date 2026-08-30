package game

type FaceID int

const (
	FaceTop FaceID = iota
	FaceBottom
	FaceFront
	FaceBack
	FaceLeft
	FaceRight
)

type Side int

const (
	SideTop Side = iota
	SideBottom
	SideLeft
	SideRight
)

type CellRef struct {
	Face FaceID
	Row  int
	Col  int
}

type Graph map[CellRef][]CellRef

type faceSide struct {
	Face FaceID
	Side Side
}

type edgeLink struct {
	A        faceSide
	B        faceSide
	Reversed bool
}

var edgeTable = []edgeLink{
	{faceSide{FaceFront, SideTop}, faceSide{FaceTop, SideBottom}, false},
	{faceSide{FaceFront, SideBottom}, faceSide{FaceBottom, SideTop}, false},
	{faceSide{FaceFront, SideLeft}, faceSide{FaceLeft, SideRight}, false},
	{faceSide{FaceFront, SideRight}, faceSide{FaceRight, SideLeft}, false},
	{faceSide{FaceRight, SideRight}, faceSide{FaceBack, SideLeft}, false},
	{faceSide{FaceBack, SideRight}, faceSide{FaceLeft, SideLeft}, false},
	{faceSide{FaceTop, SideLeft}, faceSide{FaceLeft, SideTop}, false},
	{faceSide{FaceTop, SideRight}, faceSide{FaceRight, SideTop}, true},
	{faceSide{FaceTop, SideTop}, faceSide{FaceBack, SideTop}, true},
	{faceSide{FaceBottom, SideLeft}, faceSide{FaceLeft, SideBottom}, true},
	{faceSide{FaceBottom, SideRight}, faceSide{FaceRight, SideBottom}, false},
	{faceSide{FaceBottom, SideBottom}, faceSide{FaceBack, SideBottom}, true},
}

func partnerSide(origin faceSide) (faceSide, bool, bool) {
	for _, link := range edgeTable {
		if link.A == origin {
			return link.B, link.Reversed, true
		}
		if link.B == origin {
			return link.A, link.Reversed, true
		}
	}
	return faceSide{}, false, false
}

func placeOnSide(side Side, index, size int) (int, int) {
	switch side {
	case SideTop:
		return 0, index
	case SideBottom:
		return size - 1, index
	case SideLeft:
		return index, 0
	default:
		return index, size - 1
	}
}

func crossedSide(newRow, newCol, size int) Side {
	if newRow < 0 {
		return SideTop
	}
	if newRow >= size {
		return SideBottom
	}
	if newCol < 0 {
		return SideLeft
	}
	return SideRight
}

func neighborsOf(cell CellRef, size int) []CellRef {
	neighbors := make([]CellRef, 0, 8)
	seen := make(map[CellRef]bool, 8)

	add := func(candidate CellRef) {
		if candidate != cell && !seen[candidate] {
			seen[candidate] = true
			neighbors = append(neighbors, candidate)
		}
	}

	for dRow := -1; dRow <= 1; dRow++ {
		for dCol := -1; dCol <= 1; dCol++ {
			if dRow == 0 && dCol == 0 {
				continue
			}

			newRow := cell.Row + dRow
			newCol := cell.Col + dCol
			rowOut := newRow < 0 || newRow >= size
			colOut := newCol < 0 || newCol >= size

			if rowOut && colOut {
				continue
			}

			if !rowOut && !colOut {
				add(CellRef{cell.Face, newRow, newCol})
				continue
			}

			side := crossedSide(newRow, newCol, size)
			partner, reversed, ok := partnerSide(faceSide{cell.Face, side})
			if !ok {
				continue
			}

			index := newCol
			if colOut {
				index = newRow
			}
			if reversed {
				index = size - 1 - index
			}

			row, col := placeOnSide(partner.Side, index, size)
			add(CellRef{partner.Face, row, col})
		}
	}

	return neighbors
}

func BuildAdjacencyGraph(size int) (Graph, error) {
	if size < 2 {
		return nil, ErrInvalidBoardSize
	}

	graph := make(Graph, FaceCount*size*size)
	for face := 0; face < FaceCount; face++ {
		for row := 0; row < size; row++ {
			for col := 0; col < size; col++ {
				cell := CellRef{FaceID(face), row, col}
				graph[cell] = neighborsOf(cell, size)
			}
		}
	}

	return graph, nil
}
