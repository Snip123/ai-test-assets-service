package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Snip123/ai-test-assets-service/internal/domain"
	db "github.com/Snip123/ai-test-assets-service/internal/generated/db"
)

type AssetRepo struct {
	sqlDB   *sql.DB
	queries *db.Queries
}

func NewAssetRepo(sqlDB *sql.DB) *AssetRepo {
	return &AssetRepo{sqlDB: sqlDB, queries: db.New(sqlDB)}
}

func (r *AssetRepo) Create(ctx context.Context, a domain.Asset) (domain.Asset, error) {
	row, err := r.queries.CreateAsset(ctx, db.CreateAssetParams{
		ID:           a.ID,
		TenantID:     a.TenantID,
		Name:         a.Name,
		AssetType:    a.AssetType,
		FacilityID:   a.FacilityID,
		SerialNumber: sql.NullString{String: a.SerialNumber, Valid: a.SerialNumber != ""},
		Status:       string(a.Status),
	})
	if err != nil {
		return domain.Asset{}, fmt.Errorf("insert asset: %w", err)
	}
	return toDomain(row), nil
}

func (r *AssetRepo) List(ctx context.Context, tenantID string) ([]domain.Asset, error) {
	rows, err := r.queries.ListAssets(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("list assets: %w", err)
	}
	assets := make([]domain.Asset, len(rows))
	for i, row := range rows {
		assets[i] = toDomain(row)
	}
	return assets, nil
}

func (r *AssetRepo) Get(ctx context.Context, id, tenantID string) (domain.Asset, error) {
	row, err := r.queries.GetAsset(ctx, db.GetAssetParams{ID: id, TenantID: tenantID})
	if err != nil {
		return domain.Asset{}, fmt.Errorf("get asset: %w", err)
	}
	return toDomain(row), nil
}

func (r *AssetRepo) Update(ctx context.Context, id, tenantID, name, serialNumber string) (domain.Asset, error) {
	row, err := r.queries.UpdateAsset(ctx, db.UpdateAssetParams{
		ID:           id,
		TenantID:     tenantID,
		Name:         name,
		SerialNumber: sql.NullString{String: serialNumber, Valid: serialNumber != ""},
	})
	if err != nil {
		return domain.Asset{}, fmt.Errorf("update asset: %w", err)
	}
	return toDomain(row), nil
}

func (r *AssetRepo) Decommission(ctx context.Context, id, tenantID string) (domain.Asset, error) {
	row, err := r.queries.DecommissionAsset(ctx, db.DecommissionAssetParams{ID: id, TenantID: tenantID})
	if err != nil {
		return domain.Asset{}, fmt.Errorf("decommission asset: %w", err)
	}
	return toDomain(row), nil
}

func (r *AssetRepo) SetLocation(ctx context.Context, id, tenantID, facilityID, locationID string) (domain.Asset, error) {
	row, err := r.queries.SetAssetLocation(ctx, db.SetAssetLocationParams{
		ID:         id,
		TenantID:   tenantID,
		FacilityID: facilityID,
		LocationID: sql.NullString{String: locationID, Valid: locationID != ""},
	})
	if err != nil {
		return domain.Asset{}, fmt.Errorf("set asset location: %w", err)
	}
	return toDomain(row), nil
}

func toDomain(row db.Asset) domain.Asset {
	return domain.Asset{
		ID:           row.ID,
		TenantID:     row.TenantID,
		Name:         row.Name,
		AssetType:    row.AssetType,
		FacilityID:   row.FacilityID,
		LocationID:   row.LocationID.String,
		SerialNumber: row.SerialNumber.String,
		Status:       domain.AssetStatus(row.Status),
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
}
