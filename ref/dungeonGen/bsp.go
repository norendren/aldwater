package dungeonGen

import "math/rand"

type MapNode struct {
	Left   *MapNode
	Right  *MapNode
	X      int
	Y      int
	Width  int
	Height int
}

func (m *MapNode) Split() bool {
	if (m.Left == nil) || (m.Right == nil) {
		return false
	}
	if rand.Intn(1) == 0 {
		// Vertical

	}
	return false
}
