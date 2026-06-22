package dispatcher

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/DingAnZhong/dispatcher-task/internal/config"
	"github.com/DingAnZhong/dispatcher-task/internal/entity"
)

// 调度器，封装骑手、订单通道和配置
type Dispatcher struct {
	cfg        config.Config
	riderList  []*entity.Rider
	orderChan  <-chan entity.Order
	readyQueue *entity.ReadyQueue
	wg         sync.WaitGroup
}

// 创建一个新的调度器
func NewDispatcher(cfg config.Config, riders []*entity.Rider, orders <-chan entity.Order) *Dispatcher {
	return &Dispatcher{
		cfg:        cfg,
		riderList:  riders,
		orderChan:  orders,
		readyQueue: entity.NewReadyQueue(),
		wg:         sync.WaitGroup{},
	}
}

// 任务入队，削峰(time.sleep或者MQ，这里无实现)
func (d *Dispatcher) push2rq() {
	for order := range d.orderChan {
		o := order // 不这样循环变量order一直指向同一地址
		d.readyQueue.PushOrder(&o)
	}
	// time.Sleep(1000) // sleep 1000 ns 防止DB被打爆
	d.readyQueue.Close()
}

// 开始调度，使用网格索引进行最近骑手分配，统计耗时并输出
func (d *Dispatcher) Dispatch() {
	// 启动一个协程入队
	go d.push2rq()
	concurrency := d.cfg.Concurrency
	if concurrency <= 0 {
		concurrency = 1
	}

	gridIndex := NewGrid(d.riderList)
	startTime := time.Now()

	d.wg.Add(concurrency)

	// 协程池
	for range concurrency {
		go func() {
			defer d.wg.Done()
			for {
				order := d.readyQueue.PopOrder()
				if order == nil {
					return
				}
				d.assignOrder(order, gridIndex)
			}
		}()
	}

	d.wg.Wait()
	elapsed := time.Since(startTime).Seconds()
	fmt.Printf("分配时间:%vs\n", elapsed)
}

// 订单找骑手，周围格子没有骑手话降级
func (d *Dispatcher) assignOrder(order *entity.Order, grid [][]Grid) {
	minDist := math.MaxInt
	var bestRider *entity.Rider

	cellX := order.Pos.X / 100
	cellY := order.Pos.Y / 100

	// 网格搜索（5x5 格子）总的10000个格子，1/400
	for x := cellX - 2; x <= cellX+2; x++ {
		for y := cellY - 2; y <= cellY+2; y++ {
			if x >= 0 && x < len(grid) && y >= 0 && y < len(grid[0]) {
				for _, rider := range grid[x][y].riderList {
					dist := entity.DistanceSquire(order.Pos, rider.Pos)
					if dist < minDist {
						minDist = dist
						bestRider = rider
					}
				}
			}
		}
	}

	// 降级为全局搜索
	if bestRider == nil {
		for _, rider := range d.riderList {
			dist := entity.DistanceSquire(order.Pos, rider.Pos)
			if dist < minDist {
				minDist = dist
				bestRider = rider
			}
		}
	}

	if bestRider != nil {
		bestRider.Cap.Add(1)
	}
}

// 输出接单数最少的10名骑手
func (d *Dispatcher) ShowTail10() {
	tail := findLeastOrders(d.riderList, 10)
	fmt.Println("接单数最少的10名骑手：")
	for _, rider := range tail {
		fmt.Printf("骑手ID:%d, 订单数:%d\n", rider.RID, rider.Cap.Load())
	}
}

// 订单数最小的 k 个，辅助函数
func findLeastOrders(riderList []*entity.Rider, k int) []*entity.Rider {
	if len(riderList) <= k {
		return riderList
	}
	heap := make([]*entity.Rider, k)
	copy(heap, riderList[:k])
	for i := k/2 - 1; i >= 0; i-- {
		maxHeapify(heap, k, i)
	}
	for i := k; i < len(riderList); i++ {
		if riderList[i].Cap.Load() < heap[0].Cap.Load() {
			heap[0] = riderList[i]
			maxHeapify(heap, k, 0)
		}
	}
	return heap
}

// 最大堆，辅助函数
func maxHeapify(arr []*entity.Rider, n, i int) {
	largest := i
	left := 2*i + 1
	right := 2*i + 2

	if left < n && arr[left].Cap.Load() > arr[largest].Cap.Load() {
		largest = left
	}
	if right < n && arr[right].Cap.Load() > arr[largest].Cap.Load() {
		largest = right
	}
	if largest != i {
		arr[i], arr[largest] = arr[largest], arr[i]
		maxHeapify(arr, n, largest)
	}
}
