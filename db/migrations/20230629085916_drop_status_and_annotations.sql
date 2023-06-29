-- +goose Up

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

-- +goose Down
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
