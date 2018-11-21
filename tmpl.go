package floridalist

const DYNOMITE_YML_TMPL = `
dyn_o_mite:
  #
  # name of the datacenter
  datacenter: {{.Datacenter}}
  #
  # name of the rack
  rack: {{.Rack}}
  #
  # port that dynomite nodes use to inter-communicate and gossip
  dyn_listen: {{.Address}}:8101
  #
  # seed provider implementation to provide a list of seed nodes
  dyn_seed_provider: florida_provider
  #
  #  listening address and port (name:port or ip:port) for this server pool
  listen: {{.Address}}:9101
  #
  # token(s) owned by a node
  tokens: '{{.Token}}'
  #
  # backend server ip:port:weight
  servers:
    - {{index .BackendServers 0}}
  #
  # pool speaks redis (0) or memcached (1) or other protocol. 
  data_store: 0
  # stats monitoring for admin ip:port
  stats_listen: {{.Address}}:2101
  # controls if server should be ejected temporarily when it fails consecutively 'server_failure_limit' times.
  auto_eject_hosts: true
  # number of consecutive failures on a server that would lead to it being temporarily ejected 
  server_failure_limit: 10
  # timeout value in msec to wait for before retrying on a temporarily ejected server
  server_retry_timeout: 30000
  # timeout value in msec that we wait for to establish a connection to the server or receive a response from a server.
  timeout: 5000
`
