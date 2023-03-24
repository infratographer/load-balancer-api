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
		want string
	}{
		{
			name: "example urn",
			kind: "testThing",
			id:   "9def378e-be7b-4566-83b5-20ae8ccf99cb",
			want: "urn:infratographer:testThing:9def378e-be7b-4566-83b5-20ae8ccf99cb",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newURN(tt.kind, tt.id)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_newURN(t *testing.T) {
	tests := []struct {
		name string
		kind string
		id   string
		want string
	}{
		{
			name: "example",
			kind: "foo",
			id:   "8cb89124-7954-4c98-85d5-e1fad6e3d723",
			want: "urn:infratographer:foo:8cb89124-7954-4c98-85d5-e1fad6e3d723",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newURN(tt.kind, tt.id)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewTenantURN(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want string
	}{
		{
			name: "example",
			id:   "8cb89124-7954-4c98-85d5-e1fad6e3d723",
			want: "urn:infratographer:tenant:8cb89124-7954-4c98-85d5-e1fad6e3d723",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTenantURN(tt.id)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewLoadBalancerURN(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want string
	}{
		{
			name: "example",
			id:   "8cb89124-7954-4c98-85d5-e1fad6e3d723",
			want: "urn:infratographer:load-balancer:8cb89124-7954-4c98-85d5-e1fad6e3d723",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewLoadBalancerURN(tt.id)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewPortURN(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want string
	}{
		{
			name: "example",
			id:   "8cb89124-7954-4c98-85d5-e1fad6e3d723",
			want: "urn:infratographer:load-balancer-port:8cb89124-7954-4c98-85d5-e1fad6e3d723",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPortURN(tt.id)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewAssignmentURN(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want string
	}{
		{
			name: "example",
			id:   "8cb89124-7954-4c98-85d5-e1fad6e3d723",
			want: "urn:infratographer:load-balancer-assignment:8cb89124-7954-4c98-85d5-e1fad6e3d723",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAssignmentURN(tt.id)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewOriginURN(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want string
	}{
		{
			name: "example",
			id:   "8cb89124-7954-4c98-85d5-e1fad6e3d723",
			want: "urn:infratographer:load-balancer-origin:8cb89124-7954-4c98-85d5-e1fad6e3d723",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewOriginURN(tt.id)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewPoolURN(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want string
	}{
		{
			name: "example",
			id:   "8cb89124-7954-4c98-85d5-e1fad6e3d723",
			want: "urn:infratographer:load-balancer-pool:8cb89124-7954-4c98-85d5-e1fad6e3d723",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPoolURN(tt.id)
			assert.Equal(t, tt.want, got)
		})
	}
}
