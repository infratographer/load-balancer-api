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
printf "$da_payload" $tenant_id | curl -s -X POST -H "Content-Type: application/json" -d @- http://localhost:7608/api/v1/tenant/${tenant_id}/locations | jq

# Get the location id
location_id=$(curl -s -H "Content-Type: application/json" http://localhost:7608/api/v1/tenant/${tenant_id}/locations/DA | jq -r .location.id)


# json for creating a load balancer
read -r -d '' lb_payload <<'EOF'
[{
    "display_name":"LB-01",
    "tenant_id":"%s",
    "location_id":"%s",
    "ip_addr":"1.1.1.1",
    "load_balancer_size":"small",
    "load_balancer_type":"layer-3"
}]
EOF

# Create a load balancer
printf "$lb_payload" $tenant_id $location_id | curl -X POST -H "Content-Type: application/json" -d @- http://localhost:7608/api/v1/tenant/${tenant_id}/loadbalancers | jq

# json for creating a load balancer
read -r -d '' lb_payload <<'EOF'
[{
    "display_name":"LB-02",
    "tenant_id":"%s",
    "location_id":"%s",
    "ip_addr":"1.2.1.1",
    "load_balancer_size":"small",
    "load_balancer_type":"layer-3"
}]
EOF

printf "$lb_payload" $tenant_id $location_id | curl -X POST -H "Content-Type: application/json" -d @- http://localhost:7608/api/v1/tenant/${tenant_id}/loadbalancers | jq


# Get the load balancer id
lb_id=$(curl -s -H "Content-Type: application/json" http://localhost:7608/api/v1/tenant/${tenant_id}/loadbalancers/1.1.1.1 | jq -r .load_balancer.id)

curl -s -H "Content-Type: application/json" http://localhost:7608/api/v1/tenant/${tenant_id}/loadbalancers/${lb_id}| jq

curl -s -H "Content-Type: application/json" http://localhost:7608/api/v1/tenant/${tenant_id}/loadbalancers?ip_addr=1.2.1.1 -X DELETE | jq
