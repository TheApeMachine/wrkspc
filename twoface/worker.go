package twoface

type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	disposer   Context
}

func NewWorker(
	workerPool chan chan Job,
	disposer Context,
) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		disposer:   disposer,
	}
}

func (worker Worker) Start() Worker {
	go func() {
		defer close(worker.JobChannel)

		for {
			worker.WorkerPool <- worker.JobChannel

			select {
			case job := <-worker.JobChannel:
				job.Do()
			case <-worker.disposer.Done():
				return
			}
		}
	}()

	return worker
}

func (worker Worker) Stop() {
}
