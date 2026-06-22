package entity

import (
	"sync"
)

type ReadyQueue struct {
	orders []*Order
	mu     sync.Mutex
	cond   *sync.Cond
	closed bool
}

func NewReadyQueue() *ReadyQueue {
	rq := &ReadyQueue{
		orders: make([]*Order, 0, 1024),
	}
	rq.cond = sync.NewCond(&rq.mu)
	return rq
}

// 插入订单
func (rq *ReadyQueue) PushOrder(o *Order) {
	rq.mu.Lock()
	rq.orders = append(rq.orders, o)
	rq.mu.Unlock()
	rq.cond.Signal()
}

// 取出订单，若队列为空且未关闭则阻塞，关闭后返回 nil
func (rq *ReadyQueue) PopOrder() *Order {
	rq.mu.Lock()
	defer rq.mu.Unlock()
	for len(rq.orders) == 0 && !rq.closed {
		rq.cond.Wait()
	}
	if len(rq.orders) == 0 {
		return nil
	}
	o := rq.orders[0]
	rq.orders = rq.orders[1:]
	return o
}

// 关闭队列，唤醒所有等待的协程
func (rq *ReadyQueue) Close() {
	rq.mu.Lock()
	rq.closed = true
	rq.mu.Unlock()
	rq.cond.Broadcast()
}
