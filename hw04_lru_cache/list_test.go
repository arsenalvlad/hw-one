package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("front", func(t *testing.T) {
		l := NewList()

		l.PushFront(30) // [30]
		l.PushFront(20) // [20, 30]
		l.PushFront(10) // [10, 20, 30]

		require.Equal(t, 3, l.Len())
		require.Equal(t, 30, l.Back().Value)
		require.Equal(t, 10, l.Front().Value)

		l.Remove(l.Back())
		require.Equal(t, 20, l.Back().Value)
		require.Equal(t, 10, l.Front().Value)
	})

	t.Run("back", func(t *testing.T) {
		l := NewList()

		l.PushBack(30) // [30]
		l.PushFront(5) // [5, 30]
		require.Equal(t, 30, l.Back().Value)
		require.Equal(t, 5, l.Front().Value)
		l.PushBack(20) // [5, 30, 20]
		l.PushBack(10) // [5, 30, 20, 10]
		l.PushFront(1) // [1, 5, 30, 20, 10]

		require.Equal(t, 5, l.Len())
		require.Equal(t, 10, l.Back().Value)
		require.Equal(t, 1, l.Front().Value)

		l.Remove(l.Back())
		require.Equal(t, 20, l.Back().Value)
		require.Equal(t, 1, l.Front().Value)
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		l.PushFront(3)  // [3, 10, 20, 30]
		require.Equal(t, 4, l.Len())

		two := l.Front().Next // 10
		l.Remove(two)         // [3, 20, 30]
		require.Equal(t, 3, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 3, 20, 30, 50, 70]

		require.Equal(t, 8, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 3, 20, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 3, 20, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 3, 20, 30, 50}, elems)
	})
}
