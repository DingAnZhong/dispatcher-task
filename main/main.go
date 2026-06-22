package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"github.com/DingAnZhong/dispatcher-task/internal/config"
	"github.com/DingAnZhong/dispatcher-task/internal/dispatcher"
	"github.com/DingAnZhong/dispatcher-task/internal/producer"
)

func main() {
	test(config.Config{
		Concurrency: 2,
		RiderCount:  100,
		OrderCount:  10000,
		TimeLimit:   30,
	}) // 分配时间:0.0012399s(无订单切片) ||| 分配时间:0.0021185s(有订单切片)
	test(config.Config{
		Concurrency: 2,
		RiderCount:  1000,
		OrderCount:  100000,
		TimeLimit:   60,
	}) // 分配时间:0.0182922s(无订单切片) ||| 分配时间:0.0193193s(有订单切片)
	test(config.Config{
		Concurrency: 2,
		RiderCount:  10000,
		OrderCount:  1000000,
		TimeLimit:   180,
	}) // 分配时间:0.2151369s(无订单切片) ||| 分配时间:0.2295283s(有订单切片)
	test(config.Config{
		Concurrency: 2,
		RiderCount:  100000,
		OrderCount:  10000000,
		TimeLimit:   600,
	}) // 分配时间:4.8138713s(无订单切片) ||| 分配时间:5.1135213s(有订单切片)
}

func test(cfg config.Config) {

	//1、创建骑手
	riderList := producer.CreateRiders(cfg.RiderCount)
	//2、创建订单
	orderChan := producer.CreateOrders(cfg.OrderCount, cfg.TimeLimit*time.Second)
	//3、创建调度器
	dispatcher := dispatcher.NewDispatcher(cfg, riderList, orderChan)
	//4、分配订单
	fmt.Printf("参数:骑手数量%d,订单数%d,时限%v,并发数%d\n", cfg.RiderCount, cfg.OrderCount, cfg.TimeLimit*time.Second, cfg.Concurrency)
	// cup profile
	cpuFile, err := os.Create("cpu.prof")
	if err != nil {
		fmt.Println("Create cpu.prof filed:", err)
		return
	}
	defer cpuFile.Close()
	err = pprof.StartCPUProfile(cpuFile)
	if err != nil {
		fmt.Println("start cpu.prof filed:", err)
		return
	}
	dispatcher.Dispatch()
	pprof.StopCPUProfile()
	//5、查看末尾
	dispatcher.ShowTail10()
	// heap profile
	// 无骑手订单~30MB  (极端情况)
	// 有骑手订单~500MB (极端情况)
	heapFlie, err := os.Create("heap.prof")
	if err != nil {
		fmt.Println("Create heap.prof filed:", err)
		return
	}
	defer heapFlie.Close()
	if err := pprof.WriteHeapProfile(heapFlie); err != nil {
		fmt.Println("write heap profile failed:", err)
	}
	fmt.Println("heap.prof saved")
}
