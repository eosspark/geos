package database

import (
	"container/list"
	"sync"
)

/////////////////////////////////////////////////////// Deque //////////////////////////////////////////////////////////
type deque struct {
	sync.RWMutex
	container *list.List
	capacity  int
}

func newDeque() *deque {
	return newCappedDeque(-1)
}

func newCappedDeque(capacity int) *deque {
	return &deque{
		container: list.New(),
		capacity:  capacity,
	}
}

func (s *deque) Append(item interface{}) bool {
	s.Lock()
	defer s.Unlock()

	if s.capacity < 0 || s.container.Len() < s.capacity {
		s.container.PushBack(item)
		return true
	}
	return false
}

func (s *deque) PopFront() interface{} {
	s.Lock()
	defer s.Unlock()

	var item interface{} = nil
	var firstContainerItem *list.Element = nil

	firstContainerItem = s.container.Front()
	if firstContainerItem != nil {
		item = s.container.Remove(firstContainerItem)
	}

	return item
}

func (s *deque) LastSecond() interface{} {
	last := s.Pop()
	second := s.Pop()
	s.Append(second)
	s.Append(last)
	return second
}

func (s *deque) Pop() interface{} {
	s.Lock()
	defer s.Unlock()

	var item interface{} = nil
	var lastContainerItem *list.Element = nil

	lastContainerItem = s.container.Back()
	if lastContainerItem != nil {
		item = s.container.Remove(lastContainerItem)
	}

	return item
}

func (s *deque) Size() int {
	s.RLock()
	defer s.RUnlock()

	return s.container.Len()
}

func (s *deque) First() interface{} {
	s.RLock()
	defer s.RUnlock()

	item := s.container.Front()
	if item != nil {
		return item.Value
	} else {
		return nil
	}
}

func (s *deque) Last() interface{} {
	s.RLock()
	defer s.RUnlock()

	item := s.container.Back()
	if item != nil {
		return item.Value
	} else {
		return nil
	}
}

func (s *deque) Empty() bool {
	s.RLock()
	defer s.RUnlock()

	return s.container.Len() == 0
}
