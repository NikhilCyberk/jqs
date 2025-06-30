package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"github.com/yourusername/jqs/handlers"
	"github.com/yourusername/jqs/models"
	"github.com/yourusername/jqs/repositories"
	"github.com/yourusername/jqs/services"
	"github.com/yourusername/jqs/utils"
)

func main() {
	utils.InitLogger()
	cfg := utils.LoadConfig()

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	defer db.Close()

	if err := models.Migrate(db); err != nil {
		log.Fatal("Migration failed: ", err)
	}

	repo := repositories.NewPostgresJobRepository(db)

	workerPool := services.NewWorkerPool(5, func(job models.Job) {
		// Update job status to 'processing'
		err := repo.UpdateJobStatusAndResult(context.Background(), job.ID, "processing", nil)
		if err != nil {
			utils.Logger.WithError(err).Errorf("Failed to update job %d to processing", job.ID)
			return
		}
		utils.Logger.WithField("job_id", job.ID).Info("Job set to processing")

		// Simulate job work (replace with real logic as needed)
		time.Sleep(2 * time.Second)

		// Set a dummy result (could be any JSON)
		result := []byte(`{"message": "Job completed successfully"}`)
		err = repo.UpdateJobStatusAndResult(context.Background(), job.ID, "completed", result)
		if err != nil {
			utils.Logger.WithError(err).Errorf("Failed to update job %d to completed", job.ID)
			return
		}
		utils.Logger.WithField("job_id", job.ID).Info("Job completed")
	})
	workerPool.Start()

	handlers.Init(repo, workerPool)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.POST("/jobs", handlers.SubmitJob)
	r.GET("/jobs/:id", handlers.GetJob)
	r.GET("/jobs", handlers.ListJobs)
	r.GET("/health", func(c *gin.Context) {
		c.String(200, "OK")
	})

	srv := &gin.Engine{}
	*srv = *r

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: srv,
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}
	}()

	utils.Logger.Infof("Starting server on :%s...", cfg.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
