# conveyor
Minimal Unix-like CI Runner

## REST API

```
POST /job -- Queue a job
```

```
GET /job/<job_id> -- Get status of job.
```

```
DELETE /job/<job_id> -- Stop / remove job from queue.
```

```
POST /job/<job_id> -- Restart / reinsert job into queue.
```

```
GET /job/<job_id>/log -- Get log file of job.
```

## Testing

```
curl -H "Content-Type: application/json" -d '{"name":"frontend","commands":["echo 1","echo 2","touch file_$RANDOM","sleep 15"]}' http://localhost:8080/job
```