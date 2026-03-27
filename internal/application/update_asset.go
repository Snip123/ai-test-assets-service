package application

import (
	"context"
	"fmt"

	"github.com/Snip123/ai-test-assets-service/internal/domain"
)

// UpdateAssetCommand carries the input for updating mutable Asset attributes.
type UpdateAssetCommand struct {
	TenantID     string
	ID           string
	Name         string
	SerialNumber string
}

// UpdateAsset updates mutable attributes on an existing Asset and publishes an
// AssetAttributesUpdated domain event.
func (s *AssetService) UpdateAsset(ctx context.Context, cmd UpdateAssetCommand) (domain.Asset, error) {
	if cmd.Name == "" {
		return domain.Asset{}, fmt.Errorf("name is required")
	}

	existing, err := s.GetAsset(ctx, cmd.ID, cmd.TenantID)
	if err != nil {
		return domain.Asset{}, err
	}

	updated, err := s.repo.Update(ctx, cmd.ID, cmd.TenantID, cmd.Name, cmd.SerialNumber)
	if err != nil {
		return domain.Asset{}, fmt.Errorf("update asset: %w", err)
	}

	changed := map[string]any{}
	if existing.Name != cmd.Name {
		changed["name"] = cmd.Name
	}
	if existing.SerialNumber != cmd.SerialNumber {
		changed["serial_number"] = cmd.SerialNumber
	}

	if err := s.publisher.PublishAssetAttributesUpdated(ctx, domain.AssetAttributesUpdatedEvent{
		AssetID:           updated.ID,
		TenantID:          updated.TenantID,
		ChangedAttributes: changed,
	}); err != nil {
		_ = err // best-effort; TODO: outbox pattern
	}

	return updated, nil
}
