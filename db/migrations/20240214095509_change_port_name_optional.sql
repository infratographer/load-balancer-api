-- +goose Up
-- modify "ports" table
ALTER TABLE "ports" ALTER COLUMN "name" DROP NOT NULL;

-- +goose Down
-- reverse: modify "ports" table
ALTER TABLE "ports" ALTER COLUMN "name" SET NOT NULL;
