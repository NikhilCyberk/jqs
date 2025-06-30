package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/jqs/models"
	"github.com/yourusername/jqs/repositories"
	"github.com/yourusername/jqs/services"
	"github.com/yourusername/jqs/utils"
)

var Repo repositories.JobRepository
var WorkerPool *services.WorkerPool

func Init(repo repositories.JobRepository, wp *services.WorkerPool) {
	Repo = repo
	WorkerPool = wp
}

func SubmitJob(c *gin.Context) {
	var payload json.RawMessage
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.Logger.WithError(err).Error("Invalid job payload")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}
	job := models.Job{
		Payload: payload,
		Status:  "queued",
	}
	if err := Repo.CreateJob(c.Request.Context(), &job); err != nil {
		utils.Logger.WithError(err).Error("Failed to insert job")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit job"})
		return
	}
	WorkerPool.Submit(job)
	utils.Logger.WithField("job_id", job.ID).Info("Job submitted")
	c.JSON(http.StatusCreated, job)
}

func GetJob(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}
	job, err := Repo.GetJobByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}
	c.JSON(http.StatusOK, job)
}

func ListJobs(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	jobs, err := Repo.ListJobs(c.Request.Context(), page, limit)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to list jobs")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list jobs"})
		return
	}
	c.JSON(http.StatusOK, jobs)
}
