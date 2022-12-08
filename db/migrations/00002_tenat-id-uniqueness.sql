-- +goose Up
-- +goose StatementBegin

ALTER TABLE locations ADD CONSTRAINT idx_locations_display_name_tenant_id UNIQUE (display_name, tenant_id);
DROP INDEX idx_load_balancer_ip_addr CASCADE;

ALTER TABLE load_balancers ADD CONSTRAINT idx_load_balancer_ip_addr_tenant_id UNIQUE (ip_addr, tenant_id) WHERE deleted_at IS NULL;

ALTER TABLE load_balancers ADD CONSTRAINT idx_load_balancer_ip_addr_public CHECK (NOT inet_contained_by_or_equals(ip_addr, '10.0.0.0/8'));

-- +goose StatementEnd
