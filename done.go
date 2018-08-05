package main

/**使用channel等待任务结束*/
import (
	"fmt"
	"sync"
)

func doworker(id int, c chan int, wg *sync.WaitGroup) {
	for n := range c {
		fmt.Printf("worker %d received %c\n", id, n)
		//再写一个go rountine，防止循环等待done中的数据没被拿走，又要放入数据到done中，、
		//一个等着done中数据清空，一个等着放入done数据
		wg.Done()

	}
}

func createWorker(id int, wg *sync.WaitGroup) worker {
	w := worker{
		in: make(chan int),
		wg: wg,
	}
	go doworker(id, w.in, w.wg)
	return w
}

type worker struct {
	in chan int
	wg *sync.WaitGroup
}

func chanDemo() {
	var wg sync.WaitGroup

	var workers [10]worker
	for i := 0; i < 10; i++ {
		workers[i] = createWorker(i, &wg)
	}
	//添加20个任务
	wg.Add(20)
	for i := 0; i < 10; i++ {
		workers[i].in <- 'a' + i
		//等到打印结束之后才能从done中读到东西
		//这样会按顺序打印
		// <-workers[i].done
	}

	for i := 0; i < 10; i++ {
		workers[i].in <- 'A' + i
		// <-workers[i].done
	}

	wg.Wait()
	// time.Sleep(time.Millisecond)
}

func main() {
	chanDemo()
}
