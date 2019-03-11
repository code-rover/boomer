package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/myzhan/boomer"
)

func foo() {
	start := boomer.Now()
	time.Sleep(100 * time.Millisecond)
	elapsed := boomer.Now() - start

	// Report your test result as a success, if you write it in python, it will looks like this
	// events.request_success.fire(request_type="http", name="foo", response_time=100, response_length=10)
	globalBoomer.RecordSuccess("http", "foo", elapsed, int64(10))
}

func bar() {
	start := boomer.Now()
	time.Sleep(100 * time.Millisecond)
	elapsed := boomer.Now() - start

	// Report your test result as a failure, if you write it in python, it will looks like this
	// events.request_failure.fire(request_type="udp", name="bar", response_time=100, exception=Exception("udp error"))
	globalBoomer.RecordFailure("udp", "bar", elapsed, "udp error")
}

func waitForQuit() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		globalBoomer.Quit()
		wg.Done()
	}()

	wg.Wait()
}

var globalBoomer = boomer.NewBoomer("127.0.0.1", 5557)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	task1 := &boomer.Task{
		Name:   "foo",
		Weight: 10,
		Fn:     foo,
	}

	task2 := &boomer.Task{
		Name:   "bar",
		Weight: 30,
		Fn:     bar,
	}

	globalBoomer.Run(task1, task2)

	waitForQuit()
	log.Println("shut down")
}
