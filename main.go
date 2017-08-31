package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

var wg sync.WaitGroup

func threadLoop(counter *int32) {
	id := atomic.AddInt32(counter, -1)
	// цикл не позволяет отпустить текущий поток,
	// планировщик рутин должен забрать с этого потока остальные рутины
	for atomic.LoadInt32(counter) != 0 {
	}
	// если вышли из цикла, точно запустились по рутине на поток

	// прибиваем рутину к потоку
	runtime.LockOSThread()

	// без привязки к потоку рутина при вводе.выводе могла сменить поток
	fmt.Printf("thread %d locked\n", id)

	// тут какой-то целевой код, цикл на сокеты и т.д.
	// ...

	wg.Done()
}

func main() {
	// Для main и так выполняется runtime.LockOSThread()

	n := runtime.NumCPU()
	counter := int32(n)

	wg.Add(n)

	for i := 0; i < n-1; i++ {
		go threadLoop(&counter)
	}
	threadLoop(&counter)

	wg.Wait() // ждем остальные потоки

	fmt.Println("done")
}
