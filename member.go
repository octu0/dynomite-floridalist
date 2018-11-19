package floridalist

import(
  "log"
  "fmt"
  "time"
  "context"
  "encoding/json"
  "strings"

  "github.com/hashicorp/memberlist"
)

type DefDelegate struct {
  meta          []byte
  msgs          [][]byte
  broadcasts    [][]byte
  state         []byte
  remoteState   []byte
  adv           SeedAdvertise
}

func NewMemberlistDelegate(ctx context.Context) *DefDelegate {
  adv := ctx.Value("advertise").(SeedAdvertise)

  d     := new(DefDelegate)
  d.adv  = adv
  return d
}
func (d *DefDelegate) Init() error {
  data, err := json.Marshal(d.adv)
  if err != nil {
    log.Printf("error: advertise meta marshal error: %s", err.Error())
    return err
  }
  d.meta = data
  return nil
}
func (d *DefDelegate) NodeMeta(limit int) []byte {
  return d.meta
}
func (d *DefDelegate) NotifyMsg(msg []byte) {
  // not use
}
func (d *DefDelegate) GetBroadcasts(overhead, limit int) [][]byte {
  // not use, noop
  return d.broadcasts
}
func (d *DefDelegate) LocalState(join bool) []byte {
  // not use, noop
  return d.state
}
func (d *DefDelegate) MergeRemoteState(buf []byte, join bool) {
  // not use
}

func NewMemberlistConfig(ctx context.Context) *memberlist.Config {
  config := ctx.Value("config").(Config)
  logger := ctx.Value("logger.member").(*MemberLogger)
  name   := fmt.Sprintf("%s:%d", config.MemberlistBindIp, config.MemberlistBindPort)

  c                  := memberlist.DefaultLANConfig()
  c.Name              = name
  c.BindAddr          = config.MemberlistBindIp
  c.BindPort          = config.MemberlistBindPort
  c.AdvertiseAddr     = config.MemberlistBindIp
  c.AdvertisePort     = config.MemberlistBindPort
  c.Logger            = logger.NewLogger()

  return c
}

type Member struct {
  joinAddr  string
  joined    bool
  delegate  *DefDelegate
  mconf     *memberlist.Config
  mlist     *memberlist.Memberlist
}

func NewMember(ctx context.Context) *Member {
  config := ctx.Value("config").(Config)

  m := new(Member)
  m.joinAddr  = config.MemberlistJoinAddress
  m.joined    = false
  m.mconf     = NewMemberlistConfig(ctx)
  m.delegate  = NewMemberlistDelegate(ctx)
  return m
}
func (m *Member) Init() error {
  if err := m.delegate.Init(); err != nil {
    return err
  }
  m.mconf.Delegate = m.delegate

  mlist, err := memberlist.Create(m.mconf)
  if err != nil {
    log.Printf("error: Failed to created memberlist: " + err.Error())
    return err
  }
  m.mlist = mlist
  return nil
}
func (m *Member) IsJoined() bool {
  return m.joined
}
func (m *Member) SeedList() []string {
  members := m.mlist.Members()
  slist   := make([]string, 0)
  for _, m := range members {
    adv := new(SeedAdvertise)
    if err := json.Unmarshal(m.Meta, adv); err != nil {
      log.Printf("error: node(%s) meta unmarshal error: %s", m.Name, err.Error())
      continue
    }

    //
    // Floria peer format)
    //   peer_host:peer_listen_port:rack:dc:peer_token
    //
    seed := []string{
      adv.Address,    // Address: 'ip:port'
      adv.Rack,       // Rack: 'asia-northeast1-a'
      adv.Datacenter, // Datacenter: 'asia-northeast1'
      adv.Token,      // Token: '2147483647'
    }
    slist = append(slist, strings.Join(seed, ":"))
  }
  return slist
}
func (m *Member) Join() error {
  if _, err := m.mlist.Join([]string{m.joinAddr}); err != nil {
    log.Printf("error: Failed to join cluster(%s): %s", m.joinAddr, err.Error())
    return err
  }
  m.joined = true
  log.Printf("info: cluster(%s) join successful.", m.joinAddr)

  for _, m := range m.mlist.Members() {
    log.Printf("info: available member %s(%s:%d)", m.Name, m.Addr.String(), m.Port);
  }
  return nil
}
func (m *Member) Leave(timeout time.Duration) error {
  if err := m.mlist.Leave(timeout); err != nil {
    log.Printf("error: Failed to leave cluster(%s): %s", m.joinAddr, err.Error())
    return err
  }
  m.joined = false
  return m.mlist.Shutdown()
}
