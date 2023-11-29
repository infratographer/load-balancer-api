package metadata

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.infratographer.com/load-balancer-api/pkg/client"
)

func TestGetLoadbalancerStatus(t *testing.T) {
	t.Run("valid status", func(t *testing.T) {
		statuses := client.MetadataStatuses{
			TotalCount: 2,
			Edges: []client.MetadataStatusEdges{
				{
					Node: client.MetadataStatusNode{
						StatusNamespaceID: "metasns-loadbalancer-status",
						Data:              json.RawMessage(`{"state": "active"}`),
					},
				},
				{
					Node: client.MetadataStatusNode{
						StatusNamespaceID: "metasns-some-other-namespace",
						Data:              json.RawMessage(`{"key": "value"}`),
					},
				},
			},
		}

		status, err := GetLoadbalancerStatus(statuses, "metasns-loadbalancer-status")
		require.Nil(t, err)
		assert.Equal(t, LoadBalancerStateActive, status.State)
	})

	t.Run("bad json data", func(t *testing.T) {
		statuses := client.MetadataStatuses{
			TotalCount: 1,
			Edges: []client.MetadataStatusEdges{
				{
					Node: client.MetadataStatusNode{
						StatusNamespaceID: "metasns-loadbalancer-status",
						Data:              json.RawMessage(`{"state"}`),
					},
				},
			},
		}

		status, err := GetLoadbalancerStatus(statuses, "metasns-loadbalancer-status")
		require.NotNil(t, err)
		require.Nil(t, status)
		assert.ErrorIs(t, err, ErrInvalidStatusData)
	})

	t.Run("status not found", func(t *testing.T) {
		statuses := client.MetadataStatuses{
			TotalCount: 0,
			Edges:      []client.MetadataStatusEdges{},
		}

		status, err := GetLoadbalancerStatus(statuses, "metasns-loadbalancer-status")
		require.NotNil(t, err)
		require.Nil(t, status)
		assert.ErrorIs(t, err, ErrStatusNotFound)
	})

	t.Run("no status data", func(t *testing.T) {
		statuses := client.MetadataStatuses{
			TotalCount: 1,
			Edges: []client.MetadataStatusEdges{
				{
					Node: client.MetadataStatusNode{
						StatusNamespaceID: "metasns-loadbalancer-status",
						Data:              json.RawMessage(``),
					},
				},
			},
		}

		status, err := GetLoadbalancerStatus(statuses, "metasns-loadbalancer-status")
		assert.NotNil(t, err)
		assert.Nil(t, status)
		assert.ErrorIs(t, err, ErrInvalidStatusData)
	})
}
