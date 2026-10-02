package resources

import (
	"io"
	"sync"
)

type Shared[T any] struct {
	mu       sync.Mutex
	refcount int
	resource T
}

func (s *Shared[T]) Acquire(factory func() (T, error)) (*Lease[T], error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.refcount == 0 {
		res, err := factory()
		if err != nil {
			return nil, err
		}
		s.resource = res
	}

	s.refcount++

	return &Lease[T]{
		Resource: s.resource,
		release:  s.onceReleaser(),
	}, nil
}

func (s *Shared[T]) release() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.refcount == 0 {
		panic("resources.Shared: no resource acquired")
	}

	s.refcount--
	if s.refcount > 0 {
		return nil
	}

	var zero T
	v := any(s.resource)
	s.resource = zero

	if c, ok := v.(io.Closer); ok {
		return c.Close()
	}

	return nil
}

func (s *Shared[T]) onceReleaser() func() error {
	var once sync.Once
	var err error

	return func() error {
		once.Do(func() {
			err = s.release()
		})
		return err
	}
}

type Lease[T any] struct {
	Resource T
	release  func() error
}

func (l *Lease[T]) Close() error {
	return l.release()
}
