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
	// sync.Pool 是 Go 语言标准库 sync 包中的一个类型，它提供了一个内存对象池的实现。
	// sync.Pool 主要用于存储临时对象，以便在后续需要时重用这些对象，从而避免频繁的内存分配和垃圾回收带来的性能开销。
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
	// sync.Pool 内部的管理是基于运行时机制的，它不保证对象的存活时间，系统可能在任何时间清理 pool 中的对象，特别是当内存压力较大时。
	// 因此，不应该依赖 pool 存储重要的数据，尤其是那些需要长期存在的数据。同时，存入 pool 的对象应当在其生命周期结束前确保其状态可以安全地复用。
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
	// sync.Cond 是 Go 语言标准库 sync 包中的一个并发原语，用于在多个 goroutine 之间进行同步和通信。
	// 它是一个条件变量，允许 goroutine 在某个条件满足时阻塞等待，并在条件变为真时被唤醒。
	type Button struct {
		Clicked *sync.Cond
	}

	// 使用 sync.Cond 需要配合一个互斥锁（通常是 sync.Mutex 或 sync.RWMutex），
	// 因为条件变量的 Wait() 方法需要在锁定状态下调用，而 Signal() 和 Broadcast() 方法则用来通知等待的 goroutine。
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

	// 如果有多于一个的 goroutine 在等待同一个条件变量，可以使用 cond.Broadcast() 来唤醒所有等待的 goroutine。
	button.Clicked.Broadcast()
	clickRegisterd.Wait()
}
