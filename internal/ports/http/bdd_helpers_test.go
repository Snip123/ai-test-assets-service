package http_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/Snip123/ai-test-assets-service/internal/application"
	"github.com/Snip123/ai-test-assets-service/internal/domain"
	"github.com/go-chi/chi/v5"

	httpports "github.com/Snip123/ai-test-assets-service/internal/ports/http"
)

// ── Mock repository ───────────────────────────────────────────────────────────

type mockRepo struct {
	assets map[string]domain.Asset // keyed by "id:tenantID"
}

func newMockRepo(initial ...domain.Asset) *mockRepo {
	r := &mockRepo{assets: make(map[string]domain.Asset)}
	for _, a := range initial {
		r.assets[repoKey(a.ID, a.TenantID)] = a
	}
	return r
}

func repoKey(id, tenantID string) string { return id + ":" + tenantID }

func (r *mockRepo) Create(_ context.Context, a domain.Asset) (domain.Asset, error) {
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	r.assets[repoKey(a.ID, a.TenantID)] = a
	return a, nil
}

func (r *mockRepo) List(_ context.Context, tenantID string) ([]domain.Asset, error) {
	var out []domain.Asset
	for _, a := range r.assets {
		if a.TenantID == tenantID {
			out = append(out, a)
		}
	}
	return out, nil
}

func (r *mockRepo) Get(_ context.Context, id, tenantID string) (domain.Asset, error) {
	a, ok := r.assets[repoKey(id, tenantID)]
	if !ok {
		return domain.Asset{}, fmt.Errorf("get asset: no rows in result set")
	}
	return a, nil
}

func (r *mockRepo) Update(_ context.Context, id, tenantID, name, serialNumber string) (domain.Asset, error) {
	a, ok := r.assets[repoKey(id, tenantID)]
	if !ok {
		return domain.Asset{}, fmt.Errorf("get asset: no rows in result set")
	}
	a.Name = name
	a.SerialNumber = serialNumber
	a.UpdatedAt = time.Now()
	r.assets[repoKey(id, tenantID)] = a
	return a, nil
}

func (r *mockRepo) Decommission(_ context.Context, id, tenantID string) (domain.Asset, error) {
	a, ok := r.assets[repoKey(id, tenantID)]
	if !ok {
		return domain.Asset{}, fmt.Errorf("get asset: no rows in result set")
	}
	a.Status = domain.AssetStatusDecommissioned
	a.UpdatedAt = time.Now()
	r.assets[repoKey(id, tenantID)] = a
	return a, nil
}

func (r *mockRepo) SetLocation(_ context.Context, id, tenantID, facilityID, locationID string) (domain.Asset, error) {
	a, ok := r.assets[repoKey(id, tenantID)]
	if !ok {
		return domain.Asset{}, fmt.Errorf("get asset: no rows in result set")
	}
	a.FacilityID = facilityID
	a.LocationID = locationID
	a.UpdatedAt = time.Now()
	r.assets[repoKey(id, tenantID)] = a
	return a, nil
}

// ── Mock event publisher ──────────────────────────────────────────────────────

type mockPublisher struct {
	published []string // event type names in order
}

func (p *mockPublisher) PublishAssetRegistered(_ context.Context, _ domain.AssetRegisteredEvent) error {
	p.published = append(p.published, "AssetRegistered")
	return nil
}

func (p *mockPublisher) PublishAssetAttributesUpdated(_ context.Context, _ domain.AssetAttributesUpdatedEvent) error {
	p.published = append(p.published, "AssetAttributesUpdated")
	return nil
}

func (p *mockPublisher) PublishAssetDecommissioned(_ context.Context, _ domain.AssetDecommissionedEvent) error {
	p.published = append(p.published, "AssetDecommissioned")
	return nil
}

func (p *mockPublisher) PublishAssetLocationSet(_ context.Context, _ domain.AssetLocationSetEvent) error {
	p.published = append(p.published, "AssetLocationSet")
	return nil
}

func (p *mockPublisher) published1() string {
	if len(p.published) == 0 {
		return ""
	}
	return p.published[0]
}

// ── Test server builder ───────────────────────────────────────────────────────

func newTestServer(repo *mockRepo, pub *mockPublisher) *httptest.Server {
	svc := application.NewAssetService(repo, pub)
	handler := httpports.NewAssetHandler(svc)

	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	return httptest.NewServer(r)
}

// ── Request helpers ───────────────────────────────────────────────────────────

func doRequest(method, url, tenantID, role, body string) *http.Response {
	var req *http.Request
	var err error

	if body != "" {
		req, err = http.NewRequest(method, url, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, http.NoBody)
	}
	if err != nil {
		panic(err)
	}

	req.Header.Set("X-Tenant-ID", tenantID)
	if role != "" {
		req.Header.Set("X-Platform-Role", role)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	return resp
}

func decodeBody(resp *http.Response) map[string]any {
	defer resp.Body.Close()
	var m map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&m)
	return m
}

// seedAsset is a convenience constructor for test fixtures.
func seedAsset(id, name, tenantID string, status domain.AssetStatus) domain.Asset {
	return domain.Asset{
		ID:         id,
		TenantID:   tenantID,
		Name:       name,
		AssetType:  "HVAC",
		FacilityID: "facility-001",
		Status:     status,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
