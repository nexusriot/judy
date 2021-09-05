package taskqueue

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTaskQueue_Empty(t *testing.T) {
	q := NewTaskQueue()
	assert.Equal(t, true, q.Empty())
}

func TestTaskQueue_Enqueue(t *testing.T) {
	q := NewTaskQueue()
	q.Enqueue(Task{Id: "1"})
	assert.Equal(t, 1, q.Size())
}

func TestTaskQueue_Dequeue(t *testing.T) {
	q := NewTaskQueue()
	q.Enqueue(Task{Id: "1"})
	task, err := q.Dequeue()
	assert.Equal(t, true, q.Empty())
	require.Nil(t, err)
	assert.Equal(t, "1", task.Id)
}
