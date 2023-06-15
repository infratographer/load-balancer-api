BEGIN TRANSACTION;
-- Drop "provider_tenant_id" index on "providers" table
DROP INDEX "provider_tenant_id";
-- Rename "tenant_id" column to "owner_id" in "providers" table
ALTER TABLE "providers" RENAME COLUMN "tenant_id" TO "owner_id";
-- Create index "provider_owner_id" to table: "providers"
CREATE INDEX "provider_owner_id" ON "providers" ("owner_id");

-- Drop "loadbalancer_tenant_id" to table: "load_balancers"
DROP INDEX "loadbalancer_tenant_id";
-- Rename "tenant_id" column to "owner_id" in "load_balancers" table
ALTER TABLE "load_balancers" RENAME COLUMN "tenant_id" TO "owner_id";
-- Create index "loadbalancer_owner_id" to table: "load_balancers"
CREATE INDEX "loadbalancer_owner_id" ON "load_balancers" ("owner_id");

-- Drop "pool_tenant_id" to table: "pools"
DROP INDEX "pool_tenant_id";
-- Rename "tenant_id" column to "owner_id" in "pools" table
ALTER TABLE "pools" RENAME COLUMN "tenant_id" TO "owner_id";
-- Create index "pool_owner_id" to table: "pools"
CREATE INDEX "pool_owner_id" ON "load_balancers" ("owner_id");
COMMIT;
