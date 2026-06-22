package config

import (
	"fmt"
	"time"
)

type Config struct {
	Concurrency int
	RiderCount  int
	OrderCount  int
	TimeLimit   time.Duration
}

func (cfg *Config) Validate() {
	if cfg.Concurrency <= 0 {
		cfg.Concurrency = 1
		fmt.Println("并发数有误，已设为1")
	}
	if cfg.RiderCount <= 0 {
		cfg.RiderCount = 1
		fmt.Println("骑手数有误，已设为1")
	}
	if cfg.OrderCount <= 0 {
		cfg.OrderCount = 1
		fmt.Println("订单数有误，已设为1")
	}
	if cfg.TimeLimit <= 0 {
		cfg.TimeLimit = time.Minute
		fmt.Println("时间限制有误，已设为1分钟")
	}
}
