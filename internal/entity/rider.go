package entity

import (
	"sync"
	"sync/atomic"
)

type Rider struct {
	RID int
	Pos Position
	// Orders []*Order //每个骑手维护一个订单表
	Orders []*Order
	Cap    atomic.Int64
	mu     sync.Mutex
}

// 简单工厂模式
func NewRider(rid int, pos Position) *Rider {
	return &Rider{
		RID:    rid,
		Pos:    pos,
		Orders: make([]*Order, 0),
	}
}

// 有锁加订单
func (r *Rider) AddOrder(order *Order) {
	r.mu.Lock()
	r.Orders = append(r.Orders, order)
	r.mu.Unlock()
}
