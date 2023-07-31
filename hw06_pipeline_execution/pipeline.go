package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	inin := in

	for _, stage := range stages {
		inin = fn2(stage(inin), done)
	}

	return inin
}

func fn2(inin In, done In) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case value, ok := <-inin:
				if !ok {
					return
				}
				out <- value
			}
		}
	}()

	return out
}
