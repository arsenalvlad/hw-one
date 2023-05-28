package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) (err error) {
	var wg sync.WaitGroup

	taskCh := make(chan Task)
	errCh := make(chan struct{}, m)
	quitCh := make(chan bool, n)
	defer close(taskCh)
	defer close(errCh)
	defer close(quitCh)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go Worker(m, &wg, taskCh, errCh, quitCh)
	}

	for index, item := range tasks {
		if n+m == index {
			if len(errCh) == m {
				err = ErrErrorsLimitExceeded
				break
			}
		}

		taskCh <- item
	}

	for i := 0; i < n; i++ {
		quitCh <- true
	}

	wg.Wait()
	if err != nil {
		return err
	}

	return nil
}

func Worker(m int, wg *sync.WaitGroup, jobs chan Task, errCh chan struct{}, quitCh chan bool) {
	defer wg.Done()
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}

			err := job()
			if err != nil {
				if len(errCh) == m {
					return
				}
				errCh <- struct{}{}
			}
		case <-quitCh:
			return
		}
	}
}
