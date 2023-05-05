package graphapi_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.infratographer.com/x/gidx"

	ent "go.infratographer.com/load-balancer-api/internal/ent/generated"
	"go.infratographer.com/load-balancer-api/internal/graphapi"
	"go.infratographer.com/load-balancer-api/internal/graphclient"
)

func TestQuery_pool1(t *testing.T) {
	ctx := context.Background()
	graphClient := graphclient.New(graphapi.NewResolver(EntClient))
	pool1 := (&PoolBuilder{}).MustNew(ctx)

	testCases := []struct {
		TestName     string
		QueryID      gidx.PrefixedID
		ExpectedPool *ent.Pool
		errorMsg     string
	}{
		{
			TestName:     "successful get of pool 1",
			QueryID:      pool1.ID,
			ExpectedPool: pool1,
		},
	}

	for _, tt := range testCases {
		// lint
		tt := tt

		t.Run(tt.TestName, func(t *testing.T) {
			resp, err := graphClient.QueryPool(tt.QueryID)
			if tt.errorMsg != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errorMsg)
				assert.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.EqualValues(t, tt.ExpectedPool.Name, resp.Name)
		})
	}
}
