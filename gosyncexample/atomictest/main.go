package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type Config struct {
	Nodename string
	Addr     string
	Count    int32
}

func loadNewConfig() Config {
	return Config{
		Nodename: "BEIJING",
		Addr:     "10.11.16.14",
		Count:    rand.Int31(),
	}
}

func main() {

	var config atomic.Value
	config.Store(loadNewConfig())

	mu := sync.RWMutex{}
	lock := sync.NewCond(&mu)

	go func() {
		for {
			time.Sleep(time.Duration(5 + rand.Int63n(5+rand.Int63n(5))*int64(time.Second)))
			config.Store(loadNewConfig())
			lock.Signal()
		}
	}()

	for {
		mu.Lock()
		lock.Wait()
		c := config.Load().(Config)
		fmt.Printf("config is %v\n", c)
		mu.Unlock()
	}

}
