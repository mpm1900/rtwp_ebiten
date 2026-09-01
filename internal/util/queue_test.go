package util

import "testing"

func TestQueuePushPop(t *testing.T) {
	var queue Queue[int]

	if !queue.Push(1) {
		t.Fatal("first push should report an empty queue")
	}
	if queue.Push(2) {
		t.Fatal("second push should not report an empty queue")
	}
	if queue.Len() != 2 {
		t.Fatalf("expected len 2, got %d", queue.Len())
	}

	next, ok := queue.Pop()
	if !ok || *next != 2 {
		t.Fatalf("expected next item 2, got %v, %t", next, ok)
	}
	if queue.Len() != 1 {
		t.Fatalf("expected len 1, got %d", queue.Len())
	}

	next, ok = queue.Pop()
	if ok || next != nil {
		t.Fatalf("expected empty queue, got %v, %t", next, ok)
	}
	if !queue.Push(3) {
		t.Fatal("push after draining should report an empty queue")
	}
}

func TestQueueSetReplacesItems(t *testing.T) {
	var queue Queue[int]

	queue.Push(1)
	queue.Push(2)
	queue.Set(3)

	item, ok := queue.Peek()
	if !ok || *item != 3 {
		t.Fatalf("expected item 3, got %v, %t", item, ok)
	}
	if queue.Len() != 1 {
		t.Fatalf("expected len 1, got %d", queue.Len())
	}
}

func TestQueuePeekAllowsMutation(t *testing.T) {
	type queueItem struct {
		started bool
	}

	var queue Queue[queueItem]
	queue.Push(queueItem{})

	item, ok := queue.Peek()
	if !ok {
		t.Fatal("expected queued item")
	}
	item.started = true

	item, ok = queue.Peek()
	if !ok || !item.started {
		t.Fatal("expected peeked item mutation to persist")
	}
}
