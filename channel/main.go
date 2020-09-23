package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string, 2)

	// go func() {
	// 	ch <- "hello"
	// 	ch <- "golang"
	// }()

	ch <- "hello"

	ch <- "golang"
	fmt.Println(<-ch)
	fmt.Println(<-ch)
	fmt.Println("------------------")
	done := make(chan bool, 1)
	go worker(done)

	<-done
	fmt.Println("------------------")
	pings := make(chan string, 1)
	pongs := make(chan string, 1)
	ping(pings, "passed message")
	pong(pings, pongs)
	fmt.Println(<-pongs)
	fmt.Println("------------------")

	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		time.Sleep(1 * time.Second)
		c1 <- "One"
	}()

	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "Two"
	}()

	for i := 0; i < 2; i++ {
		select {
		case d := <-c1:
			fmt.Printf("received %s\n", d)

		case <-c2:
			fmt.Println("received ")

		}
	}

	fmt.Println("----------")
	ch1 := make(chan string)

	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- "test 1"
	}()

	select {
	case m := <-ch1:
		fmt.Println(m)
	case <-time.After(2 * time.Second):
		fmt.Println("TimeOut 1")
	}

	fmt.Println("----------")

	messages := make(chan string)
	//signal := make(chan bool)

	msg := "hi"
	select {
	case messages <- msg:
		fmt.Println("received message", msg)
	default:
		fmt.Println("no message received")
	}

	fmt.Println("----------")

	jobs := make(chan int, 5)
	doed := make(chan bool)

	go func() {
		for {
			j, more := <-jobs
			if more {
				fmt.Println("received job", j)
			} else {
				fmt.Println("received all jobs")
				doed <- true
				return
			}
		}
	}()

	for j := 1; j <= 3; j++ {
		jobs <- j
		fmt.Println("sent job", j)
	}
	close(jobs)
	fmt.Println("sent all jobs")

	<-doed
}

func worker(done chan bool) {
	fmt.Print("working ...")
	time.Sleep(time.Second)
	fmt.Println("done")
	done <- true
}

func ping(pings chan<- string, msg string) {
	pings <- msg
}

func pong(pings <-chan string, pongs chan<- string) {
	msg := <-pings
	pongs <- msg
}
