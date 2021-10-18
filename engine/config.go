package engine

import (
	"time"
)

/*
Config engine配置
*/
type Config struct {
	ExecuteTime time.Duration
	Producer    Producer
	Consumer    Consumer
}

type ErrorParser interface {
	Parser(err error)
}
