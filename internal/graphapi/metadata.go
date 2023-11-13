package graphapi

import (
	"context"
	"encoding/json"
	"fmt"

	metadata "go.infratographer.com/metadata-api/pkg/client"
	"go.infratographer.com/x/gidx"

	"go.infratographer.com/load-balancer-api/internal/config"
)

const metadataStatusSource = "load-balancer-api"

// Metadata interface for the metadata client
type Metadata interface {
	StatusUpdate(ctx context.Context, input *metadata.StatusUpdateInput) (*metadata.StatusUpdate, error)
}

// LoadBalancerStatusUpdate updates the state of a load balancer in the metadata service
func (r Resolver) LoadBalancerStatusUpdate(ctx context.Context, loadBalancerID gidx.PrefixedID, state metadata.LoadBalancerState) error {
	if r.metadata == nil {
		r.logger.Warnln("metadata client not configured")
		return nil
	}

	if _, err := r.metadata.StatusUpdate(ctx, &metadata.StatusUpdateInput{
		NodeID:      loadBalancerID.String(),
		NamespaceID: config.AppConfig.Metadata.StatusNamespaceID.String(),
		Source:      metadataStatusSource,
		Data:        json.RawMessage(fmt.Sprintf(`{"state": "%s"}`, state)),
	}); err != nil {
		return err
	}

	return nil
}
