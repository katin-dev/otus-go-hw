package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		stageChan := make(Bi)
		go func(in In) {
			defer close(stageChan)
			for v := range in {
				select {
				case <-done:
					return
				default:
					stageChan <- v
				}
			}
		}(stage(in))

		in = stageChan
	}

	return in
}
