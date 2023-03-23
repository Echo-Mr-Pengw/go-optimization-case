// 单个生产者，多个消费者
package main

import (
	"fmt"
	"sync"
)

// 生产者
func producer(wg *sync.WaitGroup, ch chan int) {
	defer wg.Done()
	for i := 1; i <= 10; i++ {
		ch <- i
	}
	close(ch)
}

// 消费者
func consumer(wg *sync.WaitGroup, ch chan int, consumerGoroutineNum int) {
	defer wg.Done()
	for v := range ch {
		fmt.Printf("consumer=%d, data=%d\n", consumerGoroutineNum, v)
	}
}

func main() {

	ch := make(chan int, 10)

	var wg sync.WaitGroup
	wg.Add(10) // 1个producer 9个consumer

	go producer(&wg, ch)

	for i := 1; i < 10; i++ {
		go consumer(&wg, ch, i)
	}
	wg.Wait()
}
