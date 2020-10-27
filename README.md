# Consul-Query
Tool for simple queries to the consul. Use hashicorp official [package](https://github.com/hashicorp/consul)

## Build and install
```
make
make install
```

## Use
Configure by environment variables. See the official documentation for consul client.

```
CONSUL_HTTP_ADDR=consul.dev:8500 ./bin/consul-query service-name-in-consul -t tag1 -t tag2 -o json
```
