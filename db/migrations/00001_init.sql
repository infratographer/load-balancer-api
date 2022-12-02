-- +goose Up
-- +goose StatementBegin

CREATE TABLE load_balancers (
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    load_balancer_id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    ip_addr inet NOT NULL,
    nice_name STRING NOT NULL,
    slug STRING NOT NULL,
    UNIQUE INDEX idx_load_balancer_ip_addr (ip_addr) WHERE deleted_at IS NULL,
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
    load_balancer_id UUID NOT NULL REFERENCES load_balancers (load_balancer_id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL,
    port INT NOT NULL,
    af_inet STRING NOT NULL DEFAULT 'ipv4',
    nice_name STRING NOT NULL,
    slug STRING NOT NULL,
    UNIQUE INDEX idx_frontends_ip_port (load_balancer_id, port) WHERE deleted_at is NULL,
    INDEX idx_frontends_tenant_id (tenant_id),
    INDEX idx_frontends_created_at (created_at),
    INDEX idx_frontends_updated_at (updated_at),
    INDEX idx_frontends_deleted_at (deleted_at)
);

CREATE TABLE frontend_attributes (
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    frontend_attr_id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    frontend_id UUID NOT NULL REFERENCES frontends(frontend_id) ON UPDATE CASCADE,
    frontend_attr JSONB NOT NULL,
    UNIQUE INDEX idx_frontend_attrib_frontend_id_key (frontend_id, frontend_attr),
    INDEX idx_frontend_attr_created_at (created_at),
    INDEX idx_frontend_attr_updated_at (updated_at),
    INDEX idx_frontend_attr_deleted_at (deleted_at)
);

CREATE TABLE assignments (
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    assignment_id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    pool_id UUID REFERENCES pools(pool_id) ON UPDATE CASCADE,
    frontend_id UUID REFERENCES frontends(frontend_id) ON UPDATE CASCADE,
    load_balancer_id UUID REFERENCES load_balancers(load_balancer_id) ON DELETE CASCADE,
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
    load_balancer_id UUID NOT NULL REFERENCES load_balancers (load_balancer_id) ON DELETE CASCADE,
    tenant_id UUID  NOT NULL,
    protocol STRING NOT NULL,
    has_details BOOL NOT NULL,
    use_proxy_protocol BOOL NOT NULL DEFAULT FALSE,
    nice_name STRING NOT NULL,
    slug STRING NOT NULL,
    INDEX idx_pools_tenant_id (tenant_id),
    INDEX idx_pools_created_at (created_at),
    INDEX idx_pools_updated_at (updated_at),
    INDEX idx_pools_deleted_at (deleted_at)
);

CREATE TABLE pool_attributes (
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    pool_attribute_id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    pool_id UUID NOT NULL REFERENCES pools(pool_id) ON UPDATE CASCADE,
    pool_attr STRING NOT NULL,
    pool_attr_value STRING NOT NULL,
    UNIQUE INDEX idx_pool_attrib_pool_id_key (pool_id, attribute, value),
    INDEX idx_pool_attrib_created_at (created_at),
    INDEX idx_pool_attrib_updated_at (updated_at),
    INDEX idx_pool_attrib_deleted_at (deleted_at)
);

CREATE TABLE origins (
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    origin_id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    pool_id UUID NOT NULL REFERENCES pools(pool_id) ON UPDATE CASCADE,
    origin_target STRING NOT NULL,
    port INT NOT NULL,
    has_details BOOL NOT NULL,
    tenant_id UUID NOT NULL REFERENCES pools(tenant_id) ON UPDATE CASCADE,
    nice_name STRING NOT NULL,
    slug STRING NOT NULL,
    use_proxy_protocol BOOL NOT NULL REFERENCES pools(use_proxy_protocol) ON UPDATE CASCADE,
    af_inet STRING NOT NULL REFERENCES pools(af_inet) ON UPDATE CASCADE,
    UNIQUE INDEX idx_origins_tenant_id_pool_id_ip_addr_port (pool_id, ip_addr, port),
    INDEX idx_origin_tenant_id (tenant_id),
    INDEX idx_origins_created_at (created_at),
    INDEX idx_origins_updated_at (updated_at),
    INDEX idx_origins_deleted_at (deleted_at)
);

CREATE TABLE origin_attributes (
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    origin_attributes_id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    origin_id UUID NOT NULL REFERENCES origins(origin_id) ON UPDATE CASCADE,
    attribute STRING,
    attr_value STRING,
    UNIQUE INDEX idx_origin_attrib_origin_id_attribute_value (origin_id, attribute, value),
    INDEX idx_origin_attrib_created_at (created_at),
    INDEX idx_origin_attrib_updated_at (updated_at),
    INDEX idx_origin_attrib_deleted_at (deleted_at)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE load_balancers;
DROP TABLE frontends;
DROP TABLE frontend_attributes;
DROP TABLE pools;
DROP TABLE pool_attributes;
DROP TABLE origins;
DROP TABLE origin_attributes;

-- +goose StatementEnd
