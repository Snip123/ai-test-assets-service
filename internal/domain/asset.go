package domain

import "time"

// Asset is the central entity of this bounded context.
// It represents a physical piece of equipment owned or maintained by FSI.
type Asset struct {
	ID           string
	TenantID     string
	Name         string
	AssetType    string
	FacilityID   string
	SerialNumber string
	Status       AssetStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// AssetStatus represents the canonical lifecycle status of an Asset.
// Canonical Values — never display directly; resolve to Tenant Value in UI (ADR-0004).
type AssetStatus string

const (
	AssetStatusActive       AssetStatus = "Active"
	AssetStatusDecommission AssetStatus = "Decommissioned"
)

// AssetRegisteredEvent is the domain event published when a new Asset is registered.
// Published to NATS subject: fsi.{tenant-id}.assets.AssetRegistered (ADR-0016).
type AssetRegisteredEvent struct {
	AssetID      string
	TenantID     string
	AssetType    string
	FacilityID   string
	InstalledDate time.Time
}
