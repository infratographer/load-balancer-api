-- +goose Up
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_frontends_ip_port CASCADE;
DROP INDEX IF EXISTS idx_frontends_created_at CASCADE;
DROP INDEX IF EXISTS idx_frontends_updated_at CASCADE;
DROP INDEX IF EXISTS idx_frontends_deleted_at CASCADE;

ALTER TABLE frontends RENAME TO ports;
ALTER TABLE ports RENAME COLUMN frontend_id to port_id;
ALTER TABLE ports RENAME CONSTRAINT frontends_pkey TO ports_pkey;

ALTER TABLE ports ADD CONSTRAINT IF NOT EXISTS idx_ports_ip_port UNIQUE (load_balancer_id, port) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_ports_created_at on ports (created_at);
CREATE INDEX IF NOT EXISTS idx_ports_updated_at on ports (updated_at);
CREATE INDEX IF NOT EXISTS idx_ports_deleted_at on ports (deleted_at);

ALTER TABLE assignments ADD COLUMN port_id UUID REFERENCES ports(port_id) ON UPDATE CASCADE NOT NULL;
ALTER TABLE assignments DROP COLUMN frontend_id;

DROP INDEX IF EXISTS idx_pool_frontend_assocs_pool_id_frontend_id CASCADE;
ALTER TABLE assignments ADD CONSTRAINT IF NOT EXISTS idx_pool_port_assocs_pool_id_port_id UNIQUE (pool_id, port_id) WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_ports_ip_port CASCADE;
DROP INDEX IF EXISTS idx_ports_created_at CASCADE;
DROP INDEX IF EXISTS idx_ports_updated_at CASCADE;
DROP INDEX IF EXISTS idx_ports_deleted_at CASCADE;

ALTER TABLE ports RENAME TO frontends;
ALTER TABLE frontends RENAME COLUMN port_id TO frontend_id;
ALTER TABLE frontends RENAME CONSTRAINT ports_pkey TO frontends_pkey;

ALTER TABLE ports ADD CONSTRAINT IF NOT EXISTS idx_frontends_ip_port UNIQUE (load_balancer_id, port) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_frontends_created_at on ports (created_at);
CREATE INDEX IF NOT EXISTS idx_frontends_updated_at on ports (updated_at);
CREATE INDEX IF NOT EXISTS idx_frontends_deleted_at on ports (deleted_at);

ALTER TABLE assignments ADD COLUMN frontend_id UUID REFERENCES frontends(frontend_id) ON UPDATE CASCADE NOT NULL;
ALTER TABLE assignments DROP COLUMN port_id;

DROP INDEX IF EXISTS idx_pool_port_assocs_pool_id_port_id CASCADE;
ALTER TABLE assignments ADD CONSTRAINT IF NOT EXISTS idx_pool_frontend_assocs_pool_id_frontend_id UNIQUE (pool_id, frontend_id) WHERE deleted_at IS NULL;

-- +goose StatementEnd
