-- +goose Up
-- drop "provider_tenant_id" index on "providers" table
DROP INDEX "provider_tenant_id";
-- rename "tenant_id" column to "owner_id" in "providers" table
ALTER TABLE "providers" RENAME COLUMN "tenant_id" TO "owner_id";
-- create index "provider_owner_id" to table: "providers"
CREATE INDEX "provider_owner_id" ON "providers" ("owner_id");
-- drop "loadbalancer_tenant_id" to table: "load_balancers"
DROP INDEX "loadbalancer_tenant_id";
-- rename "tenant_id" column to "owner_id" in "load_balancers" table
ALTER TABLE "load_balancers" RENAME COLUMN "tenant_id" TO "owner_id";
-- create index "loadbalancer_owner_id" to table: "load_balancers"
CREATE INDEX "loadbalancer_owner_id" ON "load_balancers" ("owner_id");
-- drop "pool_tenant_id" to table: "pools"
DROP INDEX "pool_tenant_id";
-- rename "tenant_id" column to "owner_id" in "pools" table
ALTER TABLE "pools" RENAME COLUMN "tenant_id" TO "owner_id";
-- create index "pool_owner_id" to table: "pools"
CREATE INDEX "pool_owner_id" ON "pools" ("owner_id");

-- +goose Down
-- reverse: create index "provider_owner_id" to table: "providers"
DROP INDEX "provider_owner_id";
-- reverse: rename "tenant_id" column to "owner_id" in "providers" table
ALTER TABLE "providers" RENAME COLUMN "owner_id" TO "tenant_id";
-- reverse: drop "provider_tenant_id" index on "providers" table
CREATE INDEX "provider_tenant_id" ON "providers" ("tenant_id");
-- reverse: create index "loadbalancer_owner_id" to table: "load_balancers"
DROP INDEX "loadbalancer_owner_id";
-- reverse: rename "tenant_id" column to "owner_id" in "load_balancers" table
ALTER TABLE "load_balancers" RENAME COLUMN "owner_id" TO "tenant_id";
-- reverse: drop "loadbalancer_tenant_id" to table: "load_balancers"
CREATE INDEX "loadbalancer_tenant_id" ON "load_balancers" ("tenant_id");
-- reverse: create index "pool_owner_id" to table: "pools"
DROP INDEX "pool_owner_id";
-- reverse: rename "tenant_id" column to "owner_id" in "pools" table
ALTER TABLE "pools" RENAME COLUMN "owner_id" TO "tenant_id";
-- reverse: drop "pool_tenant_id" to table: "pools"
CREATE INDEX "pool_tenant_id" ON "pools" ("tenant_id");