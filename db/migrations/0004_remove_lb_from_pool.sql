-- +goose Up
-- +goose StatementBegin
ALTER TABLE pools DROP COLUMN load_balancer_id CASCADE;
ALTER TABLE pools DROP COLUMN use_proxy_protocol CASCADE;
-- +goose StatementEnd