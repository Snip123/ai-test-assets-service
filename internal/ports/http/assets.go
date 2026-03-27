package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Snip123/ai-test-assets-service/internal/application"
	"github.com/Snip123/ai-test-assets-service/internal/domain"
	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var tracer = otel.Tracer("ai-test-assets-service/http")

// AssetHandler wires HTTP routes to the AssetService application layer.
type AssetHandler struct {
	svc *application.AssetService
}

func NewAssetHandler(svc *application.AssetService) *AssetHandler {
	return &AssetHandler{svc: svc}
}

// RegisterRoutes mounts asset routes on the provided router.
func (h *AssetHandler) RegisterRoutes(r chi.Router) {
	r.Get("/v1/assets", h.listAssets)
	r.Post("/v1/assets", h.registerAsset)
	r.Get("/v1/assets/{id}", h.getAsset)
	r.Patch("/v1/assets/{id}", h.updateAsset)
	r.Post("/v1/assets/{id}/decommission", h.decommissionAsset)
	r.Put("/v1/assets/{id}/location", h.setAssetLocation)
}

func (h *AssetHandler) listAssets(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "AssetHandler.listAssets")
	defer span.End()

	tenantID := r.Header.Get("X-Tenant-ID")
	span.SetAttributes(attribute.String("tenant_id", tenantID))

	assets, err := h.svc.ListAssets(ctx, tenantID)
	if err != nil {
		writeProblem(w, r, http.StatusInternalServerError, "list-assets-failed", "Failed to list assets", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data": toAPIAssets(assets),
		"pagination": map[string]any{
			"has_more":    false,
			"total_count": len(assets),
		},
	})
}

