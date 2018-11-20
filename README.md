# dynomite-floridalist

[Dynomite](https://github.com/Netflix/dynomite) seed management / seed_provider.

dynomite-floridalist works as Dynomite `seed_provider` using same communicate as Florida.

features below:

- Simple configuration management. `dynomite-florialist` can using dynomite.yml
- Container friendly. `dynomite-floridalist` can be applied to both sidecar patterns and centralized pattern
- Clustering and node management with [memberlist](https://github.com/hashicorp/memberlist)

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
NAME:
   dynomite-floridalist

USAGE:
   dynomite-floridalist [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --join value, -j value       join memberlist cluster address (default: "127.0.0.1:3101") [$DYN_FLORIDALIST_JOIN_ADDR]
   --ml-name ml-ip              memberlist name(defaults: ml-ip:`ml-port`) [$DYN_FLORIDALIST_NODE_NAME]
   --ml-ip value                memberlist bind-ip (default: "0.0.0.0") [$DYN_FLORIDALIST_BIND_IP]
   --ml-port value              memberlist bind-port (default: 3101) [$DYN_FLORIDALIST_BIND_PORT]
   --conf value, -c value       path to dynomite.yml
   --address value              Dynomite node listen address
   --datacenter value           Dynomite node datacenter name
   --rack value                 Dynomite node rack name
   --token value                Dynomite node owned token
   --http-ip value, -i value    florida API http ip (default: "0.0.0.0") [$DYNOMITE_FLORIDA_IP]
   --http-port value, -p value  florida API http port (default: 8080) [$DYNOMITE_FLORIDA_PORT]
   --request value, -r value    florida API request endpoint (default: "/REST/v1/admin/get_seeds") [$DYNOMITE_FLORIDA_REQUEST]
   --http-read-timeout value    florida API read timeout (default: "500ms")
   --http-write-timeout value   florida API write timeout (default: "500ms")
   --ml-leave-timeout value     memberlist cluster leave timeout (default: "30s")
   --procs value, -P value      attach cpu(s) (default: 8)
   --debug, -d                  debug mode
   --verbose, -V                verbose. more message
   --help, -h                   show help
   --version, -v                print the version
```
