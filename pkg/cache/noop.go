package cache

// Noop is a cache that does nothing.
type Noop[T any] struct{}

func (n *Noop[T]) Get(key string) (T, bool) { return *new(T), false }
func (n *Noop[T]) Set(key string, _ T)      {}
func (n *Noop[T]) Len() int                 { return 0 }
func (n *Noop[T]) ToMap() map[string]T      { return nil }
