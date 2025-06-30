# Job Queue System (JQS)

A high-performance, asynchronous job queue system built in Go using Gin, PostgreSQL, and a goroutine-based worker pool. Designed for cloud-native deployment (Docker, Render.com, etc).

---

## Features
- **REST API** for job submission, status, and listing
- **Asynchronous processing** with a configurable worker pool (default: 5 workers)
- **Structured logging** with logrus
- **PostgreSQL** for job persistence
- **Environment variable** and `.env` file support
- **Containerized** with Docker (non-root user)
- **Graceful shutdown** and robust error handling

---

## Directory Structure
```
JQS/
├── cmd/                # Main application entrypoint
├── handlers/           # HTTP handlers (Gin)
├── services/           # Worker pool and business logic
├── models/             # Database models and migrations
├── utils/              # Logging, config, and helpers
├── Dockerfile
├── go.mod
├── go.sum
├── .env.example
└── README.md
```

---

## Database
- **PostgreSQL**
- Table: `jobs`
  - `id` (serial, primary key)
  - `payload` (JSONB, job data)
  - `status` (string: queued, processing, completed, failed)
  - `result` (JSONB, job result or error)
  - `created_at`, `updated_at` (timestamps)

---

## Environment Variables
You can use a `.env` file (see `.env.example`):
```
DATABASE_URL=postgres://postgres:mysecretpassword@localhost:5432/mydb?sslmode=disable
PORT=8080
```

---

## Setup & Run
1. **Clone and build**
   ```sh
   git clone <repo-url>
   cd JQS
   go mod tidy
   go build -o jqs ./cmd/main.go
   ```
2. **Configure environment**
   - Copy `.env.example` to `.env` and fill in your values
3. **Run**
   ```sh
   go run cmd/main.go
   # or
   ./jqs
   ```

---

## Docker
1. **Build image**
   ```sh
   docker build -t jqs .
   ```
2. **Run container**
   ```sh
   docker run --env-file .env -p 8080:8080 jqs
   ```

---

## Deployment on Render.com

Render makes it easy to deploy Dockerized web services and managed PostgreSQL databases.

### 1. **Provision a PostgreSQL Database**
- In the Render dashboard, create a new PostgreSQL database.
- Copy the connection string (it will be your `DATABASE_URL`).

### 2. **Create a New Web Service**
- Click "New Web Service" and connect your GitHub repo.
- Select "Docker" as the environment.
- Set the build and start commands (Render will use your Dockerfile automatically).
- Set environment variables:
  - `DATABASE_URL` (from your Render database)
  - `PORT` (Render uses 10000 by default, or set to 8080 and update the Render port setting)

### 3. **Expose the Port**
- Make sure your service is set to listen on the port specified by the `PORT` environment variable.

### 4. **Deploy**
- Click "Create Web Service". Render will build and deploy your app.
- You can view logs and manage redeploys from the Render dashboard.

---

## API Endpoints

### Submit a Job
- **POST** `/jobs`
- **Body:** Any JSON object (job payload)
- **Example:**
  ```json
  {"task": "send_email", "to": "user@example.com", "subject": "Hello"}
  ```
- **Response:**
  ```json
  {
    "id": 1,
    "payload": {"task": "send_email", ...},
    "status": "queued",
    "result": null,
    "created_at": "...",
    "updated_at": "..."
  }
  ```

### Get Job Status/Result
- **GET** `/jobs/{id}`
- **Response:**
  ```json
  {
    "id": 1,
    "payload": {"task": "send_email", ...},
    "status": "completed",
    "result": {"message": "Job completed successfully"},
    "created_at": "...",
    "updated_at": "..."
  }
  ```

### List Jobs
- **GET** `/jobs?page=1&limit=10`
- **Response:** Array of job objects

### Health Check
- **GET** `/health`
- **Response:** `OK`

---

## Worker Pool
- **Concurrency:** Default 5 workers (configurable in code)
- **Behavior:**
  - Picks up jobs from the queue
  - Updates status to `processing`
  - Simulates work (replace with your logic)
  - Updates status to `completed` and sets result
- **Logging:** All steps are logged with job ID and status

---

## Example Job Payloads
Here are some example JSON payloads you can POST to `/jobs`:
```json
{"task": "send_email", "to": "user1@example.com", "subject": "Welcome!", "body": "Hello User1"}
{"task": "resize_image", "image_url": "https://example.com/image1.jpg", "width": 800, "height": 600}
{"task": "generate_report", "report_type": "sales", "date": "2025-06-30"}
{"task": "backup_database", "db_name": "prod_db", "destination": "s3://backups/prod_db.bak"}
{"task": "notify_user", "user_id": 42, "message": "Your order has shipped."}
{"task": "transcode_video", "video_id": "abc123", "format": "mp4"}
{"task": "fetch_url", "url": "https://api.example.com/data", "method": "GET"}
{"task": "update_cache", "key": "homepage", "value": "new content"}
{"task": "process_payment", "user_id": 99, "amount": 49.99, "currency": "USD"}
{"task": "send_sms", "phone": "+1234567890", "message": "Your code is 123456"}
```

---

## Development & Testing
- **Unit tests:** See `handlers/job_handlers_test.go`
- **Run tests:**
  ```sh
  go test ./...
  ```
- **Logging:** Structured logs (JSON) for all job events

---

## License
MIT 