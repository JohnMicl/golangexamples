package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	//ChannelWithNoBuff()
	//ChannelWithBuff()
}

// channel
func ChannelWithNoBuff() {
	stream := make(chan interface{})
	now := time.Now()
	go func() {
		time.Sleep(1 * time.Second)
		stream <- struct{}{}
	}()
	<-stream
	fmt.Printf("Received channel time: %v\n", time.Since(now))
}

func ChannelWithBuff() {
	stream := make(chan int, 4)
	go func() {
		defer close(stream)
		defer fmt.Println("Produce One")
		for i := 0; i < 5; i++ {
			stream <- i
		}
	}()
	for k := range stream {
		fmt.Printf("Received channel %v\n", k)
	}
}

// package sync
func SyncPool() {
	mypool := sync.Pool{
		New: func() interface{} {
			fmt.Println("creating new instance.")
			return struct{}{}
		},
	}
	mypool.Get()
	instance := mypool.Get()
	mypool.Put(instance)
	mypool.Get()
}

func SyncOnce() {
	var count int
	increment := func() {
		count++
	}

	var once sync.Once

	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			//increment()
			once.Do(increment)
		}()
	}
	wg.Wait()
	fmt.Println(count)
}

func SyncBroadcast() {
	type Button struct {
		Clicked *sync.Cond
	}
	button := Button{Clicked: sync.NewCond(&sync.Mutex{})}
	subscribe := func(c *sync.Cond, fn func()) {
		var goroutineRunning sync.WaitGroup
		goroutineRunning.Add(1)
		go func() {
			goroutineRunning.Done()
			c.L.Lock()
			defer c.L.Unlock()
			c.Wait()
			fn()
		}()
		goroutineRunning.Wait()
	}

	var clickRegisterd sync.WaitGroup
	clickRegisterd.Add(3)

	subscribe(button.Clicked, func() {
		fmt.Println("maximum window")
		clickRegisterd.Done()
	})
	subscribe(button.Clicked, func() {
		fmt.Println("display window")
		clickRegisterd.Done()
	})
	subscribe(button.Clicked, func() {
		fmt.Println("mouse window")
		clickRegisterd.Done()
	})
	button.Clicked.Broadcast()
	clickRegisterd.Wait()
}
