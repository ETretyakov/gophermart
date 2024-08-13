package closer

import (
	"context"
	"gophermart/internal/log"
)

var globalCloser = New()

func Add(f ...func() error) {
	globalCloser.Add(f...)
}

func CloseAll() {
	globalCloser.CloseAll()
}

type Closer struct {
	funcs []func() error
}

func New() *Closer {
	return &Closer{
		funcs: make([]func() error, 0),
	}
}

func (c *Closer) Add(f ...func() error) {
	c.funcs = append(c.funcs, f...)
}

func (c *Closer) CloseAll() {
	for _, f := range c.funcs {
		if err := f(); err != nil {
			log.Error(context.TODO(), "error close", err)
		}
	}
}
