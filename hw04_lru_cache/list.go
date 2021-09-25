package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Prev  *ListItem
	Next  *ListItem
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	i := &ListItem{
		Value: v,
		Next:  l.front,
		Prev:  nil,
	}

	if l.front != nil {
		l.front.Prev = i
	}

	l.front = i
	if l.back == nil {
		l.back = i
	}

	l.len++

	return l.front
}

func (l *list) PushBack(v interface{}) *ListItem {
	i := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  l.back,
	}

	if l.back != nil {
		l.back.Next = i
	}

	l.back = i
	if l.front == nil {
		l.front = i
	}

	l.len++

	return l.back
}

func (l *list) Remove(i *ListItem) {
	if l.len == 0 {
		return
	}

	if i == l.front {
		l.front = i.Next
		l.front.Prev = nil
	} else if i == l.back {
		l.back = i.Prev
		l.back.Next = nil
	} else {
		i.Next.Prev = i.Prev
		i.Prev.Next = i.Next
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.front != nil {
		l.Remove(i)
		l.PushFront(i.Value)
	}
}

func NewList() List {
	return new(list)
}
