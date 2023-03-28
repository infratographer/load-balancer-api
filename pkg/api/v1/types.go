package api

import "time"

type assignment struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	ID        string     `json:"id"`
	PortID    string     `json:"port_id"`
	PoolID    string     `json:"pool_id"`
	TenantID  string     `json:"tenant_id"`
}

type assignmentSlice []*assignment

type port struct {
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	ID             string     `json:"id"`
	TenantID       string     `json:"tenant_id"`
	LoadBalancerID string     `json:"load_balancer_id"`
	Name           string     `json:"name"`
	AddressFamily  string     `json:"address_family"`
	Port           int64      `json:"port"`
	Pools          []string   `json:"pools"`
}

type portSlice []*port

type loadBalancer struct {
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	ID          string     `json:"id"`
	IPAddressID string     `json:"ip_address_id"`
	TenantID    string     `json:"tenant_id"`
	Name        string     `json:"name"`
	LocationID  string     `json:"location_id"`
	Size        string     `json:"load_balancer_size"`
	Type        string     `json:"load_balancer_type"`
	Ports       portSlice  `json:"ports"`
}

type loadBalancerSlice []*loadBalancer

type location struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	ID        string     `json:"id"`
	TenantID  string     `json:"tenant_id"`
	Name      string     `json:"name"`
}

type locationSlice []*location

type origin struct {
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	Port           int64      `json:"port"`
	OriginTarget   string     `json:"origin_target"`
	OriginDisabled bool       `json:"origin_disabled"`
}

type originSlice []*origin

type pool struct {
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt *time.Time   `json:"deleted_at,omitempty"`
	ID        string       `json:"id"`
	TenantID  string       `json:"tenant_id"`
	Name      string       `json:"name"`
	Protocol  string       `json:"protocol"`
	Origins   *originSlice `json:"origins"`
}

type poolSlice []*pool

type response struct {
	Version       string             `json:"version"`
	Kind          string             `json:"kind"`
	Assignments   *assignmentSlice   `json:"assignments,omitempty"`
	Ports         *portSlice         `json:"ports,omitempty"`
	LoadBalancers *loadBalancerSlice `json:"load_balancers,omitempty"`
	Locations     *locationSlice     `json:"locations,omitempty"`
	Origins       *originSlice       `json:"origins,omitempty"`
	Pools         *poolSlice         `json:"pools,omitempty"`
}
