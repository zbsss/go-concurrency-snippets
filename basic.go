package main

import (
	"fmt"
	"sync"
	"time"
)

func rangingOverChannel() {
	fmt.Println("---ranging over channel---")
	intStream := make(chan int)
	go func() {
		defer close(intStream)
		for i := 0; i < 5; i++ {
			intStream <- i
		}
	}()

	// iteration stops when channel is closed
	for i := range intStream {
		fmt.Println(i)
	}
}

func unblockManyGoRoutinesWithClosingChannel() {
	fmt.Println("---waiting for begin---")
	begin := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-begin
			fmt.Printf("%d has begun\n", i)
		}(i)
	}
	close(begin)
	wg.Wait()
}

func writeToChannelFromSelect() {
	fmt.Println("---write to channel from select---")
	c := make(chan interface{})

	go func() {
		select {
		case c <- struct{}{}:
			// blocks until it's possible to write to the channel
			fmt.Println("unblocked and written to channel")
			return
		}
	}()

	<-c
}

func timeout() {
	fmt.Println("---timeout in select---")

	c := make(chan int)
	select {
	case <-c:
	case <-time.After(1 * time.Second):
		fmt.Println("timed out")
	}
}
