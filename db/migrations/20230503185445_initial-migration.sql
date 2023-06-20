-- +goose Up
-- create "providers" table
CREATE TABLE "providers" ("id" character varying NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, "name" character varying NOT NULL, "tenant_id" character varying NOT NULL, PRIMARY KEY ("id"));
-- create index "provider_created_at" to table: "providers"
CREATE INDEX "provider_created_at" ON "providers" ("created_at");
-- create index "provider_tenant_id" to table: "providers"
CREATE INDEX "provider_tenant_id" ON "providers" ("tenant_id");
-- create index "provider_updated_at" to table: "providers"
CREATE INDEX "provider_updated_at" ON "providers" ("updated_at");
-- create "load_balancers" table
CREATE TABLE "load_balancers" ("id" character varying NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, "name" text NOT NULL, "tenant_id" character varying NOT NULL, "location_id" character varying NOT NULL, "provider_id" character varying NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "load_balancers_providers_provider" FOREIGN KEY ("provider_id") REFERENCES "providers" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
-- create index "loadbalancer_created_at" to table: "load_balancers"
CREATE INDEX "loadbalancer_created_at" ON "load_balancers" ("created_at");
-- create index "loadbalancer_location_id" to table: "load_balancers"
CREATE INDEX "loadbalancer_location_id" ON "load_balancers" ("location_id");
-- create index "loadbalancer_provider_id" to table: "load_balancers"
CREATE INDEX "loadbalancer_provider_id" ON "load_balancers" ("provider_id");
-- create index "loadbalancer_tenant_id" to table: "load_balancers"
CREATE INDEX "loadbalancer_tenant_id" ON "load_balancers" ("tenant_id");
-- create index "loadbalancer_updated_at" to table: "load_balancers"
CREATE INDEX "loadbalancer_updated_at" ON "load_balancers" ("updated_at");
-- create "load_balancer_annotations" table
CREATE TABLE "load_balancer_annotations" ("id" character varying NOT NULL, "namespace" character varying NOT NULL, "data" jsonb NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, "load_balancer_id" character varying NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "load_balancer_annotations_load_balancers_load_balancer" FOREIGN KEY ("load_balancer_id") REFERENCES "load_balancers" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
-- create index "loadbalancerannotation_created_at" to table: "load_balancer_annotations"
CREATE INDEX "loadbalancerannotation_created_at" ON "load_balancer_annotations" ("created_at");
-- create index "loadbalancerannotation_load_balancer_id" to table: "load_balancer_annotations"
CREATE INDEX "loadbalancerannotation_load_balancer_id" ON "load_balancer_annotations" ("load_balancer_id");
-- create index "loadbalancerannotation_load_balancer_id_namespace" to table: "load_balancer_annotations"
CREATE INDEX "loadbalancerannotation_load_balancer_id_namespace" ON "load_balancer_annotations" ("load_balancer_id", "namespace");
-- create index "loadbalancerannotation_namespace_data" to table: "load_balancer_annotations"
CREATE INDEX "loadbalancerannotation_namespace_data" ON "load_balancer_annotations" USING gin ("namespace", "data");
-- create index "loadbalancerannotation_updated_at" to table: "load_balancer_annotations"
CREATE INDEX "loadbalancerannotation_updated_at" ON "load_balancer_annotations" ("updated_at");
-- create "load_balancer_status" table
CREATE TABLE "load_balancer_status" ("id" character varying NOT NULL, "namespace" character varying NOT NULL, "data" jsonb NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, "source" character varying NOT NULL, "load_balancer_id" character varying NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "load_balancer_status_load_balancers_load_balancer" FOREIGN KEY ("load_balancer_id") REFERENCES "load_balancers" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
-- create index "loadbalancerstatus_created_at" to table: "load_balancer_status"
CREATE INDEX "loadbalancerstatus_created_at" ON "load_balancer_status" ("created_at");
-- create index "loadbalancerstatus_load_balancer_id" to table: "load_balancer_status"
CREATE INDEX "loadbalancerstatus_load_balancer_id" ON "load_balancer_status" ("load_balancer_id");
-- create index "loadbalancerstatus_load_balancer_id_namespace_source" to table: "load_balancer_status"
CREATE INDEX "loadbalancerstatus_load_balancer_id_namespace_source" ON "load_balancer_status" ("load_balancer_id", "namespace", "source");
-- create index "loadbalancerstatus_namespace_data" to table: "load_balancer_status"
CREATE INDEX "loadbalancerstatus_namespace_data" ON "load_balancer_status" USING gin ("namespace", "data");
-- create index "loadbalancerstatus_updated_at" to table: "load_balancer_status"
CREATE INDEX "loadbalancerstatus_updated_at" ON "load_balancer_status" ("updated_at");
-- create "pools" table
CREATE TABLE "pools" ("id" character varying NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, "name" character varying NOT NULL, "protocol" character varying NOT NULL, "tenant_id" character varying NOT NULL, PRIMARY KEY ("id"));
-- create index "pool_created_at" to table: "pools"
CREATE INDEX "pool_created_at" ON "pools" ("created_at");
-- create index "pool_tenant_id" to table: "pools"
CREATE INDEX "pool_tenant_id" ON "pools" ("tenant_id");
-- create index "pool_updated_at" to table: "pools"
CREATE INDEX "pool_updated_at" ON "pools" ("updated_at");
-- create "origins" table
CREATE TABLE "origins" ("id" character varying NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, "name" character varying NOT NULL, "target" character varying NOT NULL, "port_number" bigint NOT NULL, "active" boolean NOT NULL DEFAULT true, "pool_id" character varying NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "origins_pools_pool" FOREIGN KEY ("pool_id") REFERENCES "pools" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
-- create index "origin_created_at" to table: "origins"
CREATE INDEX "origin_created_at" ON "origins" ("created_at");
-- create index "origin_pool_id" to table: "origins"
CREATE INDEX "origin_pool_id" ON "origins" ("pool_id");
-- create index "origin_updated_at" to table: "origins"
CREATE INDEX "origin_updated_at" ON "origins" ("updated_at");
-- create "ports" table
CREATE TABLE "ports" ("id" character varying NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, "number" bigint NOT NULL, "name" character varying NOT NULL, "load_balancer_id" character varying NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "ports_load_balancers_load_balancer" FOREIGN KEY ("load_balancer_id") REFERENCES "load_balancers" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
-- create index "port_created_at" to table: "ports"
CREATE INDEX "port_created_at" ON "ports" ("created_at");
-- create index "port_load_balancer_id" to table: "ports"
CREATE INDEX "port_load_balancer_id" ON "ports" ("load_balancer_id");
-- create index "port_load_balancer_id_number" to table: "ports"
CREATE UNIQUE INDEX "port_load_balancer_id_number" ON "ports" ("load_balancer_id", "number");
-- create index "port_updated_at" to table: "ports"
CREATE INDEX "port_updated_at" ON "ports" ("updated_at");
-- create "pool_ports" table
CREATE TABLE "pool_ports" ("pool_id" character varying NOT NULL, "port_id" character varying NOT NULL, PRIMARY KEY ("pool_id", "port_id"), CONSTRAINT "pool_ports_pool_id" FOREIGN KEY ("pool_id") REFERENCES "pools" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "pool_ports_port_id" FOREIGN KEY ("port_id") REFERENCES "ports" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);

-- +goose Down
-- reverse: create "pool_ports" table
DROP TABLE "pool_ports";
-- reverse: create index "port_created_at" to table: "ports"
DROP INDEX "port_created_at";
-- reverse: create index "port_load_balancer_id" to table: "ports"
DROP INDEX "port_load_balancer_id";
-- reverse: create index "port_load_balancer_id_number" to table: "ports"
DROP INDEX "port_load_balancer_id_number";
-- reverse: create index "port_updated_at" to table: "ports"
DROP INDEX "port_updated_at";
-- reverse: create "ports" table
DROP TABLE "ports";
-- reverse: create index "origin_created_at" to table: "origins"
DROP INDEX "origin_created_at";
-- reverse: create index "origin_pool_id" to table: "origins"
DROP INDEX "origin_pool_id";
-- reverse: create index "origin_updated_at" to table: "origins"
DROP INDEX "origin_updated_at";
-- reverse: create "origins" table
DROP TABLE "origins";
-- reverse: create index "pool_created_at" to table: "pools"
DROP INDEX "pool_created_at";
-- reverse: create index "pool_tenant_id" to table: "pools"
DROP INDEX "pool_tenant_id";
-- reverse: create index "pool_updated_at" to table: "pools"
DROP INDEX "pool_updated_at";
-- reverse: create "pools" table
DROP TABLE "pools";
-- reverse: create index "loadbalancerstatus_created_at" to table: "load_balancer_status"
DROP INDEX "loadbalancerstatus_created_at";
-- reverse: create index "loadbalancerstatus_load_balancer_id" to table: "load_balancer_status"
DROP INDEX "loadbalancerstatus_load_balancer_id";
-- reverse: create index "loadbalancerstatus_load_balancer_id_namespace_source" to table: "load_balancer_status"
DROP INDEX "loadbalancerstatus_load_balancer_id_namespace_source";
-- reverse: create index "loadbalancerstatus_namespace_data" to table: "load_balancer_status"
DROP INDEX "loadbalancerstatus_namespace_data";
-- reverse: create index "loadbalancerstatus_updated_at" to table: "load_balancer_status"
DROP INDEX "loadbalancerstatus_updated_at";
-- reverse: create "load_balancer_status" table
DROP TABLE "load_balancer_status";
-- reverse: create index "loadbalancerannotation_created_at" to table: "load_balancer_annotations"
DROP INDEX "loadbalancerannotation_created_at";
-- reverse: create index "loadbalancerannotation_load_balancer_id" to table: "load_balancer_annotations"
DROP INDEX "loadbalancerannotation_load_balancer_id";
-- reverse: create index "loadbalancerannotation_load_balancer_id_namespace" to table: "load_balancer_annotations"
DROP INDEX "loadbalancerannotation_load_balancer_id_namespace";
-- reverse: create index "loadbalancerannotation_namespace_data" to table: "load_balancer_annotations"
DROP INDEX "loadbalancerannotation_namespace_data";
-- reverse: create index "loadbalancerannotation_updated_at" to table: "load_balancer_annotations"
DROP INDEX "loadbalancerannotation_updated_at";
-- reverse: create "load_balancer_annotations" table
DROP TABLE "load_balancer_annotations";
-- reverse: create index "loadbalancer_created_at" to table: "load_balancers"
DROP INDEX "loadbalancer_created_at";
-- reverse: create index "loadbalancer_location_id" to table: "load_balancers"
DROP INDEX "loadbalancer_location_id";
-- reverse: create index "loadbalancer_provider_id" to table: "load_balancers"
DROP INDEX "loadbalancer_provider_id";
-- reverse: create index "loadbalancer_tenant_id" to table: "load_balancers"
DROP INDEX "loadbalancer_tenant_id";
-- reverse: create index "loadbalancer_updated_at" to table: "load_balancers"
DROP INDEX "loadbalancer_updated_at";
-- reverse: create "load_balancers" table
DROP TABLE "load_balancers";
-- reverse: create index "provider_created_at" to table: "providers"
DROP INDEX "provider_created_at";
-- reverse: create index "provider_tenant_id" to table: "providers"
DROP INDEX "provider_tenant_id";
-- reverse: create index "provider_updated_at" to table: "providers"
DROP INDEX "provider_updated_at";
-- reverse: create "providers" table
DROP TABLE "providers";