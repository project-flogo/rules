package main

import (
	"fmt"
	"time"
)

func consume(c chan int) {

	for {
		fmt.Printf("trying consume..\n")

		x := <-c
		fmt.Printf("Consumed [%d]\n", x)

	}

}

func produce(c chan int) {

	for i := 0; i < 20; i++ {
		c <- i
		fmt.Printf("Produced [%d]\n", i)
		time.Sleep(1 * time.Second)
	}
	time.Sleep(2 * time.Second)
	close(c)
}

func main() {

	c := make(chan int)
	go consume(c)
	time.Sleep(1 * time.Second)
	fmt.Printf("Exiting..\n")

}
