# dynomite-floridalist

[Dynomite](https://github.com/Netflix/dynomite) seed management / seed_provider.

dynomite-floridalist works as Dynomite `seed_provider` using same communicate as Florida.

features below:

- Simple configuration management. `dynomite-florialist` can using `dynomite.yml`
- Container friendly. `dynomite-floridalist` can be applied to both sidecar patterns and centralized pattern
- `dynomite.yml` generator
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
   1.1.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --join value, -j value       join memberlist cluster address (default: "127.0.0.1:3101") [$DYN_FLORIDALIST_JOIN_ADDR]
   --ml-name ml-ip              memberlist name(defaults: ml-ip:`ml-port`) [$DYN_FLORIDALIST_NODE_NAME]
   --ml-ip value                memberlist bind-ip (default: "0.0.0.0") [$DYN_FLORIDALIST_BIND_IP]
   --ml-port value              memberlist bind-port (default: 3101) [$DYN_FLORIDALIST_BIND_PORT]
   --conf value, -c value       path to dynomite.yml [$DYNOMITE_YAML_PATH]
   --address value              Dynomite node listen address [$DYN_ADDRESS]
   --datacenter value           Dynomite node datacenter name [$DYN_DC]
   --rack value                 Dynomite node rack name [$DYN_RACK]
   --token value                Dynomite node owned token [$DYN_TOKEN]
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

## Generate dynomite.yml

Using `generate` command, can make `dynomite.yml` automatically.

```
NAME:
   dynomite-floridalist generate - generate `dynomite.yml`

USAGE:
   dynomite-floridalist generate [command options] [arguments...]

OPTIONS:
   --tmpl dynomite.yml.tmpl, -t dynomite.yml.tmpl  path to dynomite.yml.tmpl. if Empty using default value [$DYNOMITE_YAML_TMPL_PATH]
   --output FILE, -o FILE                          write output to FILE. defaults output to stdout [$DYNOMITE_YAML_PATH]
   --address value                                 Dynomite node listen address [$DYN_ADDRESS]
   --datacenter value                              Dynomite node datacenter name [$DYN_DC]
   --rack value                                    Dynomite node rack name [$DYN_RACK]
   --token value                                   Dynomite node owned token [$DYN_TOKEN]
   --backend-server servers                        servers format 'ip:port:weight' (default: "127.0.0.1:6379:100") [$DYN_BACKEND_SERVER]
```

Usage

```
$ DYNOMITE_YAML="/etc/dynomite.yml"

$ dynomite-floridalist generate \
	--address 10.16.0.123 \
	--datacenter ap-northeast-1 \
	--rack ap-northeast-1d
	--token 0 \
	--backend-server "127.0.0.1:6379:1" \
	-o $DYNOMITE_YAML

$ dynomite-floridalist --address 10.16.0.123 -c $DYNOMITE_YAML
```

Environment 

```
# defined by Dockerfile or docker-compose or any solution
$ DYN_ADDRESS="10.16.0.123"
$ DYN_DC="ap-northeast-1"
$ DYN_RACK="ap-northeast-1d"
$ DYN_TOKEN="0"
$ DYN_BACKEND_SERVER="127.0.0.1:6379:100"

# run. same source
$ dynomite-floridalist generate -o /etc/dynomite.yml
$ dynomite-floridalist -c /etc/dynomite.yml
$ dynomite -c /etc/dynomite.yml
```
