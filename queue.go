package main

import (
	"fmt"
	"log"
)

type QueueUsers struct {
	Queue []*Client
}

func (q *QueueUsers) AddtoQueue(c *Client) (*Client, bool) {
	mu.Lock()
	defer mu.Unlock()
	if len(q.Queue) == 0 {
		q.Queue = append(q.Queue, c)
		return nil, false
	} else {
		log.Println("айді до якого підєднуюсь:", q.Queue[0].ID)
		return q.Queue[0], true
	}
}

func (q *QueueUsers) DeleteFromQueue() {
	mu.Lock()
	q.Queue = q.Queue[:0]
	fmt.Println(q.Queue)
	mu.Unlock()
}
