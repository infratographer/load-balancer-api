-- +goose Up
-- modify "load_balancers" table
ALTER TABLE "load_balancers" ADD COLUMN "created_by" character varying NULL, ADD COLUMN "updated_by" character varying NULL;
-- modify "origins" table
ALTER TABLE "origins" ADD COLUMN "created_by" character varying NULL, ADD COLUMN "updated_by" character varying NULL;
-- modify "pools" table
ALTER TABLE "pools" ADD COLUMN "created_by" character varying NULL, ADD COLUMN "updated_by" character varying NULL;
-- modify "ports" table
ALTER TABLE "ports" ADD COLUMN "created_by" character varying NULL, ADD COLUMN "updated_by" character varying NULL;
-- modify "providers" table
ALTER TABLE "providers" ADD COLUMN "created_by" character varying NULL, ADD COLUMN "updated_by" character varying NULL;

-- +goose Down
-- reverse: modify "providers" table
ALTER TABLE "providers" DROP COLUMN "updated_by", DROP COLUMN "created_by";
-- reverse: modify "ports" table
ALTER TABLE "ports" DROP COLUMN "updated_by", DROP COLUMN "created_by";
-- reverse: modify "pools" table
ALTER TABLE "pools" DROP COLUMN "updated_by", DROP COLUMN "created_by";
-- reverse: modify "origins" table
ALTER TABLE "origins" DROP COLUMN "updated_by", DROP COLUMN "created_by";
-- reverse: modify "load_balancers" table
ALTER TABLE "load_balancers" DROP COLUMN "updated_by", DROP COLUMN "created_by";
