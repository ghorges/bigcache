package queue

import (
	"fmt"
	"strconv"
	"testing"
)

func Printqueuedata(q *queue) {
	fmt.Println("head = " + strconv.Itoa(q.head) +
		"   tail = " + strconv.Itoa(q.tail) +
		"  rightMargin = " + strconv.Itoa(q.rightMargin) +
		"  cap = " + strconv.Itoa(q.cap) +
		"   count = " + strconv.Itoa(q.count),
	)
}
func TestQueuePushAndPop(t *testing.T) {
	queue := NewBytesQueue(100)

	data := "hello"
	index := queue.Push([]byte(data))
	fmt.Println(index)
	Printqueuedata(queue)

	entry, err := queue.Pop()
	if err != nil {
		fmt.Println("err = ", err)
	}
	fmt.Println(string(entry))
	Printqueuedata(queue)
	/*
		entry,err = queue.Pop()
		if err != nil {
			fmt.Println("err = ", err)
		}
		fmt.Println(string(entry))
	*/
}

func TestCapToBig(t *testing.T) {
	queue := NewBytesQueue(1)

	data := "hello111"
	queue.Push([]byte(data))
	queue.Push([]byte(data))
	queue.Push([]byte(data))
	queue.Push([]byte(data))
	queue.Push([]byte(data))
	queue.Push([]byte(data))
	queue.Push([]byte(data))
	queue.Push([]byte(data))

	Printqueuedata(queue)
}

func TestTailBeforeHead(t *testing.T) {
	queue := NewBytesQueue(1)

	data := "hello111"
	queue.Push([]byte(data))
	queue.Push([]byte(data))
	queue.Push([]byte(data))
	queue.Push([]byte(data))
	queue.Push([]byte(data))
	queue.Push([]byte(data))
	queue.Push([]byte(data))
	queue.Push([]byte(data))

	queue.Pop()
	queue.Pop()
	queue.Pop()
	queue.Pop()
	queue.Pop()

	data = "1234567890"
	queue.Push([]byte(data))
	fmt.Println(queue.data)
	Printqueuedata(queue)

	// to test when change alloc,head and tail is ok
	data = "111111111111111111111111111111111111111111111111111111111111" +
		"111111111111111111111111111111111111111111111111111111111111111111"
	queue.Push([]byte(data))
	fmt.Println(queue.data)
	Printqueuedata(queue)
}

func TestQueuePeek(t *testing.T) {
	queue := NewBytesQueue(100)

	data := "hello111"
	queue.Push([]byte(data))

	fmt.Println(queue.Peek())
}

func TestQueueGet(t *testing.T) {
	queue := NewBytesQueue(100)

	data := "hello111"
	queue.Push([]byte(data))

	fmt.Println(queue.Get(1))
}