-- name: CreateAsset :one
INSERT INTO assets (id, tenant_id, name, asset_type, facility_id, serial_number, status)
VALUES (@id, @tenant_id, @name, @asset_type, @facility_id, @serial_number, @status)
RETURNING id, tenant_id, name, asset_type, facility_id, location_id, serial_number, status, created_at, updated_at;

-- name: ListAssets :many
SELECT id, tenant_id, name, asset_type, facility_id, location_id, serial_number, status, created_at, updated_at
FROM assets
WHERE tenant_id = @tenant_id
ORDER BY created_at DESC;

-- name: GetAsset :one
SELECT id, tenant_id, name, asset_type, facility_id, location_id, serial_number, status, created_at, updated_at
FROM assets
WHERE id = @id AND tenant_id = @tenant_id;

-- name: UpdateAsset :one
UPDATE assets
SET name = @name, serial_number = @serial_number, updated_at = NOW()
WHERE id = @id AND tenant_id = @tenant_id
RETURNING id, tenant_id, name, asset_type, facility_id, location_id, serial_number, status, created_at, updated_at;

-- name: DecommissionAsset :one
UPDATE assets
SET status = 'Decommissioned', updated_at = NOW()
WHERE id = @id AND tenant_id = @tenant_id
RETURNING id, tenant_id, name, asset_type, facility_id, location_id, serial_number, status, created_at, updated_at;

-- name: SetAssetLocation :one
UPDATE assets
SET facility_id = @facility_id, location_id = @location_id, updated_at = NOW()
WHERE id = @id AND tenant_id = @tenant_id
RETURNING id, tenant_id, name, asset_type, facility_id, location_id, serial_number, status, created_at, updated_at;
