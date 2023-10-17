-- +goose Up
-- modify "origins" table
ALTER TABLE "origins" ADD COLUMN "weight" bigint NOT NULL DEFAULT 100:::INT8;

-- +goose Down
-- reverse: modify "origins" table
ALTER TABLE "origins" DROP COLUMN "weight";
