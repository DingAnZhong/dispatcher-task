package config

import "time"

type Config struct {
	Concurrency int
	RiderCount  int
	OrderCount  int
	TimeLimit   time.Duration
}
