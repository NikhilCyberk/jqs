package repositories

import (
	"context"

	"github.com/yourusername/jqs/models"
)

type JobRepository interface {
	CreateJob(ctx context.Context, job *models.Job) error
	GetJobByID(ctx context.Context, id int64) (*models.Job, error)
	ListJobs(ctx context.Context, page, limit int) ([]models.Job, error)
	UpdateJobStatusAndResult(ctx context.Context, id int64, status string, result []byte) error
}
