// Copyright Infratographer, Inc. and/or licensed to Infratographer, Inc. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.
//
// Code generated by entc, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// LoadBalancersColumns holds the columns for the "load_balancers" table.
	LoadBalancersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Unique: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString, Size: 2147483647},
		{Name: "owner_id", Type: field.TypeString},
		{Name: "location_id", Type: field.TypeString},
		{Name: "provider_id", Type: field.TypeString},
	}
	// LoadBalancersTable holds the schema information for the "load_balancers" table.
	LoadBalancersTable = &schema.Table{
		Name:       "load_balancers",
		Columns:    LoadBalancersColumns,
		PrimaryKey: []*schema.Column{LoadBalancersColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "load_balancers_providers_provider",
				Columns:    []*schema.Column{LoadBalancersColumns[6]},
				RefColumns: []*schema.Column{ProvidersColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "loadbalancer_created_at",
				Unique:  false,
				Columns: []*schema.Column{LoadBalancersColumns[1]},
			},
			{
				Name:    "loadbalancer_updated_at",
				Unique:  false,
				Columns: []*schema.Column{LoadBalancersColumns[2]},
			},
			{
				Name:    "loadbalancer_provider_id",
				Unique:  false,
				Columns: []*schema.Column{LoadBalancersColumns[6]},
			},
			{
				Name:    "loadbalancer_location_id",
				Unique:  false,
				Columns: []*schema.Column{LoadBalancersColumns[5]},
			},
			{
				Name:    "loadbalancer_owner_id",
				Unique:  false,
				Columns: []*schema.Column{LoadBalancersColumns[4]},
			},
		},
	}
	// LoadBalancerAnnotationsColumns holds the columns for the "load_balancer_annotations" table.
	LoadBalancerAnnotationsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Unique: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "load_balancer_id", Type: field.TypeString},
	}
	// LoadBalancerAnnotationsTable holds the schema information for the "load_balancer_annotations" table.
	LoadBalancerAnnotationsTable = &schema.Table{
		Name:       "load_balancer_annotations",
		Columns:    LoadBalancerAnnotationsColumns,
		PrimaryKey: []*schema.Column{LoadBalancerAnnotationsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "load_balancer_annotations_load_balancers_load_balancer",
				Columns:    []*schema.Column{LoadBalancerAnnotationsColumns[3]},
				RefColumns: []*schema.Column{LoadBalancersColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "loadbalancerannotation_created_at",
				Unique:  false,
				Columns: []*schema.Column{LoadBalancerAnnotationsColumns[1]},
			},
			{
				Name:    "loadbalancerannotation_updated_at",
				Unique:  false,
				Columns: []*schema.Column{LoadBalancerAnnotationsColumns[2]},
			},
			{
				Name:    "loadbalancerannotation_load_balancer_id",
				Unique:  false,
				Columns: []*schema.Column{LoadBalancerAnnotationsColumns[3]},
			},
		},
	}
	// LoadBalancerStatusColumns holds the columns for the "load_balancer_status" table.
	LoadBalancerStatusColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Unique: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "source", Type: field.TypeString},
		{Name: "load_balancer_id", Type: field.TypeString},
	}
	// LoadBalancerStatusTable holds the schema information for the "load_balancer_status" table.
	LoadBalancerStatusTable = &schema.Table{
		Name:       "load_balancer_status",
		Columns:    LoadBalancerStatusColumns,
		PrimaryKey: []*schema.Column{LoadBalancerStatusColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "load_balancer_status_load_balancers_load_balancer",
				Columns:    []*schema.Column{LoadBalancerStatusColumns[4]},
				RefColumns: []*schema.Column{LoadBalancersColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "loadbalancerstatus_created_at",
				Unique:  false,
				Columns: []*schema.Column{LoadBalancerStatusColumns[1]},
			},
			{
				Name:    "loadbalancerstatus_updated_at",
				Unique:  false,
				Columns: []*schema.Column{LoadBalancerStatusColumns[2]},
			},
			{
				Name:    "loadbalancerstatus_load_balancer_id",
				Unique:  false,
				Columns: []*schema.Column{LoadBalancerStatusColumns[4]},
			},
			{
				Name:    "loadbalancerstatus_load_balancer_id_source",
				Unique:  false,
				Columns: []*schema.Column{LoadBalancerStatusColumns[4], LoadBalancerStatusColumns[3]},
			},
		},
	}
	// OriginsColumns holds the columns for the "origins" table.
	OriginsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Unique: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString},
		{Name: "target", Type: field.TypeString},
		{Name: "port_number", Type: field.TypeInt},
		{Name: "active", Type: field.TypeBool, Default: true},
		{Name: "pool_id", Type: field.TypeString},
	}
	// OriginsTable holds the schema information for the "origins" table.
	OriginsTable = &schema.Table{
		Name:       "origins",
		Columns:    OriginsColumns,
		PrimaryKey: []*schema.Column{OriginsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "origins_pools_pool",
				Columns:    []*schema.Column{OriginsColumns[7]},
				RefColumns: []*schema.Column{PoolsColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "origin_created_at",
				Unique:  false,
				Columns: []*schema.Column{OriginsColumns[1]},
			},
			{
				Name:    "origin_updated_at",
				Unique:  false,
				Columns: []*schema.Column{OriginsColumns[2]},
			},
			{
				Name:    "origin_pool_id",
				Unique:  false,
				Columns: []*schema.Column{OriginsColumns[7]},
			},
		},
	}
	// PoolsColumns holds the columns for the "pools" table.
	PoolsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Unique: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString},
		{Name: "protocol", Type: field.TypeEnum, Enums: []string{"tcp", "udp"}},
		{Name: "owner_id", Type: field.TypeString},
	}
	// PoolsTable holds the schema information for the "pools" table.
	PoolsTable = &schema.Table{
		Name:       "pools",
		Columns:    PoolsColumns,
		PrimaryKey: []*schema.Column{PoolsColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "pool_created_at",
				Unique:  false,
				Columns: []*schema.Column{PoolsColumns[1]},
			},
			{
				Name:    "pool_updated_at",
				Unique:  false,
				Columns: []*schema.Column{PoolsColumns[2]},
			},
			{
				Name:    "pool_owner_id",
				Unique:  false,
				Columns: []*schema.Column{PoolsColumns[5]},
			},
		},
	}
	// PortsColumns holds the columns for the "ports" table.
	PortsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Unique: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "number", Type: field.TypeInt},
		{Name: "name", Type: field.TypeString},
		{Name: "load_balancer_id", Type: field.TypeString},
	}
	// PortsTable holds the schema information for the "ports" table.
	PortsTable = &schema.Table{
		Name:       "ports",
		Columns:    PortsColumns,
		PrimaryKey: []*schema.Column{PortsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "ports_load_balancers_load_balancer",
				Columns:    []*schema.Column{PortsColumns[5]},
				RefColumns: []*schema.Column{LoadBalancersColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "port_created_at",
				Unique:  false,
				Columns: []*schema.Column{PortsColumns[1]},
			},
			{
				Name:    "port_updated_at",
				Unique:  false,
				Columns: []*schema.Column{PortsColumns[2]},
			},
			{
				Name:    "port_load_balancer_id",
				Unique:  false,
				Columns: []*schema.Column{PortsColumns[5]},
			},
			{
				Name:    "port_load_balancer_id_number",
				Unique:  true,
				Columns: []*schema.Column{PortsColumns[5], PortsColumns[3]},
			},
		},
	}
	// ProvidersColumns holds the columns for the "providers" table.
	ProvidersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Unique: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString},
		{Name: "owner_id", Type: field.TypeString},
	}
	// ProvidersTable holds the schema information for the "providers" table.
	ProvidersTable = &schema.Table{
		Name:       "providers",
		Columns:    ProvidersColumns,
		PrimaryKey: []*schema.Column{ProvidersColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "provider_created_at",
				Unique:  false,
				Columns: []*schema.Column{ProvidersColumns[1]},
			},
			{
				Name:    "provider_updated_at",
				Unique:  false,
				Columns: []*schema.Column{ProvidersColumns[2]},
			},
			{
				Name:    "provider_owner_id",
				Unique:  false,
				Columns: []*schema.Column{ProvidersColumns[4]},
			},
		},
	}
	// PoolPortsColumns holds the columns for the "pool_ports" table.
	PoolPortsColumns = []*schema.Column{
		{Name: "pool_id", Type: field.TypeString},
		{Name: "port_id", Type: field.TypeString},
	}
	// PoolPortsTable holds the schema information for the "pool_ports" table.
	PoolPortsTable = &schema.Table{
		Name:       "pool_ports",
		Columns:    PoolPortsColumns,
		PrimaryKey: []*schema.Column{PoolPortsColumns[0], PoolPortsColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "pool_ports_pool_id",
				Columns:    []*schema.Column{PoolPortsColumns[0]},
				RefColumns: []*schema.Column{PoolsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "pool_ports_port_id",
				Columns:    []*schema.Column{PoolPortsColumns[1]},
				RefColumns: []*schema.Column{PortsColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		LoadBalancersTable,
		LoadBalancerAnnotationsTable,
		LoadBalancerStatusTable,
		OriginsTable,
		PoolsTable,
		PortsTable,
		ProvidersTable,
		PoolPortsTable,
	}
)

func init() {
	LoadBalancersTable.ForeignKeys[0].RefTable = ProvidersTable
	LoadBalancerAnnotationsTable.ForeignKeys[0].RefTable = LoadBalancersTable
	LoadBalancerStatusTable.ForeignKeys[0].RefTable = LoadBalancersTable
	OriginsTable.ForeignKeys[0].RefTable = PoolsTable
	PortsTable.ForeignKeys[0].RefTable = LoadBalancersTable
	PoolPortsTable.ForeignKeys[0].RefTable = PoolsTable
	PoolPortsTable.ForeignKeys[1].RefTable = PortsTable
}
