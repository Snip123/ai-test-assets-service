CREATE TABLE assets (
    id            TEXT        NOT NULL,
    tenant_id     TEXT        NOT NULL,
    name          TEXT        NOT NULL,
    asset_type    TEXT        NOT NULL,
    facility_id   TEXT        NOT NULL,
    serial_number TEXT,
    status        TEXT        NOT NULL DEFAULT 'Active',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id, tenant_id)
);

CREATE INDEX idx_assets_tenant_id ON assets (tenant_id);
