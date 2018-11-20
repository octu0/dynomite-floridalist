package floridalist

import(
  "time"
)

type Config struct {
  DebugMode              bool
  VerboseMode            bool
  Procs                  int

  MemberlistNodeName     string
  MemberlistJoinAddress  string
  MemberlistBindIp       string
  MemberlistBindPort     int

  FloridaBindIP          string
  FloridaBindPort        int
  FloridaEndpoint        string
  HttpReadTimeout        time.Duration
  HttpWriteTimeout       time.Duration
  MemberlistLeaveTimeout time.Duration
}

type SeedAdvertise struct {
  Datacenter  string   `json:"d"`
  Rack        string   `json:"r"`
  Token       string   `json:"t"`
  Address     string   `json:"a"`
}

type DynomiteYaml struct {
  DynomiteConf   Dyn_o_mite   `yaml:"dyn_o_mite"`
}
type Dyn_o_mite struct {
  Datacenter     string       `yaml:"datacenter"`
  Rack           string       `yaml:"rack"`
  Listen         string       `yaml:"dyn_listen"`
  Tokens         string       `yaml:"tokens"`
}
