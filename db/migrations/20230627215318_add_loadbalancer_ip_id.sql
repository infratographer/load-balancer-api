-- +goose Up
-- add "ip_id" column to "load_balancers" table
ALTER TABLE "load_balancers" ADD COLUMN "ip_id" character varying;

-- +goose Down
-- revert adding "ip_id" column to "load_balancers" table
ALTER TABLE "load_balancers" DROP COLUMN "ip_id";


