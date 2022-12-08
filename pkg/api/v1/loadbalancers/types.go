package loadbalancers

import (
	"time"

	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"
)

// LoadBalancer is a load balancer model
type LoadBalancer struct {
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *null.Time `json:"deleted_at,omitempty"`
	ID         uuid.UUID  `json:"id"`
	TenantID   uuid.UUID  `json:"tenant_id"`
	IPAddress  string     `json:"ip_address"`
	Name       string     `json:"display_name"`
	LocationID uuid.UUID  `json:"location_id"`
	Size       string     `json:"size"`
	Type       string     `json:"type"`
}

// LoadBalancers is a list of load balancers
type LoadBalancers []*LoadBalancer
