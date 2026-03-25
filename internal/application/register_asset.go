package application

import (
	"context"
	"fmt"
	"time"

	"github.com/Snip123/ai-test-assets-service/internal/domain"
	"github.com/oklog/ulid/v2"
)

// AssetRepository is the port this use case depends on.
// The postgres adapter implements this interface.
type AssetRepository interface {
	Create(ctx context.Context, asset domain.Asset) (domain.Asset, error)
	List(ctx context.Context, tenantID string) ([]domain.Asset, error)
	Get(ctx context.Context, id, tenantID string) (domain.Asset, error)
}

// EventPublisher is the port for publishing domain events.
// The NATS adapter implements this interface.
type EventPublisher interface {
	PublishAssetRegistered(ctx context.Context, event domain.AssetRegisteredEvent) error
}

// RegisterAssetCommand carries the input for registering a new Asset.
type RegisterAssetCommand struct {
	TenantID     string
	Name         string
	AssetType    string
	FacilityID   string
	SerialNumber string
}

// AssetService handles Asset use cases.
type AssetService struct {
	repo      AssetRepository
	publisher EventPublisher
}

func NewAssetService(repo AssetRepository, publisher EventPublisher) *AssetService {
	return &AssetService{repo: repo, publisher: publisher}
}

// RegisterAsset creates a new Asset and publishes an AssetRegistered domain event.
func (s *AssetService) RegisterAsset(ctx context.Context, cmd RegisterAssetCommand) (domain.Asset, error) {
	if cmd.Name == "" {
		return domain.Asset{}, fmt.Errorf("name is required")
	}
	if cmd.AssetType == "" {
		return domain.Asset{}, fmt.Errorf("asset_type is required")
	}
	if cmd.FacilityID == "" {
		return domain.Asset{}, fmt.Errorf("facility_id is required")
	}

	asset := domain.Asset{
		ID:           ulid.Make().String(),
		TenantID:     cmd.TenantID,
		Name:         cmd.Name,
		AssetType:    cmd.AssetType,
		FacilityID:   cmd.FacilityID,
		SerialNumber: cmd.SerialNumber,
		Status:       domain.AssetStatusActive,
	}

	created, err := s.repo.Create(ctx, asset)
	if err != nil {
		return domain.Asset{}, fmt.Errorf("create asset: %w", err)
	}

	if err := s.publisher.PublishAssetRegistered(ctx, domain.AssetRegisteredEvent{
		AssetID:       created.ID,
		TenantID:      created.TenantID,
		AssetType:     created.AssetType,
		FacilityID:    created.FacilityID,
		InstalledDate: time.Now().UTC(),
	}); err != nil {
		// Log but don't fail — event publishing is best-effort at this stage
		// TODO: implement outbox pattern for guaranteed delivery
		_ = err
	}

	return created, nil
}

// ListAssets returns all Assets for a Tenant.
func (s *AssetService) ListAssets(ctx context.Context, tenantID string) ([]domain.Asset, error) {
	return s.repo.List(ctx, tenantID)
}
