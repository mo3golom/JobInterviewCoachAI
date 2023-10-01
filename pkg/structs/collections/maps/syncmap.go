package maps

import "sync"

type SyncMap[T1 comparable, T2 any] struct {
	m map[T1]T2
	sync.RWMutex
}

func NewSyncMap[T1 comparable, T2 any]() *SyncMap[T1, T2] {
	return &SyncMap[T1, T2]{
		m: make(map[T1]T2, 0),
	}
}

func NewSyncMapWithSize[T1 comparable, T2 any](size uint64) *SyncMap[T1, T2] {
	return &SyncMap[T1, T2]{
		m: make(map[T1]T2, size),
	}
}

func (s *SyncMap[T1, T2]) Set(key T1, value T2) {
	s.Lock()
	s.m[key] = value
	s.Unlock()
}

func (s *SyncMap[T1, T2]) Get(key T1) (T2, bool) {
	s.RLock()
	value, ok := s.m[key]
	s.RUnlock()
	return value, ok
}

func (s *SyncMap[T1, T2]) Delete(key T1) {
	s.Lock()
	delete(s.m, key)
	s.Unlock()
}

func (s *SyncMap[T1, T2]) Copy() map[T1]T2 {
	s.RLock()
	m := make(map[T1]T2, len(s.m))
	for k, v := range s.m {
		m[k] = v
	}
	s.RUnlock()
	return m
}

func (s *SyncMap[T1, T2]) Len() int {
	s.RLock()
	l := len(s.m)
	s.RUnlock()
	return l
}

func (s *SyncMap[T1, T2]) Iterate(f func(key T1, value T2)) {
	s.RLock()
	for k, v := range s.m {
		f(k, v)
	}
	s.RUnlock()
}

// GetMap - ВНИМАНИЕ!! не потокобезопасный метод, использовать можно только в тех случаях когда уверены,
// что параллельные процессы не пытаются изменить значения в s.m
func (s *SyncMap[T1, T2]) GetMap() map[T1]T2 {
	return s.m
}
