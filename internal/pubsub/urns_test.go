package pubsub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewURN(t *testing.T) {
	tests := []struct {
		name string
		kind string
		id   string
		urn  string
	}{
		{
			name: "example urn",
			kind: "testThing",
			id:   "9def378e-be7b-4566-83b5-20ae8ccf99cb",
			urn:  "urn:infratographer:testThing:9def378e-be7b-4566-83b5-20ae8ccf99cb",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := newURN(tt.kind, tt.id)
			assert.Equal(t, tt.urn, out)
		})
	}
}