func (h *AssetHandler) registerAsset(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "AssetHandler.registerAsset")
	defer span.End()

	tenantID := r.Header.Get("X-Tenant-ID")
	role := r.Header.Get("X-Platform-Role")
	span.SetAttributes(
		attribute.String("tenant_id", tenantID),
		attribute.String("platform_role", role),
	)

	if !canManageAssets(role) {
		writeProblem(w, r, http.StatusForbidden, "insufficient-role",
			"Insufficient Platform Role",
			fmt.Sprintf("Platform Role %q cannot register Assets", role))
		return
	}

	var req struct {
		Name         string `json:"name"`
		AssetType    string `json:"asset_type"`
		FacilityID   string `json:"facility_id"`
		SerialNumber string `json:"serial_number"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeProblem(w, r, http.StatusUnprocessableEntity, "invalid-body", "Invalid request body", err.Error())
		return
	}

	asset, err := h.svc.RegisterAsset(ctx, application.RegisterAssetCommand{
		TenantID:     tenantID,
		Name:         req.Name,
		AssetType:    req.AssetType,
		FacilityID:   req.FacilityID,
		SerialNumber: req.SerialNumber,
	})
	if err != nil {
		writeProblem(w, r, http.StatusUnprocessableEntity, "validation-error", "Validation failed", err.Error())
		return
	}

	span.SetAttributes(attribute.String("asset_id", asset.ID))
	w.Header().Set("Location", fmt.Sprintf("/v1/assets/%s", asset.ID))
	writeJSON(w, http.StatusCreated, toAPIAsset(asset))
}

func (h *AssetHandler) getAsset(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "AssetHandler.getAsset")
	defer span.End()

	tenantID := r.Header.Get("X-Tenant-ID")
	id := chi.URLParam(r, "id")
	span.SetAttributes(
		attribute.String("tenant_id", tenantID),
		attribute.String("asset_id", id),
	)

	asset, err := h.svc.GetAsset(ctx, id, tenantID)
	if err != nil {
		if err == application.ErrAssetNotFound {
			writeProblem(w, r, http.StatusNotFound, "asset-not-found", "Asset Not Found",
				fmt.Sprintf("Asset %q does not exist in this tenant", id))
			return
		}
		writeProblem(w, r, http.StatusInternalServerError, "get-asset-failed", "Failed to get asset", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, toAPIAsset(asset))
}

func (h *AssetHandler) updateAsset(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "AssetHandler.updateAsset")
	defer span.End()

	tenantID := r.Header.Get("X-Tenant-ID")
	role := r.Header.Get("X-Platform-Role")
	id := chi.URLParam(r, "id")
	span.SetAttributes(
		attribute.String("tenant_id", tenantID),
		attribute.String("asset_id", id),
	)

	if !canManageAssets(role) {
		writeProblem(w, r, http.StatusForbidden, "insufficient-role",
			"Insufficient Platform Role",
			fmt.Sprintf("Platform Role %q cannot update Assets", role))
		return
	}

	var req struct {
		Name         string `json:"name"`
		SerialNumber string `json:"serial_number"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeProblem(w, r, http.StatusUnprocessableEntity, "invalid-body", "Invalid request body", err.Error())
		return
	}

	asset, err := h.svc.UpdateAsset(ctx, application.UpdateAssetCommand{
		TenantID:     tenantID,
		ID:           id,
		Name:         req.Name,
		SerialNumber: req.SerialNumber,
	})
	if err != nil {
		switch err {
		case application.ErrAssetNotFound:
			writeProblem(w, r, http.StatusNotFound, "asset-not-found", "Asset Not Found",
				fmt.Sprintf("Asset %q does not exist in this tenant", id))
		default:
			writeProblem(w, r, http.StatusUnprocessableEntity, "validation-error", "Validation failed", err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, toAPIAsset(asset))
}

func (h *AssetHandler) decommissionAsset(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "AssetHandler.decommissionAsset")
	defer span.End()

	tenantID := r.Header.Get("X-Tenant-ID")
	role := r.Header.Get("X-Platform-Role")
	id := chi.URLParam(r, "id")
	span.SetAttributes(
		attribute.String("tenant_id", tenantID),
		attribute.String("asset_id", id),
	)

	if !canManageAssets(role) {
		writeProblem(w, r, http.StatusForbidden, "insufficient-role",
			"Insufficient Platform Role",
			fmt.Sprintf("Platform Role %q cannot decommission Assets", role))
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeProblem(w, r, http.StatusUnprocessableEntity, "invalid-body", "Invalid request body", err.Error())
		return
	}

	asset, err := h.svc.DecommissionAsset(ctx, application.DecommissionAssetCommand{
		TenantID: tenantID,
		ID:       id,
		Reason:   req.Reason,
	})
	if err != nil {
		switch err {
		case application.ErrAssetNotFound:
			writeProblem(w, r, http.StatusNotFound, "asset-not-found", "Asset Not Found",
				fmt.Sprintf("Asset %q does not exist in this tenant", id))
		case application.ErrAssetAlreadyDecommissioned:
			writeProblem(w, r, http.StatusConflict, "asset-already-decommissioned",
				"Asset Already Decommissioned",
				fmt.Sprintf("Asset %q is already decommissioned", id))
		default:
			writeProblem(w, r, http.StatusInternalServerError, "decommission-failed", "Failed to decommission asset", err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, toAPIAsset(asset))
}

func (h *AssetHandler) setAssetLocation(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "AssetHandler.setAssetLocation")
	defer span.End()

	tenantID := r.Header.Get("X-Tenant-ID")
	role := r.Header.Get("X-Platform-Role")
	id := chi.URLParam(r, "id")
	span.SetAttributes(
		attribute.String("tenant_id", tenantID),
		attribute.String("asset_id", id),
	)

	if !canManageAssets(role) {
		writeProblem(w, r, http.StatusForbidden, "insufficient-role",
			"Insufficient Platform Role",
			fmt.Sprintf("Platform Role %q cannot set Asset Location", role))
		return
	}

	var req struct {
		FacilityID string `json:"facility_id"`
		LocationID string `json:"location_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeProblem(w, r, http.StatusUnprocessableEntity, "invalid-body", "Invalid request body", err.Error())
		return
	}

	asset, err := h.svc.SetAssetLocation(ctx, application.SetAssetLocationCommand{
		TenantID:   tenantID,
		ID:         id,
		FacilityID: req.FacilityID,
		LocationID: req.LocationID,
	})
	if err != nil {
		switch err {
		case application.ErrAssetNotFound:
			writeProblem(w, r, http.StatusNotFound, "asset-not-found", "Asset Not Found",
				fmt.Sprintf("Asset %q does not exist in this tenant", id))
		default:
			writeProblem(w, r, http.StatusUnprocessableEntity, "validation-error", "Validation failed", err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, toAPIAsset(asset))
}

// ── helpers ───────────────────────────────────────────────────────────────────

// canManageAssets returns true for Platform Roles permitted to create/update/decommission Assets.
func canManageAssets(role string) bool {
	switch role {
	case "FacilityManager", "TenantAdmin", "FSICustomerSupport", "FSIPlatformAdmin":
		return true
	}
	return false
}

type apiAsset struct {
	ID           string         `json:"id"`
	TenantID     string         `json:"tenant_id"`
	Name         string         `json:"name"`
	AssetType    string         `json:"asset_type"`
	FacilityID   string         `json:"facility_id"`
	LocationID   string         `json:"location_id,omitempty"`
	SerialNumber string         `json:"serial_number,omitempty"`
	Status       string         `json:"status"`
	Display      map[string]any `json:"display"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

func toAPIAsset(a domain.Asset) apiAsset {
	return apiAsset{
		ID:           a.ID,
		TenantID:     a.TenantID,
		Name:         a.Name,
		AssetType:    string(a.AssetType),
		FacilityID:   a.FacilityID,
		LocationID:   a.LocationID,
		SerialNumber: a.SerialNumber,
		Status:       string(a.Status),
		Display: map[string]any{
			// Tenant Value resolution — placeholder until TenantAdmin config API is built
			"asset_type": string(a.AssetType),
			"status":     string(a.Status),
		},
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

func toAPIAssets(assets []domain.Asset) []apiAsset {
	result := make([]apiAsset, len(assets))
	for i, a := range assets {
		result[i] = toAPIAsset(a)
	}
	return result
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeProblem(w http.ResponseWriter, r *http.Request, status int, errType, title, detail string) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"type":     fmt.Sprintf("https://fsi-platform.com/errors/%s", errType),
		"title":    title,
		"status":   status,
		"detail":   detail,
		"instance": r.URL.Path,
	})
}
