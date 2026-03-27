package http_test

// BDD scenarios: docs/features/set-asset-location.feature

import (
	"net/http"
	"testing"

	"github.com/Snip123/ai-test-assets-service/internal/domain"
)

func TestBDD_SetAssetLocation(t *testing.T) {
	const (
		tenantID = "dev-tenant"
		assetID  = "asset-001"
	)

	t.Run("Successfully set the Location of an Asset", func(t *testing.T) {
		repo := newMockRepo(seedAsset(assetID, "Rooftop HVAC Unit", tenantID, domain.AssetStatusActive))
		pub := &mockPublisher{}
		srv := newTestServer(repo, pub)
		defer srv.Close()

		body := `{"facility_id":"facility-001","location_id":"roof-level-3"}`
		resp := doRequest(http.MethodPut, srv.URL+"/v1/assets/"+assetID+"/location", tenantID, "FacilityManager", body)

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("want 200, got %d", resp.StatusCode)
		}
		respBody := decodeBody(resp)
		if respBody["facility_id"] != "facility-001" {
			t.Errorf("want facility_id 'facility-001', got %v", respBody["facility_id"])
		}
		if respBody["location_id"] != "roof-level-3" {
			t.Errorf("want location_id 'roof-level-3', got %v", respBody["location_id"])
		}
		if pub.published1() != "AssetLocationSet" {
			t.Errorf("want AssetLocationSet event, got %v", pub.published)
		}
	})

	t.Run("Asset Location can be updated to a new Location", func(t *testing.T) {
		asset := seedAsset(assetID, "Rooftop HVAC Unit", tenantID, domain.AssetStatusActive)
		asset.LocationID = "roof-level-3"
		repo := newMockRepo(asset)
		pub := &mockPublisher{}
		srv := newTestServer(repo, pub)
		defer srv.Close()

		body := `{"facility_id":"facility-001","location_id":"roof-level-4"}`
		resp := doRequest(http.MethodPut, srv.URL+"/v1/assets/"+assetID+"/location", tenantID, "FacilityManager", body)

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("want 200, got %d", resp.StatusCode)
		}
		respBody := decodeBody(resp)
		if respBody["location_id"] != "roof-level-4" {
			t.Errorf("want location_id 'roof-level-4', got %v", respBody["location_id"])
		}
		if pub.published1() != "AssetLocationSet" {
			t.Errorf("want AssetLocationSet event, got %v", pub.published)
		}
	})

	t.Run("Set Asset Location returns 404 when Asset does not exist", func(t *testing.T) {
		repo := newMockRepo()
		pub := &mockPublisher{}
		srv := newTestServer(repo, pub)
		defer srv.Close()

		body := `{"facility_id":"facility-001","location_id":"roof-level-3"}`
		resp := doRequest(http.MethodPut, srv.URL+"/v1/assets/asset-does-not-exist/location", tenantID, "FacilityManager", body)

		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("want 404, got %d", resp.StatusCode)
		}
	})

	t.Run("Technician cannot set Asset Location", func(t *testing.T) {
		repo := newMockRepo(seedAsset(assetID, "Rooftop HVAC Unit", tenantID, domain.AssetStatusActive))
		pub := &mockPublisher{}
		srv := newTestServer(repo, pub)
		defer srv.Close()

		body := `{"facility_id":"facility-001","location_id":"roof-level-3"}`
		resp := doRequest(http.MethodPut, srv.URL+"/v1/assets/"+assetID+"/location", tenantID, "Technician", body)

		if resp.StatusCode != http.StatusForbidden {
			t.Fatalf("want 403, got %d", resp.StatusCode)
		}
	})
}
