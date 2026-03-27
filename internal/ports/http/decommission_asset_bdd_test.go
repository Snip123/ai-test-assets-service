package http_test

// BDD scenarios: docs/features/decommission-asset.feature

import (
	"net/http"
	"testing"

	"github.com/Snip123/ai-test-assets-service/internal/domain"
)

func TestBDD_DecommissionAsset(t *testing.T) {
	const (
		tenantID = "dev-tenant"
		assetID  = "asset-001"
	)

	t.Run("Successfully decommission an Active Asset", func(t *testing.T) {
		repo := newMockRepo(seedAsset(assetID, "Rooftop HVAC Unit", tenantID, domain.AssetStatusActive))
		pub := &mockPublisher{}
		srv := newTestServer(repo, pub)
		defer srv.Close()

		body := `{"reason":"End of service life"}`
		resp := doRequest(http.MethodPost, srv.URL+"/v1/assets/"+assetID+"/decommission", tenantID, "FacilityManager", body)

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("want 200, got %d", resp.StatusCode)
		}
		respBody := decodeBody(resp)
		if respBody["status"] != string(domain.AssetStatusDecommissioned) {
			t.Errorf("want status 'Decommissioned', got %v", respBody["status"])
		}
		if pub.published1() != "AssetDecommissioned" {
			t.Errorf("want AssetDecommissioned event, got %v", pub.published)
		}
	})

	t.Run("Cannot decommission an already-Decommissioned Asset", func(t *testing.T) {
		repo := newMockRepo(seedAsset(assetID, "Rooftop HVAC Unit", tenantID, domain.AssetStatusDecommissioned))
		pub := &mockPublisher{}
		srv := newTestServer(repo, pub)
		defer srv.Close()

		body := `{"reason":"Duplicate decommission"}`
		resp := doRequest(http.MethodPost, srv.URL+"/v1/assets/"+assetID+"/decommission", tenantID, "FacilityManager", body)

		if resp.StatusCode != http.StatusConflict {
			t.Fatalf("want 409, got %d", resp.StatusCode)
		}
	})

	t.Run("Decommission returns 404 when Asset does not exist", func(t *testing.T) {
		repo := newMockRepo()
		pub := &mockPublisher{}
		srv := newTestServer(repo, pub)
		defer srv.Close()

		body := `{"reason":"Not found"}`
		resp := doRequest(http.MethodPost, srv.URL+"/v1/assets/asset-does-not-exist/decommission", tenantID, "FacilityManager", body)

		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("want 404, got %d", resp.StatusCode)
		}
	})

	t.Run("Technician cannot decommission an Asset", func(t *testing.T) {
		repo := newMockRepo(seedAsset(assetID, "Rooftop HVAC Unit", tenantID, domain.AssetStatusActive))
		pub := &mockPublisher{}
		srv := newTestServer(repo, pub)
		defer srv.Close()

		body := `{"reason":"Technician attempt"}`
		resp := doRequest(http.MethodPost, srv.URL+"/v1/assets/"+assetID+"/decommission", tenantID, "Technician", body)

		if resp.StatusCode != http.StatusForbidden {
			t.Fatalf("want 403, got %d", resp.StatusCode)
		}
	})
}
