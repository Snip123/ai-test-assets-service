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
	LocationID   string
	SerialNumber string
	Status       AssetStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// AssetStatus represents the canonical lifecycle status of an Asset.
// Canonical Values — never display directly; resolve to Tenant Value in UI (ADR-0004).
type AssetStatus string

const (
	AssetStatusActive        AssetStatus = "Active"
	AssetStatusDecommissioned AssetStatus = "Decommissioned"
)

// AssetRegisteredEvent is published when a new Asset is registered.
// Subject: fsi.{tenant-id}.assets.AssetRegistered (ADR-0016).
type AssetRegisteredEvent struct {
	AssetID       string
	TenantID      string
	AssetType     string
	FacilityID    string
	InstalledDate time.Time
}

// AssetAttributesUpdatedEvent is published when mutable Asset attributes change.
// Subject: fsi.{tenant-id}.assets.AssetAttributesUpdated (ADR-0016).
type AssetAttributesUpdatedEvent struct {
	AssetID           string
	TenantID          string
	ChangedAttributes map[string]any
}

// AssetDecommissionedEvent is published when an Asset is decommissioned.
// Subject: fsi.{tenant-id}.assets.AssetDecommissioned (ADR-0016).
type AssetDecommissionedEvent struct {
	AssetID            string
	TenantID           string
	DecommissionedDate time.Time
	Reason             string
}

// AssetLocationSetEvent is published when an Asset's Location is set or updated.
// Subject: fsi.{tenant-id}.assets.AssetLocationSet (ADR-0016).
type AssetLocationSetEvent struct {
	AssetID    string
	TenantID   string
	FacilityID string
	LocationID string
}
