-- +goose Up
-- +goose StatementBegin

ALTER TABLE locations ADD CONSTRAINT uniq_locations_display_name UNIQUE (display_name);

-- +goose StatementEnd
