--
-- +goose Up
-- +goose StatementBegin
CREATE TABLE load_balancers_status (
  status_id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
  load_balancer_id UUID NOT NULL REFERENCES load_balancers (load_balancer_id) ON UPDATE CASCADE,
  namespace STRING NOT NULL,
  source STRING NOT NULL,
  data JSONB NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  INDEX idx_status_id (status_id),
  INVERTED INDEX idx_status_data (status_id, namespace, data),
  UNIQUE INDEX idx_namespace_source_lb_id (namespace, source, load_balancer_id),
  INDEX idx_load_balancer_created_at (created_at),
  INDEX idx_load_balancer_updated_at (updated_at)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE load_balancers_status
-- +goose StatementEnd