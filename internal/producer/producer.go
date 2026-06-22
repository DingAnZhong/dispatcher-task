package producer

import (
	"math/rand/v2"
	"time"

	"github.com/DingAnZhong/dispatcher-task/internal/entity"
)

// 创建订单，忽略创建订单时间，假设均匀分布
func CreateOrders(orderCount int, timeLimit time.Duration) <-chan entity.Order {
	orderChan := make(chan entity.Order, orderCount)
	go func() {
		defer close(orderChan)
		for i := range orderCount {
			orderChan <- entity.NewOrder(i, entity.Position{
				X: rand.IntN(10000),
				Y: rand.IntN(10000),
			})
		}
	}()
	return orderChan
}

// 创建骑手
func CreateRiders(riderCount int) []*entity.Rider {
	riderList := make([]*entity.Rider, riderCount)
	for i := range riderCount {
		riderList[i] = entity.NewRider(i, entity.Position{
			X: rand.IntN(10000),
			Y: rand.IntN(10000),
		})
	}
	return riderList
}
