package main

import (
	"fmt"
	"math/rand"
	"time"
)

func generator() chan int {
	out := make(chan int)
	go func() {
		i := 0
		for {
			//随机sleep一段时间
			time.Sleep(
				time.Duration(rand.Intn(1500)) * time.Millisecond)
			out <- i
			i++
		}
	}()
	return out
}

func worker(id int, c chan int) {
	for n := range c {
		fmt.Printf("worker %d received %d\n", id, n)
	}
}

func createWorker(id int) chan<- int {
	c := make(chan int)
	go worker(id, c)
	return c
}

func main() {
	//c1, c2被赋值的channel
	var c1, c2 = generator(), generator()

	// w := createWorker(0)
	//worker worker0的channel
	var worker = createWorker(0)
	n := 0
	hasValue := false
	for {
		var activeWorker chan<- int //activeWorker是nilchannel
		if hasValue {
			activeWorker = worker

		}
		select {
		case n = <-c1:
			hasValue = true
		case n = <-c2:
			hasValue = true
		//case w <- n:,这个case会阻塞，
		//在n还没放入数据的时候
		//为了解决这个问题，加入activeworker
		//activeWorker是nil的，会block，因为select的避免永久等待机制，case不被select到（这个case被ignore）
		case activeWorker <- n:
			hasValue = false
		}
	}
}
