package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/NikhilCyberk/jqs/internal/models"
	"github.com/NikhilCyberk/jqs/internal/repositories"
	"github.com/NikhilCyberk/jqs/internal/services"
	"github.com/NikhilCyberk/jqs/internal/utils"
	"github.com/gin-gonic/gin"
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
		utils.Logger.WithError(err).Error(utils.ErrInvalidPayload)
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrInvalidPayload})
		return
	}
	job := models.Job{
		Payload: payload,
		Status:  utils.JobStatusQueued,
	}
	if err := Repo.CreateJob(c.Request.Context(), &job); err != nil {
		utils.Logger.WithError(err).Error(utils.ErrFailedSubmitJob)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.ErrFailedSubmitJob})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrInvalidJobID})
		return
	}
	job, err := Repo.GetJobByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": utils.ErrJobNotFound})
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
		utils.Logger.WithError(err).Error(utils.ErrFailedListJobs)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.ErrFailedListJobs})
		return
	}
	c.JSON(http.StatusOK, jobs)
}
