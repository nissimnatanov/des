package cache

type Args[T any] struct {
	Factory func() T
	Reset   func(T) T
	MaxSize int
}

func New[T any](args Args[T]) *Cache[T] {
	if args.Factory == nil {
		panic("factory function must be provided")
	}
	if args.MaxSize <= 0 {
		args.MaxSize = 100
	}
	return &Cache[T]{
		c:       make([]T, 0, args.MaxSize),
		factory: args.Factory,
		reset:   args.Reset,
		maxSize: args.MaxSize,
	}
}

type Cache[T any] struct {
	c       []T
	factory func() T
	reset   func(T) T
	maxSize int
}

func (c *Cache[T]) Get() T {
	if len(c.c) == 0 {
		return c.factory()
	}
	v := c.c[len(c.c)-1]
	c.c = c.c[:len(c.c)-1]
	if c.reset != nil {
		v = c.reset(v)
	}
	return v
}

func (c *Cache[T]) Put(v T) {
	if len(c.c) >= c.maxSize {
		return
	}
	c.c = append(c.c, v)
}
