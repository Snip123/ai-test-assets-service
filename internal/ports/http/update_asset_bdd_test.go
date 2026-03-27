package http_test

// BDD scenarios: docs/features/update-asset.feature

import (
	"net/http"
	"testing"

	"github.com/Snip123/ai-test-assets-service/internal/domain"
)

func TestBDD_UpdateAsset(t *testing.T) {
	const (
		tenantID = "dev-tenant"
		assetID  = "asset-001"
	)

	t.Run("Successfully update Asset attributes", func(t *testing.T) {
		repo := newMockRepo(seedAsset(assetID, "Rooftop HVAC Unit", tenantID, domain.AssetStatusActive))
		pub := &mockPublisher{}
		srv := newTestServer(repo, pub)
		defer srv.Close()

		body := `{"name":"Rooftop HVAC Unit - R1","serial_number":"SN-2024-9999"}`
		resp := doRequest(http.MethodPatch, srv.URL+"/v1/assets/"+assetID, tenantID, "FacilityManager", body)

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("want 200, got %d", resp.StatusCode)
		}
		respBody := decodeBody(resp)
		if respBody["name"] != "Rooftop HVAC Unit - R1" {
			t.Errorf("want name 'Rooftop HVAC Unit - R1', got %v", respBody["name"])
		}
		if pub.published1() != "AssetAttributesUpdated" {
			t.Errorf("want AssetAttributesUpdated event, got %v", pub.published)
		}
	})

	t.Run("Update Asset fails when name is blank", func(t *testing.T) {
		repo := newMockRepo(seedAsset(assetID, "Rooftop HVAC Unit", tenantID, domain.AssetStatusActive))
		pub := &mockPublisher{}
		srv := newTestServer(repo, pub)
		defer srv.Close()

		body := `{"name":""}`
		resp := doRequest(http.MethodPatch, srv.URL+"/v1/assets/"+assetID, tenantID, "FacilityManager", body)

		if resp.StatusCode != http.StatusUnprocessableEntity {
			t.Fatalf("want 422, got %d", resp.StatusCode)
		}
	})

	t.Run("Update Asset returns 404 when Asset does not exist", func(t *testing.T) {
		repo := newMockRepo()
		pub := &mockPublisher{}
		srv := newTestServer(repo, pub)
		defer srv.Close()

		body := `{"name":"New Name"}`
		resp := doRequest(http.MethodPatch, srv.URL+"/v1/assets/asset-does-not-exist", tenantID, "FacilityManager", body)

		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("want 404, got %d", resp.StatusCode)
		}
	})

	t.Run("Technician cannot update an Asset", func(t *testing.T) {
		repo := newMockRepo(seedAsset(assetID, "Rooftop HVAC Unit", tenantID, domain.AssetStatusActive))
		pub := &mockPublisher{}
		srv := newTestServer(repo, pub)
		defer srv.Close()

		body := `{"name":"New Name"}`
		resp := doRequest(http.MethodPatch, srv.URL+"/v1/assets/"+assetID, tenantID, "Technician", body)

		if resp.StatusCode != http.StatusForbidden {
			t.Fatalf("want 403, got %d", resp.StatusCode)
		}
	})
}
