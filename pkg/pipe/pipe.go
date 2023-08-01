package pipe

type Pipe[T any] struct {
	c chan T
}

func (p *Pipe[T]) init() {
	p.c = make(chan T)
}
