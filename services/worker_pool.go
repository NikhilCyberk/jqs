package services

import (
	"github.com/yourusername/jqs/models"
	"github.com/yourusername/jqs/utils"
)

// The JobHandler should update job status and result in the database.
// Example job processing logic:
// 1. Update job status to 'processing' in DB
// 2. Perform the job (simulate with sleep or actual logic)
// 3. Update job status to 'completed' and set result in DB
type JobHandler func(job models.Job)

type WorkerPool struct {
	JobQueue   chan models.Job
	NumWorkers int
	Handler    JobHandler
}

func NewWorkerPool(numWorkers int, handler JobHandler) *WorkerPool {
	return &WorkerPool{
		JobQueue:   make(chan models.Job, 100),
		NumWorkers: numWorkers,
		Handler:    handler,
	}
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.NumWorkers; i++ {
		go func(workerID int) {
			for job := range wp.JobQueue {
				utils.Logger.WithField("worker", workerID).Infof("Processing job %d", job.ID)
				wp.Handler(job)
			}
		}(i)
	}
}

func (wp *WorkerPool) Submit(job models.Job) {
	wp.JobQueue <- job
}
