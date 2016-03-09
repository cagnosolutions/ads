package xdb

import "sync"

// single queue entry
type qnode struct {
	data interface{}
	next *qnode
}

// struct representing a fifo queue
// all methods of the queue are safe
// for multiple go-routines / threads
type queue struct {
	head, tail *qnode
	size       int
	sync.RWMutex
}

// return a new fifo queue instance
func NewQueue() *queue {
	return &queue{}
}

// returns number of entries in queue
func (q *queue) Size() int {
	q.RLock()
	n := q.size
	q.RUnlock()
	return n
}

// pushes an entry at the tail of the queue
func (q *queue) Push(v interface{}) {
	q.Lock()
	n := &qnode{data: v}
	if q.tail == nil {
		q.tail, q.head = n, n
	} else {
		q.tail.next, q.tail = n, n
	}
	q.size++
	q.Unlock()
}

//  returns entry at the front, ie. the oldest
func (q *queue) Poll() interface{} {
	q.Lock()
	if q.head == nil {
		return nil
	}
	n := q.head
	q.head = n.next
	if q.head == nil {
		q.tail = nil
	}
	v := n.data
	q.size--
	q.Unlock()
	return v
}

//  returns entry at front of queue, does NOT mutate
func (q *queue) Peek() interface{} {
	q.RLock()
	if q.head == nil {
		return nil
	}
	v := q.head.data
	q.RUnlock()
	return v
}
