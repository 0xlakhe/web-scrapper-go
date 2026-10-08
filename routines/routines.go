package main

import (
	"fmt"
	"sync"
	"time"
)

func r(s string) {
	for i := 0; i < 3; i++ {
		fmt.Printf("%s, %d\n", s, i)
	}
	time.Sleep(1 * time.Second)

}

func aa(a chan<- string, b string) {
	a <- b
}

func bb(a <-chan string, b chan<- string) {

	b <- <-a
}

func main() {

	var wg sync.WaitGroup
	a := make(chan string, 1)
	b := make(chan string, 1)
	aa(a, "hello from go")
	bb(a, b)
	fmt.Println(<-b)
	messages := make(chan string)
	msgs := make(chan int)

	wg.Add(1)
	go func() {
		for range 5 {
			msg := <-msgs
			fmt.Println(msg)
		}
		defer wg.Done()
	}()

	wg.Add(1)
	go func() {
		for i := range 5 {
			msgs <- i
		}
		defer wg.Done()
	}()

	wg.Add(1)
	go func() {
		s1 := <-messages
		r(s1)
		defer wg.Done()
	}()
	messages <- "hello"

	wg.Add(1)

	go func(s2 string) {
		r(s2)
		defer wg.Done()
	}("bye")

	wg.Wait()
	r("Hello")

	fmt.Println("done")
}
