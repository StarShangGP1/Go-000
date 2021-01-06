package main

import (
	"container/ring"
	"fmt"
	"math/rand"
	"time"
)

type slidingWindow struct {
	maxRequest   int
	unitTime     time.Duration
	windowTime   time.Duration
	windowsCount int
}

type window struct {
	index              int
	maxRequestCount    int
	handleRequestCount int
}

var requestChan chan string
var request []string

func init() {
	rand.Seed(42)
	requestChan = make(chan string, 0)
	request = make([]string, 0)
}

func NewSlidingWindow(maxRequest, windowsCount int, unitTime time.Duration) *slidingWindow {
	return &slidingWindow{
		maxRequest:   maxRequest,
		unitTime:     unitTime,
		windowTime:   unitTime / time.Duration(windowsCount),
		windowsCount: windowsCount,
	}
}

func getRequest() {
	for {
		time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
		requestChan <- "url"
	}
}

func getRequestCountSum(r *ring.Ring) int {
	var sum int
	for i := 0; i < r.Len(); i++ {
		w := r.Value.(*window)
		sum += w.handleRequestCount
		r = r.Next()
	}
	return sum
}

func main() {
	swl := NewSlidingWindow(100, 10, 10*time.Second)

	r := ring.New(swl.windowsCount)
	for i := 0; i < r.Len(); i++ {
		r.Value = &window{i, swl.maxRequest, 0}
		r = r.Next()
	}

	ticker := time.NewTicker(swl.windowTime)
	defer ticker.Stop()

	go getRequest()

	w := r.Value.(*window)
	for {
		select {
		case url := <-requestChan:
			if w.maxRequestCount > 0 && w.handleRequestCount < w.maxRequestCount {
				request = append(request, url)
				w.handleRequestCount++
				fmt.Printf("this %d window, can handle request max:%d, current has handled request:%d\n", w.index, w.maxRequestCount, w.handleRequestCount)
				continue
			}
		case <-ticker.C:
			fmt.Println("one windowtime has passed ")
			r = r.Move(1)
			w = r.Value.(*window)
			w.maxRequestCount = swl.maxRequest - (getRequestCountSum(r) - w.handleRequestCount)
			w.handleRequestCount = 0
		}
	}
}
