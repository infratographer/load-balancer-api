// Package config provides a struct to store the applications config
package config

import (
	"time"

	"go.infratographer.com/permissions-api/pkg/permissions"
	"go.infratographer.com/x/crdbx"
	"go.infratographer.com/x/echojwtx"
	"go.infratographer.com/x/echox"
	"go.infratographer.com/x/events"
	"go.infratographer.com/x/gidx"
	"go.infratographer.com/x/loggingx"
	"go.infratographer.com/x/oauth2x"
	"go.infratographer.com/x/otelx"
)

// AppConfig stores all the config values for our application
var AppConfig struct {
	OIDC            echojwtx.AuthConfig `mapstructure:"oidc"`
	OIDCClient      OIDCClientConfig    `mapstructure:"oidc"`
	CRDB            crdbx.Config
	Logging         loggingx.Config
	Server          echox.Config
	Tracing         otelx.Config
	Events          events.Config
	Permissions     permissions.Config
	Metadata        Metadata
	RestrictedPorts []int

	SupergraphURL string `mapstructure:"supergraph-url"`
}

// Metadata stores the configuration for metadata client
type Metadata struct {
	StatusNamespaceID gidx.PrefixedID `mapstructure:"status-namespace-id"`
	Timeout           time.Duration
}

// OIDCClientConfig stores the configuration for an OIDC client
type OIDCClientConfig struct {
	Config oauth2x.Config `mapstructure:"client"`
}
