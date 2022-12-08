package locations

import (
	"time"

	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"
)

// Location is the API model for a location
type Location struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *null.Time `json:"deleted_at,omitempty"`
	ID        uuid.UUID  `json:"id"`
	TenantID  uuid.UUID  `json:"tenant_id"`
	Name      string     `json:"display_name"`
}

// Locations is a slice of Location
type Locations []*Location
