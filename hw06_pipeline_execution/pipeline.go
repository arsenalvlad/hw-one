package hw06pipelineexecution

import (
	"sync"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	var wg sync.WaitGroup

	out := make(Bi)
	inin := make(Bi, 1)

	for ko := range in {
		select {
		case <-done:
			return nil
		default:
			inin <- ko
			wg.Add(1)
			go func() {
				out <- <-fn1(inin, &done, stages...)
				wg.Done()
			}()
		}

	}

	go func() {
		wg.Wait()
		close(inin)
		close(out)
	}()

	return out
}

func fn1(in In, done *In, stages ...Stage) Out {
	select {
	case <-*done:
		return nil
	default:
		if len(stages) == 1 {
			return stages[0](in)
		}

		return fn1(stages[0](in), done, stages[1:]...)
	}
}
