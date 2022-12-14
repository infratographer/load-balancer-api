-- +goose Up
-- +goose StatementBegin

ALTER TABLE locations ADD CONSTRAINT idx_locations_display_name_tenant_id UNIQUE (display_name, tenant_id);

DROP INDEX idx_load_balancer_ip_addr CASCADE;

ALTER TABLE load_balancers ADD CONSTRAINT idx_load_balancer_ip_addr_tenant_id UNIQUE (ip_addr, tenant_id) WHERE deleted_at IS NULL;
ALTER TABLE load_balancers ADD CONSTRAINT idx_load_balancer_display_name_tenant_id UNIQUE (display_name, tenant_id) WHERE deleted_at IS NULL;

-- +goose StatementEnd
