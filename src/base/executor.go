package base

import (
	"time"
)

type Closure func()

type ExecutorStopOpt bool

const (
	StopOptDiscard = false
	StopOptRunAll  = true
)

type Framer interface {
	Frame()
}

type Executor struct {
	tasks chan Closure
	stop  chan ExecutorStopOpt
}

func NewExecutor(taskSize int) *Executor {
	e := &Executor{
		tasks: make(chan Closure, taskSize),
		stop:  make(chan ExecutorStopOpt),
	}
	return e
}

func (e *Executor) Post(c Closure) {
	e.tasks <- c
}

func (e *Executor) Dispatch(c Closure) {
	signal := make(chan bool)
	e.tasks <- func() {
		c()
		close(signal)
	}
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
			if opt != StopOptRunAll {
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
			if opt != StopOptRunAll {
				return
			}

			for task := range e.tasks {
				task()
			}
			return
		}
	}
}

func (e *Executor) Stop(opt ExecutorStopOpt) {
	close(e.tasks)
	e.stop <- opt
}

func (e *Executor) WaitForStop() {
	<-e.stop
}
