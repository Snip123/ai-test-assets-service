package http_test

// BDD scenarios: docs/features/get-asset.feature

import (
	"net/http"
	"testing"

	"github.com/Snip123/ai-test-assets-service/internal/domain"
)

func TestBDD_GetAsset(t *testing.T) {
	const (
		tenantID      = "dev-tenant"
		otherTenantID = "other-tenant"
		assetID       = "asset-001"
	)

	t.Run("Successfully get an Asset by ID", func(t *testing.T) {
		repo := newMockRepo(seedAsset(assetID, "Rooftop HVAC Unit", tenantID, domain.AssetStatusActive))
		pub := &mockPublisher{}
		srv := newTestServer(repo, pub)
		defer srv.Close()

		resp := doRequest(http.MethodGet, srv.URL+"/v1/assets/"+assetID, tenantID, "FacilityManager", "")

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("want 200, got %d", resp.StatusCode)
		}
		body := decodeBody(resp)
		if body["name"] != "Rooftop HVAC Unit" {
			t.Errorf("want name 'Rooftop HVAC Unit', got %v", body["name"])
		}
		if body["status"] != string(domain.AssetStatusActive) {
			t.Errorf("want status 'Active', got %v", body["status"])
		}
	})

	t.Run("Get Asset returns 404 when Asset does not exist", func(t *testing.T) {
		repo := newMockRepo()
		pub := &mockPublisher{}
		srv := newTestServer(repo, pub)
		defer srv.Close()

		resp := doRequest(http.MethodGet, srv.URL+"/v1/assets/asset-does-not-exist", tenantID, "FacilityManager", "")

		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("want 404, got %d", resp.StatusCode)
		}
	})

	t.Run("Get Asset returns 404 when Asset belongs to a different Tenant", func(t *testing.T) {
		// Asset exists only in other-tenant, not in dev-tenant
		repo := newMockRepo(seedAsset("asset-other", "Other Tenant Asset", otherTenantID, domain.AssetStatusActive))
		pub := &mockPublisher{}
		srv := newTestServer(repo, pub)
		defer srv.Close()

		resp := doRequest(http.MethodGet, srv.URL+"/v1/assets/asset-other", tenantID, "FacilityManager", "")

		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("want 404, got %d", resp.StatusCode)
		}
	})
}
