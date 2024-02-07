-- +goose Up
-- modify "load_balancers" table
ALTER TABLE "load_balancers" ADD COLUMN "deleted_at" timestamptz NULL, ADD COLUMN "deleted_by" character varying NULL;

-- +goose Down
-- reverse: modify "load_balancers" table
ALTER TABLE "load_balancers" DROP COLUMN "deleted_by", DROP COLUMN "deleted_at";
