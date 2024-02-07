-- +goose Up
-- modify "load_balancers" table
ALTER TABLE "load_balancers" ADD COLUMN "created_by" character varying NULL, ADD COLUMN "updated_by" character varying NULL;

-- +goose Down
-- reverse: modify "load_balancers" table
ALTER TABLE "load_balancers" DROP COLUMN "updated_by", DROP COLUMN "created_by";
