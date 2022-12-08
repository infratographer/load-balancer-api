#!/bin/bash
set -x

make dev-database

# Create a tenant
tenant_id=$(uuidgen)

# json for creating a location
read -r -d '' da_payload <<'EOF'
{
    "display_name":"DA",
    "tenant_id":"%s"
}
EOF

# Create a location
printf "$da_payload" $tenant_id | curl -X POST -H "Content-Type: application/json" -d @- http://localhost:8000/api/v1/locations | jq

# Get the location id
location_id=$(curl -H "Content-Type: application/json" http://localhost:8000/api/v1/tenant/${tenant_id}/locations/DA | jq -r .location.id)

# json for creating a load balancer
read -r -d '' lb_payload <<'EOF'
{
    "display_name":"LB",
    "tenant_id":"%s",
    "location_id":"%s",
    "ip_address":"1.1.1.1",
    "size":"small",
    "type":"layer-3"
}
EOF

# Create a load balancer
printf "$lb_payload" $tenant_id $location_id | curl -X POST -H "Content-Type: application/json" -d @- http://localhost:8000/api/v1/loadbalancers | jq

# Get the load balancer id
lb_id=$(curl -H "Content-Type: application/json" http://localhost:8000/api/v1/tenant/${tenant_id}/loadbalancers/1.1.1.1 | jq -r .load_balancer.id)

curl -H "Content-Type: application/json" http://localhost:8000/api/v1/tenant/${tenant_id}/loadbalancers/1.1.1.1 -X DELETE | jq