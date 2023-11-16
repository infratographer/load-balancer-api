package graphapi

import (
	"context"
	"encoding/json"

	metacli "go.infratographer.com/metadata-api/pkg/client"
	"go.infratographer.com/x/gidx"

	metastatus "go.infratographer.com/load-balancer-api/pkg/metadata"

	"go.infratographer.com/load-balancer-api/internal/config"
)

const metadataStatusSource = "load-balancer-api"

// Metadata interface for the metadata client
type Metadata interface {
	StatusUpdate(ctx context.Context, input *metacli.StatusUpdateInput) (*metacli.StatusUpdate, error)
}

// LoadBalancerStatusUpdate updates the state of a load balancer in the metadata service
func (r Resolver) LoadBalancerStatusUpdate(ctx context.Context, loadBalancerID gidx.PrefixedID, status *metastatus.LoadBalancerStatus) error {
	if r.metadata == nil {
		r.logger.Warnln("metadata client not configured")
		return nil
	}

	jsonBytes, err := json.Marshal(status)
	if err != nil {
		return err
	}

	if _, err := r.metadata.StatusUpdate(ctx, &metacli.StatusUpdateInput{
		NodeID:      loadBalancerID.String(),
		NamespaceID: config.AppConfig.Metadata.StatusNamespaceID.String(),
		Source:      metadataStatusSource,
		Data:        json.RawMessage(jsonBytes),
	}); err != nil {
		return err
	}

	return nil
}
