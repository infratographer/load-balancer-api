-- +goose Up
-- modify "origins" table
ALTER TABLE "ports" ALTER COLUMN "name" DROP NOT NULL;

-- +goose Down
-- reverse: modify "origins" table
ALTER TABLE "ports" ALTER COLUMN "name" SET NOT NULL;