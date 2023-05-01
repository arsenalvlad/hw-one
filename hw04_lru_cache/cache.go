package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	if _, ok := l.items[key]; !ok {
		if l.queue.Len() > l.capacity {
			l.queue.Remove(l.queue.Back())
			delete(l.items, l.queue.Back().Key)
		}

		l.items[key] = l.queue.PushFront(value)
		l.items[key].Key = key

		return false
	}

	l.queue.Remove(l.items[key])
	l.items[key] = l.queue.PushFront(value)

	return true
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	if item, ok := l.items[key]; ok {
		l.queue.MoveToFront(item)

		return item.Value, true
	}

	return nil, false
}

func (l *lruCache) Clear() {
	l.items = nil
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
