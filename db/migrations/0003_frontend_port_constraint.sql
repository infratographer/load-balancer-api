-- +goose Up
-- +goose StatementBegin
ALTER TABLE frontends ADD CONSTRAINT idx_frontends_port_tenant_id CHECK (port > 1 AND port < 65536);
-- +goose StatementEnd