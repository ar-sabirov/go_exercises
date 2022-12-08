package main

import (
	"fmt"
	"sync"
)

type Foobar struct {
	n        int
	wg       *sync.WaitGroup
	fooToBar chan bool
	barToFoo chan bool
}

func New(n int) *Foobar {
	fooToBar := make(chan bool)
	barToFoo := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(2)

	return &Foobar{
		n:        n,
		wg:       &wg,
		fooToBar: fooToBar,
		barToFoo: barToFoo,
	}
}

func (fb *Foobar) foo() {
	defer fb.wg.Done()

	for i := 0; i < fb.n; i++ {
		<-fb.barToFoo
		fmt.Print("foo")
		fb.fooToBar <- true
	}
	<-fb.barToFoo
	close(fb.fooToBar)
}

func (fb *Foobar) bar() {
	defer fb.wg.Done()

	fb.barToFoo <- true
	for i := 0; i < fb.n; i++ {
		<-fb.fooToBar
		fmt.Print("bar")
		fb.barToFoo <- true
	}
	close(fb.barToFoo)
}

func main() {
	fb := New(5)
	go fb.foo()
	go fb.bar()

	fb.wg.Wait()
}
