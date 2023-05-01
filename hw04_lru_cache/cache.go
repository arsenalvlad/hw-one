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
		if l.queue.Len() >= l.capacity {
			delete(l.items, l.queue.Back().Key)
			l.queue.Remove(l.queue.Back())
		}

		l.items[key] = l.queue.PushFront(value)
		l.items[key].Key = key

		return false
	}

	l.queue.Remove(l.items[key])
	delete(l.items, key)

	l.items[key] = l.queue.PushFront(value)
	l.items[key].Key = key

	return true
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	if item, ok := l.items[key]; ok {
		l.queue.MoveToFront(item)
		l.items[key].Key = key
		l.queue.Front().Key = key

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
