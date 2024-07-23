package main

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/marusama/cyclicbarrier"
)

func main() {
	cnt := 0
	b := cyclicbarrier.NewWithAction(10, func() error {
		cnt++
		return nil
	})

	wg := sync.WaitGroup{}
	wg.Add(10)

	for i := 0; i < 10; i++ {
		tmp := i
		go func() {
			defer wg.Done()
			for j := 0; j < 5; j++ {
				time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
				log.Printf("goroutine %d: comes %d\n", tmp, j)
				err := b.Await(context.TODO())
				log.Printf("goroutine %d: break %d\n", tmp, j)
				if err != nil {
					panic(err)
				}
				println("Counter incremented by", cnt)
			}
		}()
	}

	wg.Wait()
}
