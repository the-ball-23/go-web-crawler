package main

import "sync"

type SafeUrlQueue struct {
	urls map[string]int
	mux  sync.Mutex
}

func (queue *SafeUrlQueue) Push(url string, depth int) {
	queue.mux.Lock()
	queue.urls[url] = depth
	queue.mux.Unlock()
}

func (queue *SafeUrlQueue) Pop() (string, int) {
	queue.mux.Lock()
	defer queue.mux.Unlock()
	var url string
	var depth int
	for k, v := range queue.urls {
		url = k
		depth = v
		break
	}
	delete(queue.urls, url)
	return url, depth
}

func (queue *SafeUrlQueue) Len() int {
	queue.mux.Lock()
	defer queue.mux.Unlock()
	length := len(queue.urls)
	return length
}
	