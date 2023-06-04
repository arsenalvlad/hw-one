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
	var sm sync.Mutex

	out := make(Bi)
	mergeCh := make([]Out, 0, len(in))
	inin := make(Bi, 1)

	waitCh := make(Bi)

	go func() {
		select {
		case <-done:
			close(out)
		case <-waitCh:
			close(out)
		}
	}()

	go func() {
		for ko := range in {
			wg.Add(1)
			inin <- ko
			go func(sm *sync.Mutex) {
				sm.Lock()
				mergeCh = append(mergeCh, fn1(inin, stages...))
				wg.Done()
				sm.Unlock()
			}(&sm)
		}

		wg.Wait()

		for _, p := range mergeCh {
			out <- <-p
		}

		close(waitCh)
	}()

	return out
}

func fn1(in In, stages ...Stage) Out {
	if len(stages) == 1 {
		return stages[0](in)
	}

	return fn1(stages[0](in), stages[1:]...)
}
