package solver

type cache[T any] struct {
	c       []T
	factory func() T
	reset   func(T) T
}

func (c *cache[T]) get() T {
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

func (c *cache[T]) put(v T) {
	c.c = append(c.c, v)
}
