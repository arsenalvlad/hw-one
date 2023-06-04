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
	mergeCh := make([]Out, 0, len(in))
	inin := make(Bi, 1)

	go func() {
		<-done
		close(inin)
		close(out)
	}()

	for ko := range in {
		inin <- ko
		wg.Add(1)
		go func() {
			mergeCh = append(mergeCh, fn1(inin, &done, stages...))
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		for _, p := range mergeCh {
			out <- <-p
		}
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
