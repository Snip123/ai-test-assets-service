package application

import (
	"context"
	"fmt"
	"time"

	"github.com/Snip123/ai-test-assets-service/internal/domain"
)

// DecommissionAssetCommand carries the input for decommissioning an Asset.
type DecommissionAssetCommand struct {
	TenantID string
	ID       string
	Reason   string
}

// DecommissionAsset transitions an Active Asset to Decommissioned and publishes
// an AssetDecommissioned domain event.
func (s *AssetService) DecommissionAsset(ctx context.Context, cmd DecommissionAssetCommand) (domain.Asset, error) {
	existing, err := s.GetAsset(ctx, cmd.ID, cmd.TenantID)
	if err != nil {
		return domain.Asset{}, err
	}

	if existing.Status == domain.AssetStatusDecommissioned {
		return domain.Asset{}, ErrAssetAlreadyDecommissioned
	}

	decommissioned, err := s.repo.Decommission(ctx, cmd.ID, cmd.TenantID)
	if err != nil {
		return domain.Asset{}, fmt.Errorf("decommission asset: %w", err)
	}

	if err := s.publisher.PublishAssetDecommissioned(ctx, domain.AssetDecommissionedEvent{
		AssetID:            decommissioned.ID,
		TenantID:           decommissioned.TenantID,
		DecommissionedDate: time.Now().UTC(),
		Reason:             cmd.Reason,
	}); err != nil {
		_ = err // best-effort; TODO: outbox pattern
	}

	return decommissioned, nil
}
