package main

import (
	"log"
	"strconv"
	"time"
)

func main() {
	a_chan := make(chan string, 2)
	b_chan := make(chan string, 2)
	c_chan := make(chan string, 2)
	msg_chan := make(chan string, 1)
	drop_chan := make(chan bool, 1)
	count := 1
	go pass_b(a_chan, b_chan, msg_chan)
	go pass_c(b_chan, c_chan, msg_chan)
	go c_drop(c_chan, drop_chan)
	go msg_sender(msg_chan, drop_chan)

	for range time.Tick(time.Second * 1) {
		log.Printf("before ,a has %d balls, b has %d balls, c has %d balls.", len(a_chan), len(b_chan), len(c_chan))

		if count == 15 {
			for c := range c_chan {
				log.Println("c", c)
			}
		}
		log.Println("-------------")
		count_s := strconv.Itoa(count)
		// 發球
		a_chan <- "ball" + count_s
		log.Printf("發球%d", count)
		count++
	}

}

func pass_b(a_chan chan string, b_chan chan string, msg_chan chan string) {
	for {
		b_ball := <-a_chan
	SEND:
		for {
			select {
			case b_chan <- b_ball:
				break SEND
			case <-time.After(time.Second * 1):
				msg_chan <- "b快點接球"
			}
		}
	}
}

func pass_c(b_chan chan string, c_chan chan string, msg_chan chan string) {
	for {
		c_ball := <-b_chan
	SEND:
		for {
			select {
			case c_chan <- c_ball:
				break SEND
			case <-time.After(time.Second * 1):
				msg_chan <- "c快點接球"
			}
		}

	}
}

func msg_sender(msg_chan chan string, drop_chan chan bool) {
	for {
		msg := <-msg_chan
		log.Println(msg)
		if msg == "c快點接球" {
			drop_chan <- true
		}
	}
}

func c_drop(c_chan chan string, drop_chan chan bool) {
	for {
		drop := <-drop_chan
		if drop {
			<-c_chan
			log.Println("丟球")
			time.Sleep(time.Second * 2)
		}
	}
}
