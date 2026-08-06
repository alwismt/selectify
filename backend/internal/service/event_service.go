package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/google/uuid"

	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/repo"
	"alwis.dev/selectify/internal/types"
)

const eventEnvelopeVersion = "1.0"

type EventEnvelope struct {
	EventID       string          `json:"event_id"`
	EventType     types.EventType `json:"event_type"`
	Version       string          `json:"version"`
	OccurredAt    time.Time       `json:"occurred_at"`
	CorrelationID string          `json:"correlation_id"`
	Payload       json.RawMessage `json:"payload"`
}

type EventPublisher interface {
	Publish(ctx context.Context, eventType types.EventType, payload any) error
}

type sqsEventPublisher struct {
	client    *sqs.Client
	queueURL  string
	eventRepo repo.EventRepo
}

func NewSQSEventPublisher(eventRepo repo.EventRepo) EventPublisher {
	queueURL := os.Getenv("API_SQS_QUEUE_URL")
	region := os.Getenv("API_SQS_REGION")
	accessKeyID := os.Getenv("API_SQS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("API_SQS_SECRET_ACCESS_KEY")

	if queueURL == "" {
		panic("API_SQS_QUEUE_URL is required")
	}
	if region == "" {
		panic("API_SQS_REGION is required")
	}
	if accessKeyID == "" {
		panic("API_SQS_ACCESS_KEY_ID is required")
	}
	if secretAccessKey == "" {
		panic("API_SQS_SECRET_ACCESS_KEY is required")
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

	return &sqsEventPublisher{
		client:    sqs.NewFromConfig(awsConfig),
		queueURL:  queueURL,
		eventRepo: eventRepo,
	}
}

func (p *sqsEventPublisher) Publish(ctx context.Context, eventType types.EventType, payload any) error {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return logger.Error(ctx, err, "failed to marshal event payload")
	}

	occurredAt := time.Now().UTC()
	envelope := EventEnvelope{
		EventID:       uuid.New().String(),
		EventType:     eventType,
		Version:       eventEnvelopeVersion,
		OccurredAt:    occurredAt,
		CorrelationID: uuid.New().String(),
		Payload:       payloadJSON,
	}

	event := &model.Event{
		ID: envelope.EventID,
		Data: &model.EventData{
			Type:    eventType,
			Payload: payloadJSON,
			Date:    &occurredAt,
		},
		ReceivedDate: occurredAt,
	}

	if err = p.eventRepo.InsertEvent(ctx, event); err != nil {
		return logger.Error(ctx, err, "failed to persist event before SQS publish")
	}

	body, err := json.Marshal(envelope)
	if err != nil {
		return logger.Error(ctx, err, "failed to marshal event envelope")
	}

	_, err = p.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(p.queueURL),
		MessageBody: aws.String(string(body)),
	})
	if err != nil {
		return logger.Errorf(ctx, err, "failed to send event %q to SQS", envelope.EventID)
	}
	logger.Infof(ctx, "sent event %q to SQS", envelope.EventID)
	return nil
}
