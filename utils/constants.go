package utils

// Job status constants
const (
	JobStatusQueued     = "queued"
	JobStatusProcessing = "processing"
	JobStatusCompleted  = "completed"
	JobStatusFailed     = "failed"
)

// Error message constants
const (
	ErrInvalidPayload  = "Invalid payload"
	ErrFailedSubmitJob = "Failed to submit job"
	ErrJobNotFound     = "Job not found"
	ErrFailedListJobs  = "Failed to list jobs"
	ErrInvalidJobID    = "Invalid job ID"
)

// Success message constants
const (
	MsgJobCompleted = "Job completed successfully"
)
