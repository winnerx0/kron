package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"github.com/winnerx0/kron/internal/config"
	"github.com/winnerx0/kron/internal/database"
	"github.com/winnerx0/kron/internal/domain"
	"github.com/winnerx0/kron/internal/execution"
	"github.com/winnerx0/kron/internal/job"
	rabbitmq "github.com/winnerx0/kron/internal/queue"
	"github.com/winnerx0/kron/internal/scheduler"
	"github.com/winnerx0/kron/internal/secret"
)

func main() {

	httpClient := &http.Client{}

	config := config.Load()

	conn := rabbitmq.NewRabbitMQClient(config.RabbitMQURL)

	err := rabbitmq.Setup(conn.Ch)

	if err != nil {
		log.Fatal("Failed to setup RabbitMQ Queues ", err)
	}

	database := database.NewDatabase(config.DBHost, config.DBUser, config.DBPassword, config.DBPort, config.DBName)

	db := database.Start()

	rdb := redis.NewClient(&redis.Options{Addr: "redis:6379", DB: 0, Protocol: 2})

	locker := scheduler.NewLocker(rdb, time.Minute*5, "")

	jobRepo := job.NewRepository(db)

	executionRepo := execution.NewPostgresRepository(db)

	msgs, err := conn.ConsumeQueue("jobs_queue")

	if err != nil {
		log.Fatal("Failed to consume queue ", err)
	}

	log.Println("Worker instance started")

	for msg := range msgs {

		var job domain.Job

		err := json.Unmarshal(msg.Body, &job)
		if err != nil {
			log.Println("Failed to unmarshal job", err)
			continue
		}

		// s.ensureActiveExecutionTracking()

		executionCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		newExecution := domain.Execution{
			ID:      uuid.NewString(),
			JobID:   job.ID,
			Status:  domain.RUNNING,
			Started: time.Now(),
		}

		executionRepo.Save(executionCtx, newExecution)

		finish := func(status domain.ExecutionStatus, responseBody string) {
			newExecution.Finished = time.Now()
			newExecution.Status = status
			newExecution.ResponseBody = responseBody
			if err := executionRepo.Update(executionCtx, newExecution); err != nil {
				log.Println("Error updating execution", err)
				return
			}

			if status == domain.FAILED {
				job.Status = false
				jobRepo.Update(executionCtx, job)
				return
			}

			scheduler, err := cron.ParseStandard(job.Schedule)

			if err != nil {
				log.Println("Error parsing schedule", err)
				return
			}
			locker.Release(executionCtx, job.ID, job.NextRunAt)
			job.NextRunAt = scheduler.Next(time.Now())
			jobRepo.Update(executionCtx, job)
		}

		select {
		case <-executionCtx.Done():
			finish(domain.STOPPED, "")
			continue
		default:
		}

		for attempt := range 5 {
			req, err := http.NewRequestWithContext(executionCtx, job.Method, job.Endpoint, bytes.NewReader([]byte(job.Body)))
			if err != nil {
				log.Println("Error creating request", err)
				finish(domain.FAILED, err.Error())
				continue
			}

			headers, err := decryptHeaders(job.Headers)
			if err != nil {
				log.Println("Error decrypting job headers", err)
				finish(domain.FAILED, err.Error())
				continue
			}

			for key, value := range headers {
				req.Header.Set(key, value)
			}

			resp, err := httpClient.Do(req)
			if err != nil {
				log.Println("Error sending request", err)
				if errors.Is(executionCtx.Err(), context.Canceled) || errors.Is(executionCtx.Err(), context.DeadlineExceeded) {
					finish(domain.STOPPED, "")
					continue
				}
				if attempt < 4 {
					select {
					case <-executionCtx.Done():
						finish(domain.STOPPED, "")
						continue
					case <-time.After(30 * time.Second):
					}
					continue
				}

				//TODO: for this add a dead letter queue to send failed jobs due to bad requests
				finish(domain.FAILED, err.Error())

				continue
			}

			body := readResponseBody(resp.Body)
			resp.Body.Close()

			if resp.StatusCode < 200 || resp.StatusCode > 299 {
				log.Printf("Error: received status code %d for job %s", resp.StatusCode, job.Name)
				if resp.StatusCode >= 500 && attempt < 4 {
					select {
					case <-executionCtx.Done():
						finish(domain.STOPPED, "")
						continue
					case <-time.After(30 * time.Second):
					}
					continue
				}
				finish(domain.FAILED, fmt.Sprintf("status %d: %s", resp.StatusCode, body))
				err := conn.Ch.Nack(msg.DeliveryTag, false, true)

				if err != nil {
					log.Println("Failed to nack message: ", err)
				}
				break
			}

			finish(domain.SUCCESS, body)
			jobRepo.Update(executionCtx, job)
			err = conn.Ch.Ack(msg.DeliveryTag, false)
			if err != nil {
				log.Println("Failed to ack message: ", err)
			}

			break
		}
	}
}

func readResponseBody(r io.Reader) string {
	const maxBodyBytes = 64 * 1024
	data, err := io.ReadAll(io.LimitReader(r, maxBodyBytes))
	if err != nil {
		return ""
	}
	return string(data)
}

func decryptJobSecrets(job domain.Job) (domain.Job, error) {
	headers, err := decryptHeaders(job.Headers)
	if err != nil {
		return domain.Job{}, err
	}

	decryptedHeaders := make(map[string]any, len(headers))
	for key, value := range headers {
		decryptedHeaders[key] = value
	}

	job.Headers = decryptedHeaders
	return job, nil
}

func decryptHeaders(headers map[string]any) (map[string]string, error) {
	decrypted := make(map[string]string, len(headers))
	for key, value := range headers {
		valueString := fmt.Sprintf("%v", value)
		if !secret.IsSensitiveHeader(key) {
			decrypted[key] = valueString
			continue
		}

		decryptedValue, err := secretManager().Decrypt(valueString)
		if err != nil {
			return nil, err
		}
		decrypted[key] = decryptedValue
	}
	return decrypted, nil
}

func secretManager() secret.Manager {
	return secret.NoopManager{}
}

func encryptJobSecrets(job domain.Job, existingHeaders map[string]any) (domain.Job, error) {
	headers := make(map[string]any, len(job.Headers))
	for key, value := range job.Headers {
		valueString := fmt.Sprintf("%v", value)
		if !secret.IsSensitiveHeader(key) {
			headers[key] = value
			continue
		}

		if secret.IsMasked(valueString) && existingHeaders != nil {
			if existingValue, ok := existingHeaders[key]; ok {
				headers[key] = existingValue
				continue
			}
		}

		encryptedValue, err := secretManager().Encrypt(valueString)
		if err != nil {
			return domain.Job{}, err
		}
		headers[key] = encryptedValue
	}

	job.Headers = headers
	return job, nil
}
