package dispatcher

import "github.com/DingAnZhong/dispatcher-task/internal/entity"

type Grid struct {
	riderList []*entity.Rider
}

// 建立骑手的网格索引，加一圈padding
func NewGrid(riderList []*entity.Rider) [][]Grid {
	res := make([][]Grid, 102)
	for i := range 102 {
		res[i] = make([]Grid, 102)
	}
	for _, rider := range riderList {
		cellX := rider.Pos.X/100 + 1
		cellY := rider.Pos.Y/100 + 1
		res[cellX][cellY].riderList = append(res[cellX][cellY].riderList, rider)
	}
	return res
}
