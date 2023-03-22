-- +goose Up
-- +goose StatementBegin
ALTER TABLE load_balancers RENAME COLUMN "display_name" to "name";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE load_balancers RENAME COLUMN "name" to "display_name";
-- +goose StatementEnd
