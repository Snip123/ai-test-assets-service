package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Snip123/ai-test-assets-service/internal/domain"
	"github.com/oklog/ulid/v2"
)

var ErrAssetNotFound = errors.New("asset not found")
var ErrAssetAlreadyDecommissioned = errors.New("asset is already decommissioned")

// AssetRepository is the port this bounded context depends on for persistence.
type AssetRepository interface {
	Create(ctx context.Context, asset domain.Asset) (domain.Asset, error)
	List(ctx context.Context, tenantID string) ([]domain.Asset, error)
	Get(ctx context.Context, id, tenantID string) (domain.Asset, error)
	Update(ctx context.Context, id, tenantID, name, serialNumber string) (domain.Asset, error)
	Decommission(ctx context.Context, id, tenantID string) (domain.Asset, error)
	SetLocation(ctx context.Context, id, tenantID, facilityID, locationID string) (domain.Asset, error)
}

// EventPublisher is the port for publishing domain events to NATS JetStream.
type EventPublisher interface {
	PublishAssetRegistered(ctx context.Context, event domain.AssetRegisteredEvent) error
	PublishAssetAttributesUpdated(ctx context.Context, event domain.AssetAttributesUpdatedEvent) error
	PublishAssetDecommissioned(ctx context.Context, event domain.AssetDecommissionedEvent) error
	PublishAssetLocationSet(ctx context.Context, event domain.AssetLocationSetEvent) error
}

// AssetService handles all Asset use cases for this bounded context.
type AssetService struct {
	repo      AssetRepository
	publisher EventPublisher
}

func NewAssetService(repo AssetRepository, publisher EventPublisher) *AssetService {
	return &AssetService{repo: repo, publisher: publisher}
}

// RegisterAssetCommand carries the input for registering a new Asset.
type RegisterAssetCommand struct {
	TenantID     string
	Name         string
	AssetType    string
	FacilityID   string
	SerialNumber string
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

// GetAsset returns a single Asset by ID within a Tenant.
func (s *AssetService) GetAsset(ctx context.Context, id, tenantID string) (domain.Asset, error) {
	asset, err := s.repo.Get(ctx, id, tenantID)
	if err != nil {
		if err.Error() == "get asset: no rows in result set" {
			return domain.Asset{}, ErrAssetNotFound
		}
		return domain.Asset{}, fmt.Errorf("get asset: %w", err)
	}
	return asset, nil
}
