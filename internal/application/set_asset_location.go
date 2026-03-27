package application

import (
	"context"
	"fmt"

	"github.com/Snip123/ai-test-assets-service/internal/domain"
)

// SetAssetLocationCommand carries the input for setting or updating an Asset's Location.
type SetAssetLocationCommand struct {
	TenantID   string
	ID         string
	FacilityID string
	LocationID string
}

// SetAssetLocation sets or updates the Location of an Asset and publishes an
// AssetLocationSet domain event.
func (s *AssetService) SetAssetLocation(ctx context.Context, cmd SetAssetLocationCommand) (domain.Asset, error) {
	if cmd.FacilityID == "" {
		return domain.Asset{}, fmt.Errorf("facility_id is required")
	}
	if cmd.LocationID == "" {
		return domain.Asset{}, fmt.Errorf("location_id is required")
	}

	// Confirm the Asset exists and belongs to this Tenant.
	if _, err := s.GetAsset(ctx, cmd.ID, cmd.TenantID); err != nil {
		return domain.Asset{}, err
	}

	updated, err := s.repo.SetLocation(ctx, cmd.ID, cmd.TenantID, cmd.FacilityID, cmd.LocationID)
	if err != nil {
		return domain.Asset{}, fmt.Errorf("set asset location: %w", err)
	}

	if err := s.publisher.PublishAssetLocationSet(ctx, domain.AssetLocationSetEvent{
		AssetID:    updated.ID,
		TenantID:   updated.TenantID,
		FacilityID: updated.FacilityID,
		LocationID: updated.LocationID,
	}); err != nil {
		_ = err // best-effort; TODO: outbox pattern
	}

	return updated, nil
}
