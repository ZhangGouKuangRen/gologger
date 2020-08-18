package util

import (
	"errors"
	"fmt"
)

type node struct {
	value interface{}
	next *node
}

type Queue struct {
	head *node
	tail *node
	length int
}

func NewQueue()*Queue  {
	queue := Queue{
		nil,
		nil,
		0,
	}
	return &queue
}

func (q *Queue)Size()int  {
	return q.length
}

func (q *Queue)Push(a interface{})  {
	node := &node{value: a, next: nil}
    if q.length == 0 {
    	q.head = node
    	q.tail = node
	} else {
		q.tail.next = node
		q.tail = node
	}
	q.length++
}

func (q *Queue)Poll()(interface{}, error)  {
	if q.length == 0 {
		return 0, errors.New("队列为空")
	} else {
		node := q.head
		q.head = q.head.next
		q.length--
		return node.value, nil
	}
}

func (q *Queue)Peek()(interface{}, error)  {
	if q.length == 0 {
		return 0, errors.New("队列为空")
	} else {
		return q.head.value, nil
	}
}

func (q *Queue)Show()  {
	temp := q.head
	for i:=0; i<q.length; i++ {
		fmt.Print(temp.value, " ")
		temp = temp.next
	}
	fmt.Println()
}