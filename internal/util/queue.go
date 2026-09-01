package util

type Queue[T any] struct {
	items []T
	head  int
}

func (q *Queue[T]) Len() int {
	return len(q.items) - q.head
}

func (q *Queue[T]) Peek() (*T, bool) {
	if q.Len() == 0 {
		return nil, false
	}

	return &q.items[q.head], true
}

func (q *Queue[T]) Set(item T) {
	q.Clear()
	q.items = append(q.items, item)
}

func (q *Queue[T]) Push(item T) bool {
	was_empty := q.Len() == 0
	if was_empty && q.head > 0 {
		q.items = q.items[:0]
		q.head = 0
	}

	q.items = append(q.items, item)
	return was_empty
}

func (q *Queue[T]) Pop() (*T, bool) {
	if q.Len() == 0 {
		return nil, false
	}

	var zero T
	q.items[q.head] = zero
	q.head++

	if q.Len() == 0 {
		q.items = q.items[:0]
		q.head = 0
		return nil, false
	}
	if q.head > 32 && q.head*2 >= len(q.items) {
		q.compact()
	}

	return q.Peek()
}

func (q *Queue[T]) Clear() {
	clear(q.items[q.head:])
	q.items = q.items[:0]
	q.head = 0
}

func (q *Queue[T]) compact() {
	next_len := copy(q.items, q.items[q.head:])
	clear(q.items[next_len:])
	q.items = q.items[:next_len]
	q.head = 0
}
