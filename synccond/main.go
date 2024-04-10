package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type SubTask struct {
	id         int
	mutex      *sync.Mutex
	cond       *sync.Cond
	parentTask *Task
}

func (s *SubTask) Handle() {
	fmt.Printf("subtask task id=%+v start..., parent task id=%+v\n", s.id, s.parentTask.id)
	n := rand.Intn(3)
	time.Sleep(time.Duration(n) * time.Second)
	s.mutex.Lock()
	s.parentTask.finishedSubTaskNum += 1
	s.cond.Broadcast()
	s.mutex.Unlock()
	fmt.Printf("subtask task id=%+v finished, parent task id=%+v\n", s.id, s.parentTask.id)
}

type Task struct {
	id                 int
	subTaskNum         int
	finishedSubTaskNum int
	subTaskList        []*SubTask
	mutex              *sync.Mutex
	cond               *sync.Cond
	send               chan<- *SubTask
}

func NewTask(taskId, num int, send chan<- *SubTask) *Task {
	m := &sync.Mutex{}
	c := sync.NewCond(m)
	task := &Task{
		id:                 taskId,
		subTaskNum:         num,
		finishedSubTaskNum: 0,
		subTaskList:        make([]*SubTask, num),
		mutex:              m,
		cond:               c,
		send:               send,
	}

	for i := 0; i < num; i++ {
		subTask := &SubTask{
			id:         taskId*100 + i,
			mutex:      m,
			cond:       c,
			parentTask: task,
		}
		task.subTaskList[i] = subTask
	}
	return task
}

func (t *Task) HandleTask() {
	for _, v := range t.subTaskList {
		t.send <- v
	}
}

func (t *Task) WaitForMerge() {
	t.mutex.Lock()
	for t.finishedSubTaskNum != t.subTaskNum {
		fmt.Printf("waiting for merge task info id=%+v, total=%+v, finished=%+v\n", t.id, t.subTaskNum, t.finishedSubTaskNum)
		t.cond.Wait()
	}
	t.MergeWork()
	t.mutex.Unlock()
}

func (t *Task) MergeWork() {
	fmt.Printf("finish merge task info id=%+v, total=%+v, finished=%+v\n", t.id, t.subTaskNum, t.finishedSubTaskNum)
}

type Worker struct {
	id       int
	received <-chan *SubTask
}

func (w *Worker) dowork() {
	fmt.Printf("worker %+v statr work\n", w.id)
	for {
		select {
		case subtask := <-w.received:
			subtask.Handle()
		}
	}
}

type WorkPool struct {
	worker []*Worker
}

func NewWorkPool(n int, received <-chan *SubTask) *WorkPool {
	workePool := &WorkPool{
		worker: make([]*Worker, n),
	}
	for i := 0; i < n; i++ {
		workePool.worker[i] = &Worker{
			id:       i,
			received: received,
		}
	}
	return workePool
}

func (w *WorkPool) Start() {
	for i := 0; i < len(w.worker); i++ {
		go w.worker[i].dowork()
	}
}

func NewTaskWorkChan() chan *SubTask {
	ch := make(chan *SubTask)
	return ch
}

func main() {
	newch := NewTaskWorkChan()
	pool := NewWorkPool(4, newch)
	pool.Start()

	task1 := NewTask(1, 10, newch)
	task2 := NewTask(2, 15, newch)
	task3 := NewTask(3, 12, newch)

	go task1.HandleTask()
	go task1.WaitForMerge()

	go task2.HandleTask()
	go task2.WaitForMerge()

	go task3.HandleTask()
	go task3.WaitForMerge()

	// main process not exit
	select {}
}
