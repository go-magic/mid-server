package engine

import (
	"time"
)

type Config struct {
	SuccessWaitTime time.Duration
	ErrorWaitTime   time.Duration
	ExecuteTime     time.Duration
	Producer        Producer
	Consumer        Consumer
}
