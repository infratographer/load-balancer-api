-- +goose Up
-- +goose StatementBegin

CREATE TABLE load_balancers (
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    load_balancer_id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    location_id UUID NOT NULL,
    tenant_id UUID NOT NULL,
    ip_addr inet NOT NULL,
    display_name STRING NOT NULL,
    slug STRING NOT NULL,
    load_balancer_size STRING NOT NULL,
    load_balancer_type STRING NOT NULL,
    state_changed_at TIMESTAMPTZ,
    current_state STRING NOT NULL,
    previous_state STRING NOT NULL,
    UNIQUE INDEX idx_load_balancers_tenant_id_slug (tenant_id, slug) WHERE deleted_at IS NULL,
    UNIQUE INDEX idx_load_balancer_display_name_tenant_id (tenant_id,ip_addr) WHERE deleted_at IS NULL,
    INDEX idx_load_balancer_tenant_id (tenant_id),
    INDEX idx_load_balancer_created_at (created_at),
    INDEX idx_load_balancer_updated_at (updated_at),
    INDEX idx_load_balancer_deleted_at (deleted_at)
);


CREATE TABLE frontends (
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    frontend_id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    load_balancer_id UUID NOT NULL REFERENCES load_balancers (load_balancer_id) ON UPDATE CASCADE,
    tenant_id UUID NOT NULL,
    port INT NOT NULL CHECK (port > 1 AND port < 65536),
    af_inet STRING NOT NULL DEFAULT 'ipv4',
    display_name STRING NOT NULL,
    slug STRING NOT NULL,
    state_changed_at TIMESTAMPTZ,
    current_state STRING NOT NULL,
    previous_state STRING NOT NULL,
    UNIQUE INDEX idx_frontends_tenant_id_slug  (tenant_id, slug) WHERE deleted_at IS NULL,
    UNIQUE INDEX idx_frontends_ip_port (load_balancer_id, port) WHERE deleted_at is NULL,
    INDEX idx_frontends_tenant_id (tenant_id),
    INDEX idx_frontends_created_at (created_at),
    INDEX idx_frontends_updated_at (updated_at),
    INDEX idx_frontends_deleted_at (deleted_at)
);



CREATE TABLE pools (
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    pool_id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    protocol STRING NOT NULL,
    display_name STRING NOT NULL,
    slug STRING NOT NULL,
    tenant_id UUID  NOT NULL,
    UNIQUE INDEX idx_pools_tenant_id_slug (tenant_id, slug) WHERE deleted_at IS NULL,
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
    tenant_id UUID NOT NULL,
    display_name STRING NOT NULL,
    slug STRING NOT NULL,
    origin_user_setting_disabled BOOL NOT NULL DEFAULT TRUE,
    state_changed_at TIMESTAMPTZ,
    current_state STRING NOT NULL,
    previous_state STRING NOT NULL,
    UNIQUE INDEX idx_origins_tenant_id_slug (tenant_id, slug) WHERE deleted_at IS NULL,
    UNIQUE INDEX idx_origins_pool_id_origin_target_port (pool_id, origin_target, port) WHERE deleted_at is NULL,
    INDEX idx_origin_tenant_id (tenant_id),
    INDEX idx_origins_created_at (created_at),
    INDEX idx_origins_updated_at (updated_at),
    INDEX idx_origins_deleted_at (deleted_at)
);


CREATE TABLE assignments (
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    assignment_id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    pool_id UUID REFERENCES pools(pool_id) ON UPDATE CASCADE NOT NULL,
    frontend_id UUID REFERENCES frontends(frontend_id) ON UPDATE CASCADE NOT NULL,
    load_balancer_id UUID REFERENCES load_balancers(load_balancer_id) ON UPDATE CASCADE NOT NULL,
    tenant_id UUID NOT NULL,
    UNIQUE INDEX idx_pool_frontend_assocs_pool_id_frontend_id (pool_id, frontend_id) WHERE deleted_at is NULL,
    INDEX idx_assocs_tenant_id (tenant_id),
    INDEX idx_assocs_created_at (created_at),
    INDEX idx_assocs_updated_at (updated_at),
    INDEX idx_assocs_deleted_at (deleted_at)
);



-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE assignments;
DROP TABLE origins;
DROP TABLE pools;
DROP TABLE frontends;
DROP TABLE load_balancers;

-- +goose StatementEnd
