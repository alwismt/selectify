package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"

	"alwis.dev/selectify/event-worker-pkg/handlers"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/repo"
)

type Worker interface {
	Run(ctx context.Context)
}

type sqsWorker struct {
	client       *sqs.Client
	queueURL     string
	waitTime     int32
	maxMessages  int32
	pollInterval time.Duration
	eventRepo    repo.EventRepo
}

func NewSQSWorker(eventRepo repo.EventRepo) Worker {
	queueURL := os.Getenv("EVT_SQS_QUEUE_URL")
	region := os.Getenv("EVT_SQS_REGION")
	accessKeyID := os.Getenv("EVT_SQS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("EVT_SQS_SECRET_ACCESS_KEY")

	if queueURL == "" {
		panic("EVT_SQS_QUEUE_URL is required")
	}
	if region == "" {
		panic("EVT_SQS_REGION is required")
	}
	if accessKeyID == "" {
		panic("EVT_SQS_ACCESS_KEY_ID is required")
	}
	if secretAccessKey == "" {
		panic("EVT_SQS_SECRET_ACCESS_KEY is required")
	}

	waitTime := int32(0)
	if v := os.Getenv("EVT_SQS_WAIT_TIME_SECONDS"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 || n > 20 {
			panic(fmt.Sprintf("EVT_SQS_WAIT_TIME_SECONDS must be 0-20, got %q", v))
		}
		waitTime = int32(n)
	}

	maxMessages := int32(10)
	if v := os.Getenv("EVT_SQS_MAX_MESSAGES"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 || n > 10 {
			panic(fmt.Sprintf("EVT_SQS_MAX_MESSAGES must be 1-10, got %q", v))
		}
		maxMessages = int32(n)
	}

	pollIntervalSec := 30
	if v := os.Getenv("EVT_SQS_POLL_INTERVAL_SECONDS"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 {
			panic(fmt.Sprintf("EVT_SQS_POLL_INTERVAL_SECONDS must be >= 1, got %q", v))
		}
		pollIntervalSec = n
	}

	awsConfig, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, ""),
		),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to load SQS configuration: %v", err))
	}

	return &sqsWorker{
		client:       sqs.NewFromConfig(awsConfig),
		queueURL:     queueURL,
		waitTime:     waitTime,
		maxMessages:  maxMessages,
		pollInterval: time.Duration(pollIntervalSec) * time.Second,
		eventRepo:    eventRepo,
	}
}

func (w *sqsWorker) Run(ctx context.Context) {
	logger.Info(ctx, "event worker started")
	for {
		select {
		case <-ctx.Done():
			logger.Info(ctx, "event worker shutting down")
			return
		default:
		}

		logger.Infof(ctx, "polling SQS queue (interval=%ds wait=%ds max=%d)",
			int(w.pollInterval.Seconds()), w.waitTime, w.maxMessages)

		out, err := w.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(w.queueURL),
			MaxNumberOfMessages: w.maxMessages,
			WaitTimeSeconds:     w.waitTime,
		})
		if err != nil {
			if ctx.Err() != nil {
				logger.Info(ctx, "event worker shutting down")
				return
			}
			_ = logger.Error(ctx, err, "failed to receive SQS messages")
			if !w.sleep(ctx) {
				return
			}
			continue
		}

		if len(out.Messages) == 0 {
			logger.Info(ctx, "poll complete: no messages")
		} else {
			logger.Infof(ctx, "poll complete: received %d message(s)", len(out.Messages))
		}

		for _, msg := range out.Messages {
			if err := w.processMessage(ctx, msg); err != nil {
				_ = logger.Error(ctx, err, "failed to process SQS message")
			}
		}

		if !w.sleep(ctx) {
			return
		}
	}
}

func (w *sqsWorker) sleep(ctx context.Context) bool {
	timer := time.NewTimer(w.pollInterval)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		logger.Info(ctx, "event worker shutting down")
		return false
	case <-timer.C:
		return true
	}
}

func (w *sqsWorker) processMessage(ctx context.Context, msg types.Message) error {
	messageID := aws.ToString(msg.MessageId)
	body := aws.ToString(msg.Body)

	var envelope model.EventEnvelope
	if err := json.Unmarshal([]byte(body), &envelope); err != nil {
		return logger.Error(ctx, err, "failed to unmarshal event envelope")
	}

	logger.Infof(ctx, "event received: event_id=%s type=%s message_id=%s",
		envelope.EventID, envelope.EventType, messageID)

	event, err := w.eventRepo.GetByID(ctx, envelope.EventID)
	if err != nil {
		return err
	}
	if event == nil {
		return fmt.Errorf("event %q not found in database", envelope.EventID)
	}

	eventHandler, ok := handlers.HandlerRegistry[envelope.EventType]
	if !ok {
		return fmt.Errorf("no handler registered for event type %q", envelope.EventType)
	}

	if err = eventHandler.Handle(ctx, event); err != nil {
		return logger.Errorf(ctx, err, "handler failed for event %q", envelope.EventID)
	}

	processedAt := time.Now().UTC()
	if err = w.eventRepo.MarkProcessed(ctx, envelope.EventID, processedAt); err != nil {
		return err
	}

	_, err = w.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(w.queueURL),
		ReceiptHandle: msg.ReceiptHandle,
	})
	if err != nil {
		return logger.Errorf(ctx, err, "failed to delete SQS message for event %q", envelope.EventID)
	}

	logger.Infof(ctx, "event processed: event_id=%s type=%s", envelope.EventID, envelope.EventType)
	return nil
}
