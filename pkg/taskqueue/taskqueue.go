package taskqueue

import (
	"container/list"
	"fmt"
	"time"
)

type Task struct {
	Id      string
	Payload []byte
	Taken   time.Time
}

type TaskQueue struct {
	queue *list.List
}

func NewTaskQueue() TaskQueue {
	return TaskQueue{queue: list.New()}
}

func (t *TaskQueue) Enqueue(value Task) {
	t.queue.PushBack(value)
}

func (t *TaskQueue) Dequeue() (*Task, error) {
	if t.Size() > 0 {
		elem := t.queue.Front()
		task := elem.Value.(Task)
		t.queue.Remove(elem)
		return &task, nil
	}
	return nil, fmt.Errorf("queue is empty")
}

func (t *TaskQueue) Size() int {
	return t.queue.Len()
}

func (t *TaskQueue) Empty() bool {
	return t.Size() == 0
}
