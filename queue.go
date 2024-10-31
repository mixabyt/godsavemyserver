package main

type QueueUsers struct {
	Queue []*Client
}

func (q *QueueUsers) AddtoQueue(c *Client) bool {
	if len(q.Queue) == 0 {
		q.Queue = append(q.Queue, c)
		return false
	} else {
		return true
	}
}
