package entity

import "time"

type Order struct {
	OID         int
	Pos         Position
	createdTime time.Time // 时间越早优先级越高
}

// 简单工厂模式
func NewOrder(oid int, pos Position) Order {
	return Order{
		OID:         oid,
		Pos:         pos,
		createdTime: time.Now(),
	}
}
