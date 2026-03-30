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
	return p.publish(ctx, event.TenantID, "AssetRegistered", "1", event)
}

func (p *Publisher) PublishAssetAttributesUpdated(ctx context.Context, event domain.AssetAttributesUpdatedEvent) error {
	return p.publish(ctx, event.TenantID, "AssetAttributesUpdated", "1", event)
}

func (p *Publisher) PublishAssetDecommissioned(ctx context.Context, event domain.AssetDecommissionedEvent) error {
	return p.publish(ctx, event.TenantID, "AssetDecommissioned", "1", event)
}

func (p *Publisher) PublishAssetLocationSet(ctx context.Context, event domain.AssetLocationSetEvent) error {
	return p.publish(ctx, event.TenantID, "AssetLocationSet", "1", event)
}

func (p *Publisher) publish(ctx context.Context, tenantID, eventType, version string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal %s: %w", eventType, err)
	}

	env := Envelope{
		ID:         ulid.Make().String(),
		Type:       eventType,
		Source:     source,
		Version:    version,
		TenantID:   tenantID,
		OccurredAt: time.Now().UTC(),
		TraceID:    traceID(ctx),
		Data:       data,
	}

	msg, err := json.Marshal(env)
	if err != nil {
		return fmt.Errorf("marshal envelope: %w", err)
	}

	subject := fmt.Sprintf("platform.%s.assets.%s", tenantID, eventType)
	if _, err := p.js.Publish(ctx, subject, msg); err != nil {
		return fmt.Errorf("publish %s: %w", eventType, err)
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
