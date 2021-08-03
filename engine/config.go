package engine

import (
	"time"
)

/*
Config engine配置
*/
type Config struct {
	SuccessWaitTime time.Duration
	ErrorWaitTime   time.Duration
	ExecuteTime     time.Duration
	Producer        Producer
	Consumer        Consumer
}
