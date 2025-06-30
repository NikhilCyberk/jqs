package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/NikhilCyberk/jqs/internal/models"
	"github.com/NikhilCyberk/jqs/internal/services"
	"github.com/gin-gonic/gin"
)

type mockDB struct{}

type mockWorkerPool struct{}

func setupTestDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock db: %v", err)
	}
	return db, mock
}

func TestSubmitJob_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _ := setupTestDB(t)
	wp := services.NewWorkerPool(1, func(job models.Job) {})
	Init(db, wp)
	r := gin.Default()
	r.POST("/jobs", SubmitJob)
	w := httptest.NewRecorder()
	body := bytes.NewBufferString(`invalid-json`)
	req, _ := http.NewRequest("POST", "/jobs", body)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestSubmitJob_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := setupTestDB(t)
	wp := services.NewWorkerPool(1, func(job models.Job) {})
	Init(db, wp)
	r := gin.Default()
	r.POST("/jobs", SubmitJob)
	w := httptest.NewRecorder()
	payload := `{"task": "test"}`
	now := time.Now()
	mock.ExpectQuery("INSERT INTO jobs").
		WithArgs(sqlmock.AnyArg(), "queued").
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(1, now, now))
	req, _ := http.NewRequest("POST", "/jobs", bytes.NewBufferString(payload))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
}

func TestGetJob_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := setupTestDB(t)
	Init(db, nil)
	r := gin.Default()
	r.GET("/jobs/:id", GetJob)
	w := httptest.NewRecorder()
	mock.ExpectQuery(`SELECT id, payload, status, result, created_at, updated_at FROM jobs WHERE id = \$1`).
		WithArgs("123").
		WillReturnError(sql.ErrNoRows)
	req, _ := http.NewRequest("GET", "/jobs/123", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestGetJob_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := setupTestDB(t)
	Init(db, nil)
	r := gin.Default()
	r.GET("/jobs/:id", GetJob)
	w := httptest.NewRecorder()
	now := time.Now()
	payload := []byte(`{"task":"test"}`)
	result := []byte(`{"message":"done"}`)
	mock.ExpectQuery(`SELECT id, payload, status, result, created_at, updated_at FROM jobs WHERE id = \$1`).
		WithArgs("1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "payload", "status", "result", "created_at", "updated_at"}).
			AddRow(1, payload, "completed", result, now, now))
	req, _ := http.NewRequest("GET", "/jobs/1", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var job map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &job)
	if job["status"] != "completed" {
		t.Errorf("expected status completed, got %v", job["status"])
	}
}

func TestListJobs_Empty(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := setupTestDB(t)
	Init(db, nil)
	r := gin.Default()
	r.GET("/jobs", ListJobs)
	w := httptest.NewRecorder()
	mock.ExpectQuery(`SELECT id, payload, status, result, created_at, updated_at FROM jobs ORDER BY created_at DESC LIMIT \$1 OFFSET \$2`).
		WithArgs(10, 0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "payload", "status", "result", "created_at", "updated_at"}))
	req, _ := http.NewRequest("GET", "/jobs", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var jobs []interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &jobs)
	if len(jobs) != 0 {
		t.Errorf("expected 0 jobs, got %d", len(jobs))
	}
}
