-- +goose Up
-- +goose StatementBegin
ALTER TABLE load_balancers ADD COLUMN ip_address_id UUID NOT NULL;
ALTER TABLE load_balancers DROP COLUMN IF EXISTS ip_addr;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE load_balancers DROP COLUMN IF EXISTS ip_address_id;
ALTER TABLE load_balancers ADD COLUMN ip_addr inet NOT NULL;
-- +goose StatementEnd
