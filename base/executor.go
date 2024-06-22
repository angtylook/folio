package base

import (
	"sync"
	"time"
)

type Closure func()

const (
	kStopOptDiscard = false
	kStopOptRunAll  = true
)

type Framer interface {
	Frame()
}

type Executor struct {
	tasks  chan Closure
	stop   chan bool
	closed int32
	l      sync.Mutex
}

func NewExecutor(taskSize int) *Executor {
	e := &Executor{
		tasks:  make(chan Closure, taskSize),
		stop:   make(chan bool),
		closed: 0,
	}
	return e
}

func (e *Executor) Post(c Closure) {
	e.l.Lock()
	defer e.l.Unlock()
	if e.closed == 1 {
		return
	}

	e.tasks <- c
}

func (e *Executor) Dispatch(c Closure) {
	e.l.Lock()
	if e.closed == 1 {
		e.l.Unlock()
		return
	}
	signal := make(chan bool)
	e.tasks <- func() {
		c()
		close(signal)
	}
	e.l.Unlock()
	<-signal
}

func (e *Executor) Start() {
	for {
		select {
		case task, ok := <-e.tasks:
			if ok && task != nil {
				task()
			}
		case opt := <-e.stop:
			if opt != kStopOptRunAll {
				close(e.stop)
				return
			}

			for task := range e.tasks {
				task()
			}
			close(e.stop)
			return
		}
	}

}

func (e *Executor) StartWithFrame(framer Framer, d time.Duration) {
	ticker := time.NewTicker(d)
	defer func() {
		ticker.Stop()
		close(e.stop)
	}()

	for {
		select {
		case <-ticker.C:
			framer.Frame()
		case task, ok := <-e.tasks:
			if ok && task != nil {
				task()
			}
		case opt := <-e.stop:
			if opt != kStopOptRunAll {
				return
			}

			for task := range e.tasks {
				task()
			}
			return
		}
	}
}

func (e *Executor) Stop() {
	e.l.Lock()
	if e.closed == 1 {
		e.l.Unlock()
		return
	}

	e.closed = 1
	close(e.tasks)
	e.l.Unlock()
	e.stop <- kStopOptRunAll
}

func (e *Executor) StopNow() {
	e.l.Lock()
	if e.closed == 1 {
		e.l.Unlock()
		return
	}

	e.closed = 1
	close(e.tasks)
	e.l.Unlock()
	e.stop <- kStopOptDiscard
}

func (e *Executor) Wait() {
	<-e.stop
}
