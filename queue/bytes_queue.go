package queue

import (
	"log"
	"encoding/binary"
	"time"
)

const (
	// to storage entry lens
	headEntrySize = 4
	// not zero.if is zero,to return errInvalidIndex
	leftMargin = 1
	// between tail and head,when tail < head,use this
	emptyBlob = 20
)

var (
	errEmptyQueue        = &queueError{error: "bigcache queue is empty"}
	errInvalidIndex      = &queueError{error: "invalid index, index must big than zero"}
	errIndexMoreThanSize = &queueError{error: "index more than size"}
)

type Queue struct {
	data         []byte
	head         int
	tail         int
	rightMargin  int
	cap          int
	count        int
	headByte     []byte
	initDataSize int
}

type queueError struct {
	error string
}

func NewBytesQueue(initDataSize int) *Queue {
	return &Queue{
		data:         make([]byte, initDataSize),
		head:         leftMargin,
		tail:         leftMargin,
		rightMargin:  leftMargin,
		cap:          initDataSize,
		count:        0,
		headByte:     make([]byte, headEntrySize),
		initDataSize: initDataSize,
	}
}

func (q *Queue) Push(wrap []byte) int {
	dataLen := len(wrap)

	if q.afterTailSpace() < dataLen+headEntrySize {
		if q.beforeHeadSpace() >= dataLen+headEntrySize {
			q.tail = leftMargin
		} else {
			q.allocMemory(dataLen + headEntrySize)
		}
	}

	index := q.tail

	q.push(wrap, len(wrap))

	return index
}

func (q *Queue) Pop() ([]byte, error) {
	data, len, err := q.peek(q.head)

	if err != nil {
		return data, err
	}

	q.head += len + headEntrySize
	q.count--

	if q.head == q.rightMargin {
		q.head = leftMargin
		if q.tail == q.rightMargin {
			q.tail = leftMargin
		}
		q.rightMargin = q.tail
	}

	return data, nil
}

func (q *Queue) checkOut(index int) error {

	if q.count == 0 {
		return errEmptyQueue
	}

	if index < leftMargin {
		return errInvalidIndex
	}

	if index > q.rightMargin {
		return errIndexMoreThanSize
	}

	return nil
}

// return head entry
func (q *Queue) Peek() ([]byte, error) {
	data, _, err := q.peek(q.head)
	return data, err
}

func (q *Queue) Get(index int) ([]byte, error) {
	data, _, err := q.peek(index)
	return data, err
}

func (q *Queue) peek(index int) ([]byte, int, error) {
	err := q.checkOut(index)
	if err != nil {
		return nil, 0, err
	}

	len := binary.LittleEndian.Uint32(q.data[index : index+headEntrySize])
	return q.data[index+headEntrySize : index+headEntrySize+int(len)], int(len), nil
}

func (q *Queue) allocMemory(len int) {
	start := time.Now()

	if q.cap < len {
		q.cap += len
	}
	q.cap = 2 * q.cap
	oldData := q.data
	q.data = make([]byte, q.cap)

	if q.rightMargin != leftMargin {
		copy(q.data, oldData[:q.rightMargin])
		if q.tail < q.head {
			paddingSize := q.head - q.tail - headEntrySize
			q.push(make([]byte, paddingSize), paddingSize)
			q.head = leftMargin
			q.tail = q.rightMargin
		}
	}

	log.Println("alloc in ", time.Since(start), " cap is ", q.cap)
}

func (q *Queue) push(data []byte, len int) {
	binary.LittleEndian.PutUint32(q.headByte, uint32(len))

	copy(q.data[q.tail:], q.headByte[:headEntrySize])
	copy(q.data[q.tail+headEntrySize:], data[:len])

	q.tail += headEntrySize + len

	if q.tail > q.head {
		q.rightMargin = q.tail
	}

	q.count++
}

func (q *Queue) afterTailSpace() int {
	if q.tail >= q.head {
		return q.cap - q.tail
	}
	return q.head - q.tail - emptyBlob
}

func (q *Queue) beforeHeadSpace() int {
	if q.tail >= q.head {
		return q.head - leftMargin - emptyBlob
	}
	return q.head - q.tail - emptyBlob
}

// write error interface
func (err *queueError) Error() string {
	return err.error
}
