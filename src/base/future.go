package base

type FutureFunc func() (interface{}, error)
type FutureFuncSucc func(arg interface{})
type FutureFuncFail func(err error)
type FutureFuncReply func(arg interface{}, err error)

type Future struct {
	succFunc  func(arg interface{})
	succExec  *Executor
	failFunc  func(err error)
	failExec  *Executor
	replyFunc func(arg interface{}, err error)
	replyExec *Executor

	done chan bool
}

func NewFuture() *Future {
	f := &Future{
		done: make(chan bool, 1),
	}
	return f
}

func (f *Future) Succ(cb FutureFuncSucc) *Future {
	f.succFunc = cb
	return f
}

func (f *Future) SuccAt(exector *Executor, cb FutureFuncSucc) *Future {
	f.succExec = exector
	f.succFunc = cb
	return f
}

func (f *Future) Fail(cb FutureFuncFail) *Future {
	f.failFunc = cb
	return f
}

func (f *Future) FailAt(exector *Executor, cb FutureFuncFail) *Future {
	f.failExec = exector
	f.failFunc = cb
	return f
}

func (f *Future) Reply(cb FutureFuncReply) *Future {
	f.replyFunc = cb
	return f
}

func (f *Future) ReplyAt(executor *Executor, cb FutureFuncReply) *Future {
	f.replyExec = executor
	f.replyFunc = cb
	return f
}

func (f *Future) Execute(cb *FutureFunc) {
	go f.execute(cb)
}

func (f *Future) ExecuteAt(executor *Executor, cb *FutureFunc) {
	executor.Post(func() {
		f.execute(cb)
	})
}

func (f *Future) execute(cb *FutureFunc) {
	arg, err := cb()
	if f.succFunc != nil {
		f.onSucc(arg, err)
	}
	if f.failFunc != nil {
		f.onFail(arg, err)
	}
	if f.replyFunc != nil {
		f.onReply(arg, err)
	}
	f.done <- true
	close(f.done)
}

func (f *Future) onSucc(arg interface{}, err error) {
	if err != nil {
		return
	}

	if f.succExec == nil {
		go f.succFunc(arg)
		return
	}

	f.succExec.Post(func() {
		f.succFunc(arg)
	})
}

func (f *Future) onFail(arg interface{}, err error) {
	if err == nil {
		return
	}

	if f.failExec == nil {
		go f.failFunc(err)
		return
	}

	f.failExec.Post(func() {
		f.failFunc(err)
	})
}

func (f *Future) onReply(arg interface{}, err error) {
	if f.replyExec == nil {
		go f.replyFunc(arg, err)
		return
	}

	f.replyExec.Post(func() {
		f.replyExec(arg, err)
	})
}

func (f *Future) Wait() {
	_, _ := <-f.done
}
