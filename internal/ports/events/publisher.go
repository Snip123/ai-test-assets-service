package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Snip123/ai-test-assets-service/internal/domain"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/oklog/ulid/v2"
	"go.opentelemetry.io/otel/trace"
)

// Envelope is the standard domain event wrapper (ADR-0016).
type Envelope struct {
	ID         string          `json:"id"`
	Type       string          `json:"type"`
	Source     string          `json:"source"`
	Version    string          `json:"version"`
	TenantID   string          `json:"tenant_id"`
	OccurredAt time.Time       `json:"occurred_at"`
	TraceID    string          `json:"trace_id"`
	Data       json.RawMessage `json:"data"`
}

const source = "ai-test-assets-service"

type Publisher struct {
	js jetstream.JetStream
}

func NewPublisher(js jetstream.JetStream) *Publisher {
	return &Publisher{js: js}
}

func (p *Publisher) PublishAssetRegistered(ctx context.Context, event domain.AssetRegisteredEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal AssetRegistered: %w", err)
	}

	env := Envelope{
		ID:         ulid.Make().String(),
		Type:       "AssetRegistered",
		Source:     source,
		Version:    "1",
		TenantID:   event.TenantID,
		OccurredAt: time.Now().UTC(),
		TraceID:    traceID(ctx),
		Data:       data,
	}

	msg, err := json.Marshal(env)
	if err != nil {
		return fmt.Errorf("marshal envelope: %w", err)
	}

	subject := fmt.Sprintf("fsi.%s.assets.AssetRegistered", event.TenantID)
	if _, err := p.js.Publish(ctx, subject, msg); err != nil {
		return fmt.Errorf("publish AssetRegistered: %w", err)
	}
	return nil
}

func traceID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		return span.SpanContext().TraceID().String()
	}
	return "00000000000000000000000000000000"
}
