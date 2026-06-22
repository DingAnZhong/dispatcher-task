package entity

import "sync/atomic"

type Rider struct {
	RID int
	Pos Position
	// Orders []Order //每个骑手维护一个订单表
	Cap atomic.Int64
}

// 简单工厂模式
func NewRider(rid int, pos Position) *Rider {
	return &Rider{
		RID: rid,
		Pos: pos,
	}
}
