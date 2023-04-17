--
-- +goose Up
-- +goose StatementBegin
CREATE TABLE load_balancers_metadata (
  metadata_id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
  load_balancer_id UUID NOT NULL REFERENCES load_balancers (load_balancer_id) ON UPDATE CASCADE,
  source STRING NOT NULL,
  data JSONB NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  INDEX idx_metadata_id (metadata_id),
  INVERTED INDEX idx_metadata_data (metadata_id, source, data),
  UNIQUE INDEX idx_source_lb_id (source, load_balancer_id),
  INDEX idx_load_balancer_created_at (created_at),
  INDEX idx_load_balancer_updated_at (updated_at)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE load_balancers_metadata;
-- +goose StatementEnd