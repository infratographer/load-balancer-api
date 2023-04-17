// Package api provides the API for the load balancers
package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sliceCompare(t *testing.T) {
	tests := []struct {
		name string
		s1   []string
		s2   []string
		want map[string]int
	}{
		{
			name: "happy path",
			s1:   []string{"p1", "p2", "p3"},
			s2:   []string{"p3", "p4"},
			want: map[string]int{
				"p1": -1,
				"p2": -1,
				"p3": 0,
				"p4": 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sliceCompare(tt.s1, tt.s2)
			assert.Equal(t, tt.want, got)
		})
	}
}
