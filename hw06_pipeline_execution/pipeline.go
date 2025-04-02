package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		return in
	}
	for _, stage := range stages {
		stageOut := stage(in)
		out := make(Bi, 1)
		go func(stageOut In, out Bi) {
			defer close(out)
			for {
				select {
				case <-done:
					go func() {
						for range stageOut {
						}
					}()
					return
				case v, ok := <-stageOut:
					if !ok {
						return
					}
					select {
					case out <- v:
					case <-done:
						go func() {
							for range stageOut {
							}
						}()
						return
					}
				}
			}
		}(stageOut, out)
		in = out
	}
	return in
}
