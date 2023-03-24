[![Build status](https://badge.buildkite.com/e73aabfa52747216e2c5fd2f02ac801f034ec0f2e982905443.svg)](https://buildkite.com/infratographer/load-balancer-api)

# Load Balancer API

This is an API for managing load balancer configurations for the Infratographer platform.

## Load Balancer Components

Load balancer --> Port --> Pool --> Origin

### Load Balancers

Load balancers are the top level component managed by the load balancer API.  Load balancers are assigned to a tenant.

### Ports

Ports are the listening port of the load balancer. Ports are assigned to a load balancer.

### Pools

Pools are a collection of origins. Pools are assigned to a tenant, and linked to a port through assignments.

### Origins

Origins define a backend service IP and port. Origins belong to a pool.

### Assignments

Assignments link ports to pools.

## Development and Contributing

* [Development Guide](docs/development.md)
* [Contributing](https://infratographer.com/community/contributing/)

## Code of Conduct

[Contributor Code of Conduct](https://infratographer.com/community/code-of-conduct/). By participating in this project you agree to abide by its terms.

## Contact

To contact the maintainers, please open a [GithHub Issue](https://github.com/infratographer/load-balancer-api/issues/new)

## License

[Apache 2.0](LICENSE)
