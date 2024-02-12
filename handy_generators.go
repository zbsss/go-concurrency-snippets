package main

import "fmt"

func repeat(done <-chan interface{}, values ...interface{}) <-chan interface{} {
	valueStream := make(chan interface{})

	go func() {
		defer close(valueStream)
		for {
			for _, value := range values {
				select {
				case <-done:
					return
				case valueStream <- value:
				}
			}
		}
	}()

	return valueStream
}

func take(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
	takeStream := make(chan interface{})

	go func() {
		defer close(takeStream)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:
				// pass value from one channel to another
			}
		}
	}()

	return takeStream
}

func takeRepeatExample() {
	fmt.Println("---take repeat generator pipeline---")
	done := make(chan interface{})
	defer close(done)

	for num := range take(done, repeat(done, 1, 2), 10) {
		fmt.Printf("%d ", num)
	}
}
