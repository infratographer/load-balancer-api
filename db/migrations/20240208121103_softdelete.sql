-- +goose Up
-- modify "load_balancers" table
ALTER TABLE "load_balancers" ADD COLUMN "deleted_at" timestamptz NULL, ADD COLUMN "deleted_by" character varying NULL;
-- modify "origins" table
ALTER TABLE "origins" ADD COLUMN "deleted_at" timestamptz NULL, ADD COLUMN "deleted_by" character varying NULL;
-- modify "pools" table
ALTER TABLE "pools" ADD COLUMN "deleted_at" timestamptz NULL, ADD COLUMN "deleted_by" character varying NULL;
-- modify "ports" table
ALTER TABLE "ports" ADD COLUMN "deleted_at" timestamptz NULL, ADD COLUMN "deleted_by" character varying NULL;
-- modify "providers" table
ALTER TABLE "providers" ADD COLUMN "deleted_at" timestamptz NULL, ADD COLUMN "deleted_by" character varying NULL;

-- +goose Down
-- reverse: modify "providers" table
ALTER TABLE "providers" DROP COLUMN "deleted_by", DROP COLUMN "deleted_at";
-- reverse: modify "ports" table
ALTER TABLE "ports" DROP COLUMN "deleted_by", DROP COLUMN "deleted_at";
-- reverse: modify "pools" table
ALTER TABLE "pools" DROP COLUMN "deleted_by", DROP COLUMN "deleted_at";
-- reverse: modify "origins" table
ALTER TABLE "origins" DROP COLUMN "deleted_by", DROP COLUMN "deleted_at";
-- reverse: modify "load_balancers" table
ALTER TABLE "load_balancers" DROP COLUMN "deleted_by", DROP COLUMN "deleted_at";
