package postgres

import (
	"context"
	"fmt"

	"github.com/Snip123/ai-test-assets-service/internal/domain"
	db "github.com/Snip123/ai-test-assets-service/internal/generated/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AssetRepo struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewAssetRepo(pool *pgxpool.Pool) *AssetRepo {
	return &AssetRepo{pool: pool, queries: db.New(pool)}
}

func (r *AssetRepo) Create(ctx context.Context, a domain.Asset) (domain.Asset, error) {
	row, err := r.queries.CreateAsset(ctx, db.CreateAssetParams{
		ID:           a.ID,
		TenantID:     a.TenantID,
		Name:         a.Name,
		AssetType:    a.AssetType,
		FacilityID:   a.FacilityID,
		SerialNumber: pgText(a.SerialNumber),
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

func toDomain(row db.Asset) domain.Asset {
	return domain.Asset{
		ID:           row.ID,
		TenantID:     row.TenantID,
		Name:         row.Name,
		AssetType:    row.AssetType,
		FacilityID:   row.FacilityID,
		SerialNumber: row.SerialNumber.String,
		Status:       domain.AssetStatus(row.Status),
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
	}
}

func pgText(s string) interface{ Valid bool } {
	// pgx nullable text helper — returns pgtype.Text
	// Replaced by sqlc generated type in practice; stub for compilation
	type pgText struct {
		String string
		Valid  bool
	}
	return pgText{String: s, Valid: s != ""}
}
