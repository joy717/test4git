package poolasync

import (
	"fmt"
	"runtime/debug"
	"sync/atomic"
)

const DefaultLimit = 10

// PoolAsync has a pool for async.
type PoolAsync struct {
	// pool properties
	poolChan chan int
	poolSize int
	// job properties
	jobChan chan int
	jobSize int32
	// job error list. the cursor is the index of jobs.
	// if job doesn't have error, then the element is nil
	errList []error
}

func NewPoolAsync(count int) *PoolAsync {
	if count <= 0 {
		count = DefaultLimit
	}
	result := &PoolAsync{
		poolChan: make(chan int, count),
		jobChan:  make(chan int),
		poolSize: count,
		errList:  make([]error, 0),
	}
	// init pool chan.
	for i := 0; i < count; i++ {
		result.poolChan <- i
	}
	return result
}

func NewDefaultPoolAsync() *PoolAsync {
	return NewPoolAsync(DefaultLimit)
}

// Deprecated
//
// get a ticket from poolChan. if got, then do f(), if not, wait for poolChan.
func (a *PoolAsync) Do(f func()) *PoolAsync {
	ticket := <-a.poolChan
	atomic.AddInt32(&a.jobSize, 1)
	go func() {
		defer func() {
			a.poolChan <- ticket            //release pool ticket
			a.jobChan <- ticket             // done job.
			atomic.AddInt32(&a.jobSize, -1) // decrease job size
		}()
		f()
	}()
	return a
}

// get a ticket from poolChan. if got, then do f(), if not, wait for poolChan.
func (a *PoolAsync) DoWitError(f func() error) *PoolAsync {
	ticket := <-a.poolChan
	currentIdx := atomic.AddInt32(&a.jobSize, 1) - 1
	// add a new element for later use
	a.errList = append(a.errList, nil)
	go func() {
		defer func() {
			a.poolChan <- ticket            //release pool ticket
			a.jobChan <- ticket             // done job.
			atomic.AddInt32(&a.jobSize, -1) // decrease job size
			if err := recover(); err != nil {
				debug.PrintStack()
				a.errList[currentIdx] = fmt.Errorf("PoolAsync.DoWitError %v panic: %v", currentIdx, err)
			}
		}()
		err := f()
		if err != nil {
			// log error
			//fmt.Printf("PoolAsync.Do called function return error: %v\n", err)
		}
		a.errList[currentIdx] = err
	}()
	return a
}

// wait will block until all jobs done. and return the first none-nil error.
//
// IMPORTANT: Wait must be called after ALL Do/DoWitError
func (a *PoolAsync) Wait() error {
	var i int32
	// should get size before loop
	loopCount := a.GetCurrentJobSize()
	for i = 0; i < loopCount; i++ {
		select {
		case <-a.jobChan:
		}
	}
	for _, er := range a.errList {
		if er != nil {
			return er
		}
	}
	return nil
}

func (a *PoolAsync) GetCurrentJobSize() int32 {
	return atomic.LoadInt32(&a.jobSize)
}

func (a *PoolAsync) GetErrors() []error {
	result := make([]error, 0)
	for _, v := range a.errList {
		result = append(result, v)
	}
	return result
}
