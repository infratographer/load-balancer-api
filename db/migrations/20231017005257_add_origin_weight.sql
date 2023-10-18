-- +goose Up
-- modify "origins" table
ALTER TABLE "origins" ADD COLUMN "weight" int NOT NULL DEFAULT 100;

-- +goose Down
-- reverse: modify "origins" table
ALTER TABLE "origins" DROP COLUMN "weight";
