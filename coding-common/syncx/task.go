package syncx

import (
	"fmt"
	"sync"
)

type TaskRunner struct {
	maxConcurrent int
	taskQueue     chan func()
	wg            sync.WaitGroup
	stopChan      chan struct{}
}

// maxConcurrent: 最大并发数
func NewTaskRunner(maxConcurrent, queueSize int) *TaskRunner {
	return &TaskRunner{
		maxConcurrent: maxConcurrent,
		taskQueue:     make(chan func(), queueSize),
		stopChan:      make(chan struct{}),
	}
}

func (ts *TaskRunner) Start() {
	for range ts.maxConcurrent {
		ts.wg.Add(1)
		go ts.worker()
	}
}

func (ts *TaskRunner) worker() {
	defer ts.wg.Done()
	for {
		select {
		case task, ok := <-ts.taskQueue:
			if !ok {
				return
			}
			task()
		case <-ts.stopChan:
			return
		}
	}
}

func (ts *TaskRunner) AddTask(task func()) error {
	select {
	case ts.taskQueue <- task:
		return nil
	default:
		return fmt.Errorf("任务队列已满")
	}
}

func (ts *TaskRunner) Stop() {
	close(ts.stopChan)
	ts.wg.Wait()
}
