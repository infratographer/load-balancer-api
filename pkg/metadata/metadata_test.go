package metadata

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.infratographer.com/load-balancer-api/pkg/client"
)

func TestGetLoadbalancerState(t *testing.T) {
	t.Run("valid status", func(t *testing.T) {
		statuses := client.MetadataStatuses{
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

		state, err := GetLoadbalancerState(statuses, "metasns-loadbalancer-status")
		assert.Nil(t, err)
		assert.Equal(t, LoadBalancerStateActive, state)
	})

	t.Run("bad json data", func(t *testing.T) {
		statuses := client.MetadataStatuses{
			Edges: []client.MetadataStatusEdges{
				{
					Node: client.MetadataStatusNode{
						StatusNamespaceID: "metasns-loadbalancer-status",
						Data:              json.RawMessage(`{"state"}`),
					},
				},
			},
		}

		state, err := GetLoadbalancerState(statuses, "metasns-loadbalancer-status")
		assert.NotNil(t, err)
		assert.Empty(t, state)
		assert.ErrorIs(t, err, ErrInvalidStatusData)
	})

	t.Run("status not found", func(t *testing.T) {
		statuses := client.MetadataStatuses{
			Edges: []client.MetadataStatusEdges{},
		}

		state, err := GetLoadbalancerState(statuses, "metasns-loadbalancer-status")
		assert.NotNil(t, err)
		assert.Empty(t, state)
		assert.ErrorIs(t, err, ErrStatusNotFound)
	})

	t.Run("unknown state", func(t *testing.T) {
		statuses := client.MetadataStatuses{
			Edges: []client.MetadataStatusEdges{
				{
					Node: client.MetadataStatusNode{
						StatusNamespaceID: "metasns-loadbalancer-status",
						Data:              json.RawMessage(`{"state": "unknown"}`),
					},
				},
			},
		}

		state, err := GetLoadbalancerState(statuses, "metasns-loadbalancer-status")
		assert.NotNil(t, err)
		assert.Empty(t, state)
		assert.ErrorIs(t, err, ErrUnknownLoadBalancerState{State: "unknown"})
	})
}
