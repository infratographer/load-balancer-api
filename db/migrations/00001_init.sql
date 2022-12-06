-- +goose Up
-- +goose StatementBegin

CREATE TABLE load_balancers (
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    load_balancer_id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    ip_addr inet NOT NULL,
    display_name STRING NOT NULL,
    location_id UUID NOT NULL REFERENCES locations(location_id) ON UPDATE CASCADE,
    load_balancer_size STRING NOT NULL,
    load_balancer_type STRING NOT NULL,
    UNIQUE INDEX idx_load_balancer_ip_addr (ip_addr) WHERE deleted_at IS NULL,
    INDEX idx_load_balancer_tenant_id (tenant_id),
    INDEX idx_load_balancer_created_at (created_at),
    INDEX idx_load_balancer_updated_at (updated_at),
    INDEX idx_load_balancer_deleted_at (deleted_at)
);

CREATE TABLE locations (
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    location_id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    display_name STRING NOT NULL,
    INDEX idx_location_tenant_id (tenant_id),
    INDEX idx_location_created_at (created_at),
    INDEX idx_location_updated_at (updated_at),
    INDEX idx_location_deleted_at (deleted_at)
);

CREATE TABLE frontends (
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    frontend_id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    load_balancer_id UUID NOT NULL REFERENCES load_balancers (load_balancer_id) ON UPDATE CASCADE,
    tenant_id UUID NOT NULL,
    port INT NOT NULL,
    af_inet STRING NOT NULL DEFAULT 'ipv4',
    display_name STRING NOT NULL,
    UNIQUE INDEX idx_frontends_ip_port (load_balancer_id, port) WHERE deleted_at is NULL,
    INDEX idx_frontends_tenant_id (tenant_id),
    INDEX idx_frontends_created_at (created_at),
    INDEX idx_frontends_updated_at (updated_at),
    INDEX idx_frontends_deleted_at (deleted_at)
);


CREATE TABLE assignments (
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    assignment_id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    pool_id UUID REFERENCES pools(pool_id) ON UPDATE CASCADE,
    frontend_id UUID REFERENCES frontends(frontend_id) ON UPDATE CASCADE,
    load_balancer_id UUID REFERENCES load_balancers(load_balancer_id) ON UPDATE CASCADE,
    tenant_id UUID NOT NULL,
    UNIQUE INDEX idx_pool_frontend_assocs_pool_id_frontend_id (pool_id, frontend_id) WHERE deleted_at is NULL,
    INDEX idx_assocs_tenant_id (tenant_id),
    INDEX idx_assocs_created_at (created_at),
    INDEX idx_assocs_updated_at (updated_at),
    INDEX idx_assocs_deleted_at (deleted_at)
)

CREATE TABLE pools (
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    pool_id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    load_balancer_id UUID NOT NULL REFERENCES load_balancers (load_balancer_id) ON UPDATE CASCADE,
    tenant_id UUID  NOT NULL,
    protocol STRING NOT NULL,
    use_proxy_protocol BOOL NOT NULL DEFAULT FALSE,
    display_name STRING NOT NULL,
    INDEX idx_pools_tenant_id (tenant_id),
    INDEX idx_pools_created_at (created_at),
    INDEX idx_pools_updated_at (updated_at),
    INDEX idx_pools_deleted_at (deleted_at)
);


CREATE TABLE origins (
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    origin_id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    pool_id UUID NOT NULL REFERENCES pools(pool_id) ON UPDATE CASCADE,
    origin_target INET NOT NULL,
    port INT NOT NULL,
    tenant_id UUID NOT NULL REFERENCES pools(tenant_id) ON UPDATE CASCADE,
    display_name STRING NOT NULL,
    origin_disabled BOOL NOT NULL DEFAULT TRUE,
    UNIQUE INDEX idx_origins_tenant_id_pool_id_ip_addr_port (pool_id, ip_addr, port),
    INDEX idx_origin_tenant_id (tenant_id),
    INDEX idx_origins_created_at (created_at),
    INDEX idx_origins_updated_at (updated_at),
    INDEX idx_origins_deleted_at (deleted_at)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE load_balancers;
DROP TABLE frontends;
DROP TABLE pools;
DROP TABLE origins;
DROP TABLE assignments;

-- +goose StatementEnd
