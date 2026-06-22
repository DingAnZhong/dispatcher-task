package entity

type Position struct {
	X int
	Y int
}

// 两点距离
func DistanceSquire(pos1, pos2 Position) int {
	dx := pos1.X - pos2.X
	dy := pos1.Y - pos2.Y
	return dx*dx + dy*dy
}
