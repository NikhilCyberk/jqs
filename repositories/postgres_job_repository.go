package repositories

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/yourusername/jqs/models"
)

type PostgresJobRepository struct {
	DB *sql.DB
}

func NewPostgresJobRepository(db *sql.DB) *PostgresJobRepository {
	return &PostgresJobRepository{DB: db}
}

func (r *PostgresJobRepository) CreateJob(ctx context.Context, job *models.Job) error {
	return r.DB.QueryRowContext(ctx,
		`INSERT INTO jobs (payload, status) VALUES ($1, $2) RETURNING id, created_at, updated_at`,
		job.Payload, job.Status,
	).Scan(&job.ID, &job.CreatedAt, &job.UpdatedAt)
}

func (r *PostgresJobRepository) GetJobByID(ctx context.Context, id int64) (*models.Job, error) {
	var job models.Job
	var result sql.NullString
	err := r.DB.QueryRowContext(ctx,
		`SELECT id, payload, status, result, created_at, updated_at FROM jobs WHERE id = $1`, id,
	).Scan(&job.ID, &job.Payload, &job.Status, &result, &job.CreatedAt, &job.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if result.Valid {
		job.Result = json.RawMessage(result.String)
	}
	return &job, nil
}

func (r *PostgresJobRepository) ListJobs(ctx context.Context, page, limit int) ([]models.Job, error) {
	offset := (page - 1) * limit
	rows, err := r.DB.QueryContext(ctx,
		`SELECT id, payload, status, result, created_at, updated_at FROM jobs ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	jobs := []models.Job{}
	for rows.Next() {
		var job models.Job
		var result sql.NullString
		if err := rows.Scan(&job.ID, &job.Payload, &job.Status, &result, &job.CreatedAt, &job.UpdatedAt); err != nil {
			continue
		}
		if result.Valid {
			job.Result = json.RawMessage(result.String)
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (r *PostgresJobRepository) UpdateJobStatusAndResult(ctx context.Context, id int64, status string, result []byte) error {
	_, err := r.DB.ExecContext(ctx,
		`UPDATE jobs SET status = $1, result = $2, updated_at = NOW() WHERE id = $3`, status, result, id)
	return err
}
