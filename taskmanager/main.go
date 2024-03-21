package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	//ChannelWithNoBuff()
	//ChannelWithBuff()
	DealWithTask()
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

type TaskStatus int

const (
	TaskInit     TaskStatus = iota // 0
	TaskRuning                     // 1
	TaskFinished                   //2
	TaskAborted                    //3
)

type Task struct {
	ID     int
	Status TaskStatus
}

type TaskManagerImpl interface {
	Stop()
	Run()
	TryAddTask(task *Task)
}

type TaskManager struct {
	done          chan interface{}
	taskQueue     chan *Task
	informerQueue chan *Task
}

func NewTaskManager(workerN int) *TaskManager {
	return &TaskManager{
		done:          make(chan interface{}),
		taskQueue:     make(chan *Task, workerN),
		informerQueue: make(chan *Task, workerN),
	}
}

func (t *TaskManager) Stop() {
	defer close(t.done)
	defer close(t.taskQueue)
	defer close(t.informerQueue)
}

func (t *TaskManager) Run() {
	for i := 0; i < 3; i++ {
		go t.dowork(i)
		go t.rollbackWork()
	}
}

func (t *TaskManager) dowork(workid int) {
	for {
		select {
		case <-t.done:
			fmt.Printf("Child %+v Received Closed by main func\n", workid)
			return
		case task := <-t.taskQueue:
			fmt.Printf("Child %+v Received task %+v\n", workid, t)
			handleTask(task, t.informerQueue)
		}
	}
}

func (t *TaskManager) rollbackWork() {
	for {
		select {
		case <-t.done:
			return
		case task := <-t.informerQueue:
			// do others to handle task
			fmt.Printf("task received finished messages %v\n", task)
		}
	}
}

func (t *TaskManager) TryAddTask(task *Task) {
	// 缓冲区满了会导致调用阻塞，可在此添加超时机制
	t.taskQueue <- task
	fmt.Println("Message enqueued successfully.")
}

func DealWithTask() {
	now := time.Now()
	taskmg := NewTaskManager(4)
	taskmg.Run()

	// simulate task generation
	for i := 0; i < 10; i++ {
		task := Task{ID: i, Status: TaskInit}
		taskmg.TryAddTask(&task)
		time.Sleep(time.Millisecond * 100)
	}
	fmt.Printf("task finised end %v\n", time.Since(now))
	time.Sleep(5 * time.Second)
}

func handleTask(task *Task, informer chan<- *Task) {
	// simulate task handle
	now := time.Now()
	time.Sleep(1 * time.Second)
	switch task.Status {
	case TaskInit:
		fmt.Printf("Start Hanle task %v\n", task)
	case TaskFinished:
		fmt.Printf("Start Hanle task %v\n", task)
	case TaskRuning:
		fmt.Printf("Start Hanle task %v\n", task)
	case TaskAborted:
		fmt.Printf("Start Hanle task %v\n", task)
	}
	informer <- task
	fmt.Printf("Task %v finished informer time now = %v\n", task, time.Since(now))
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
