package base

type FutureFunc func() (interface{}, error)
type FutureFuncSuccess func(arg interface{})
type FutureFuncFail func(err error)
type FutureFuncReply func(arg interface{}, err error)

type Future struct {
	successFunc func(arg interface{})
	successExec *Executor
	failFunc    func(err error)
	failExec    *Executor
	replyFunc   func(arg interface{}, err error)
	replyExec   *Executor

	done chan bool
}

func NewFuture() *Future {
	f := &Future{
		done: make(chan bool, 1),
	}
	return f
}

func (f *Future) Success(cb FutureFuncSuccess) *Future {
	f.successFunc = cb
	return f
}

func (f *Future) SuccessAt(executor *Executor, cb FutureFuncSuccess) *Future {
	f.successExec = executor
	f.successFunc = cb
	return f
}

func (f *Future) Fail(cb FutureFuncFail) *Future {
	f.failFunc = cb
	return f
}

func (f *Future) FailAt(executor *Executor, cb FutureFuncFail) *Future {
	f.failExec = executor
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

func (f *Future) Execute(cb FutureFunc) {
	go f.execute(cb)
}

func (f *Future) ExecuteAt(executor *Executor, cb FutureFunc) {
	executor.Post(func() {
		f.execute(cb)
	})
}

func (f *Future) execute(cb FutureFunc) {
	arg, err := cb()
	if f.successFunc != nil {
		f.onSuccess(arg, err)
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

func (f *Future) onSuccess(arg interface{}, err error) {
	if err != nil {
		return
	}

	if f.successExec == nil {
		go f.successFunc(arg)
		return
	}

	f.successExec.Post(func() {
		f.successFunc(arg)
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
		f.replyFunc(arg, err)
	})
}

func (f *Future) Wait() bool {
	_, ok := <-f.done
	return ok
}
