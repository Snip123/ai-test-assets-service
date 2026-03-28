package http_test

// BDD scenarios: docs/features/register-asset.feature

import (
	"net/http"
	"strings"
	"testing"

	"github.com/Snip123/ai-test-assets-service/internal/domain"
)

func TestBDD_RegisterAsset(t *testing.T) {
	const (
		tenantID      = "dev-tenant"
		otherTenantID = "other-tenant"
	)

	t.Run("Successfully register an Asset", func(t *testing.T) {
		repo := newMockRepo()
		pub := &mockPublisher{}
		srv := newTestServer(repo, pub)
		defer srv.Close()

		body := `{"name":"Rooftop HVAC Unit","asset_type":"HVAC","facility_id":"facility-001","serial_number":"SN-2024-0042"}`
		resp := doRequest(http.MethodPost, srv.URL+"/v1/assets", tenantID, "FacilityManager", body)

		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("want 201, got %d", resp.StatusCode)
		}
		respBody := decodeBody(resp)
		if respBody["id"] == "" || respBody["id"] == nil {
			t.Errorf("want non-empty asset ID in response, got %v", respBody["id"])
		}
		if respBody["status"] != string(domain.AssetStatusActive) {
			t.Errorf("want status 'Active', got %v", respBody["status"])
		}
		if pub.published1() != "AssetRegistered" {
			t.Errorf("want AssetRegistered event, got %v", pub.published)
		}
	})

	t.Run("List Assets returns only Assets belonging to the authenticated Tenant", func(t *testing.T) {
		repo := newMockRepo(
			seedAsset("asset-dev", "Rooftop HVAC Unit", tenantID, domain.AssetStatusActive),
			seedAsset("asset-other", "Boiler Room Unit", otherTenantID, domain.AssetStatusActive),
		)
		pub := &mockPublisher{}
		srv := newTestServer(repo, pub)
		defer srv.Close()

		resp := doRequest(http.MethodGet, srv.URL+"/v1/assets", tenantID, "FacilityManager", "")

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("want 200, got %d", resp.StatusCode)
		}
		respBody := decodeBody(resp)
		data, _ := respBody["data"].([]any)

		var names []string
		for _, item := range data {
			if a, ok := item.(map[string]any); ok {
				if n, ok := a["name"].(string); ok {
					names = append(names, n)
				}
			}
		}

		found := false
		for _, n := range names {
			if n == "Rooftop HVAC Unit" {
				found = true
			}
			if n == "Boiler Room Unit" {
				t.Errorf("response must not contain asset from other-tenant, but found 'Boiler Room Unit'")
			}
		}
		if !found {
			t.Errorf("want 'Rooftop HVAC Unit' in response, got %v", names)
		}
	})

	t.Run("Register Asset fails when required fields are missing", func(t *testing.T) {
		repo := newMockRepo()
		pub := &mockPublisher{}
		srv := newTestServer(repo, pub)
		defer srv.Close()

		body := `{"name":"","asset_type":"HVAC"}`
		resp := doRequest(http.MethodPost, srv.URL+"/v1/assets", tenantID, "FacilityManager", body)

		if resp.StatusCode != http.StatusUnprocessableEntity {
			t.Fatalf("want 422, got %d", resp.StatusCode)
		}
		respBody := decodeBody(resp)
		detail, _ := respBody["detail"].(string)
		if !strings.Contains(detail, "name") {
			t.Errorf("want validation error detail to reference field 'name', got %q", detail)
		}
	})

	t.Run("Technician cannot register an Asset", func(t *testing.T) {
		repo := newMockRepo()
		pub := &mockPublisher{}
		srv := newTestServer(repo, pub)
		defer srv.Close()

		body := `{"name":"Rooftop HVAC Unit","asset_type":"HVAC","facility_id":"facility-001"}`
		resp := doRequest(http.MethodPost, srv.URL+"/v1/assets", tenantID, "Technician", body)

		if resp.StatusCode != http.StatusForbidden {
			t.Fatalf("want 403, got %d", resp.StatusCode)
		}
	})
}
