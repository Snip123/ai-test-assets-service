-- name: CreateAsset :one
INSERT INTO assets (id, tenant_id, name, asset_type, facility_id, serial_number, status)
VALUES (@id, @tenant_id, @name, @asset_type, @facility_id, @serial_number, @status)
RETURNING *;

-- name: ListAssets :many
SELECT * FROM assets
WHERE tenant_id = @tenant_id
ORDER BY created_at DESC;

-- name: GetAsset :one
SELECT * FROM assets
WHERE id = @id AND tenant_id = @tenant_id;
