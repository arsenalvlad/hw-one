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
	Next  *ListItem
	Prev  *ListItem
	Key   Key
}

type list struct {
	length int
	first  *ListItem
	last   *ListItem
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.first
}

func (l *list) Back() *ListItem {
	return l.last
}

func (l *list) PushFront(v interface{}) *ListItem {
	data := &ListItem{}
	data.Value = v

	if l.first != nil {
		data.Next = l.first
		l.first.Prev = data
		if l.last == nil {
			l.last = l.first
			data.Key = l.first.Key
		}
	} else if l.last != nil {
		l.last.Prev = data
		data.Next = l.last
	}

	l.first = data
	l.length++

	return data
}

func (l *list) PushBack(v interface{}) *ListItem {
	data := &ListItem{}
	data.Value = v

	if l.last != nil {
		data.Prev = l.last
		l.last.Next = data
		if l.first == nil {
			l.first = l.last
		}
	} else if l.first != nil {
		l.first.Next = data
		data.Prev = l.first
	}

	l.last = data
	l.length++

	return data
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}
	if i.Prev == nil { //nolint: gocritic
		l.first = i.Next
		i.Next.Prev = nil
	} else if i.Next == nil {
		l.last = i.Prev
		i.Prev.Next = nil
	} else {
		i.Next.Prev = i.Prev
		i.Prev.Next = i.Next
	}
	i = nil
	l.length--
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}

func NewList() List {
	return new(list)
}
