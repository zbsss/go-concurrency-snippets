package main

import (
	"fmt"
)

func pipelineExample() {
	fmt.Println("---pipeline---")

	generator := func(done <-chan interface{}, nums ...int) <-chan int {
		intStream := make(chan int)

		go func() {
			defer close(intStream)
			for _, num := range nums {
				select {
				case <-done:
					return
				case intStream <- num:
				}
			}
		}()

		return intStream
	}

	multiply := func(done <-chan interface{}, intStream <-chan int, multiplier int) <-chan int {
		multipliedStream := make(chan int)
		go func() {
			defer close(multipliedStream)
			for num := range intStream {
				select {
				case <-done:
					return
				case multipliedStream <- num * multiplier:
				}
			}
		}()

		return multipliedStream
	}

	add := func(done <-chan interface{}, intStream <-chan int, additive int) <-chan int {
		addedStream := make(chan int)
		go func() {
			defer close(addedStream)
			for num := range intStream {
				select {
				case <-done:
					return
				case addedStream <- num + additive:
				}
			}
		}()

		return addedStream
	}

	done := make(chan interface{})
	defer close(done)

	intStream := generator(done, 1, 2, 3, 4)

	pipeline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)

	for v := range pipeline {
		fmt.Println(v)
	}
}
