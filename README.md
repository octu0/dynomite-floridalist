# dynomite-floridalist

Dynomite seed management tool.

dynomite-floridalist works as Dynomite `seed_provider` using same communicate as Florida.

features below:

- Simple configuration management. `dynomite-florialist` can using dynomite.yml
- Container friendly. `dynomite-floridalist` can be applied to both sidecar patterns and centralized pattern
- Clustering and node management with (memberlist)[https://github.com/hashicorp/memberlist]

## Build

Build requires Go version 1.11+ installed.

```
$ go version
```

Run `make pkg` to Build and package for linux, darwin.

```
$ git clone https://github.com/octu0/dynomite-floridalist
$ make pkg
```

## Usage

using dynomite.yml (e.g. side-car pattern)

```
$ dynomite-floridalist -j memberlist:port -c /path/to/dynomite.conf
```

Or specify DC,Rack,Token

```
$ dynomite-floridalist --join member:port --address "$(hostname -i):8101" --datacenter asia-northeast1 --rack asia-northeast1-c --token 2147483647
```

## Help

```
T.B.D.
```
