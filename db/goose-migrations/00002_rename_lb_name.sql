-- +goose Up
-- +goose StatementBegin
ALTER TABLE load_balancers RENAME COLUMN "display_name" to "name";
ALTER TABLE frontends RENAME COLUMN "display_name" to "name";
ALTER TABLE origins RENAME COLUMN "display_name" to "name";
ALTER TABLE pools RENAME COLUMN "display_name" to "name";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE pools RENAME COLUMN "name" to "display_name";
ALTER TABLE origins RENAME COLUMN "name" to "display_name";
ALTER TABLE frontends RENAME COLUMN "name" to "display_name";
ALTER TABLE load_balancers RENAME COLUMN "name" to "display_name";
-- +goose StatementEnd
