package main

import(
  "log"
  "fmt"
  "runtime"
  "net"
  "context"
  "time"
  "strings"
  "os"
  "os/signal"
  "syscall"
  "io/ioutil"

  "github.com/comail/colog"
  "gopkg.in/urfave/cli.v1"
  "gopkg.in/yaml.v2"

  "github.com/octu0/dynomite-floridalist"
)

var (
  Commands = make([]cli.Command, 0)
)
func AddCommand(cmd cli.Command){
  Commands = append(Commands, cmd)
}

func SeedAdvertise(c *cli.Context) (floridalist.SeedAdvertise, error) {
  adv      := floridalist.SeedAdvertise{}
  yamlPath := c.String("conf")
  if yamlPath != "" {
    if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
      log.Printf("error: conf '%s' is not exists: %s", yamlPath, err.Error())
      return adv, err
    }

    data, err := ioutil.ReadFile(yamlPath)
    if err != nil {
      log.Printf("error: conf '%s' read error: %s", yamlPath, err.Error())
      return adv, err
    }

    d := new(floridalist.DynomiteYaml)
    if err := yaml.Unmarshal(data, d); err != nil {
      log.Printf("error: yaml '%s' parse error: %s", yamlPath, err.Error())
      return adv, err
    }
    conf          := d.DynomiteConf
    adv.Datacenter = conf.Datacenter
    adv.Rack       = conf.Rack
    adv.Token      = conf.Tokens

    if strings.Contains(conf.Listen, "0.0.0.0") {
      msg := fmt.Sprintf(
        "%s does not support '0.0.0.0:PORT' address pattern(%s). please set 'hostname:port' or 'ip:port' (ip are lookupable)",
        floridalist.UA,
        conf.Listen,
      )
      log.Printf("error: %s", msg)
      return adv, fmt.Errorf(msg)
    }
    ip, err := net.ResolveTCPAddr("tcp4", conf.Listen)
    if err != nil {
      return adv, err
    }
    adv.Address = ip.String()
  }

  datacenter := c.String("datacenter")
  if datacenter != "" {
    adv.Datacenter = datacenter
  }
  rack := c.String("rack")
  if rack != "" {
    adv.Rack = rack
  }
  token := c.String("token")
  if token != "" {
    adv.Token = token
  }
  address := c.String("address")
  if address != "" {
    adv.Address = address
  }
  return adv, nil
}

func action(c *cli.Context) error {
  httpReadTimeout, err := time.ParseDuration(c.String("http-read-timeout"))
  if err != nil {
    return err
  }
  httpWriteTimeout, err := time.ParseDuration(c.String("http-write-timeout"))
  if err != nil {
    return err
  }
  leaveTimeout, err := time.ParseDuration(c.String("ml-leave-timeout"))
  if err != nil {
    return err
  }

  adv, err := SeedAdvertise(c)
  if err != nil {
    return err
  }
  if adv.Datacenter == "" || adv.Rack == "" || adv.Token == "" || adv.Address == "" {
    msg := fmt.Sprintf(
      `
      error: seed Advertise doesnt include empty value.
      Datacenter: '%s'
      Rack: '%s'
      Token: '%s'
      Address: '%s'
      `,
      adv.Datacenter,
      adv.Rack,
      adv.Token,
      adv.Address,
    )
    return fmt.Errorf(msg)
  }

  config := floridalist.Config{
    DebugMode:              c.Bool("debug"),
    VerboseMode:            c.Bool("verbose"),
    Procs:                  c.Int("procs"),
    MemberlistJoinAddress:  c.String("join"),
    MemberlistNodeName:     c.String("ml-name"),
    MemberlistBindIp:       c.String("ml-ip"),
    MemberlistBindPort:     c.Int("ml-port"),
    FloridaBindIP:          c.String("http-ip"),
    FloridaBindPort:        c.Int("http-port"),
    FloridaEndpoint:        c.String("request"),
    HttpReadTimeout:        httpReadTimeout,
    HttpWriteTimeout:       httpWriteTimeout,
    MemberlistLeaveTimeout: leaveTimeout,
    UseWANConfig:           c.Bool("wan"),
  }
  if config.Procs < 1 {
    config.Procs = 1
  }

  if config.DebugMode {
    colog.SetMinLevel(colog.LDebug)
    if config.VerboseMode {
      colog.SetMinLevel(colog.LTrace)
    }
  }

  ctx := context.Background()
  ctx  = context.WithValue(ctx, "config", config)
  ctx  = context.WithValue(ctx, "advertise", adv)
  ctx  = context.WithValue(ctx, "logger.http", floridalist.NewHttpLogger(config))
  ctx  = context.WithValue(ctx, "logger.member", floridalist.NewMemberLogger(config))

  m   := floridalist.NewMember(ctx)
  if err := m.Init(); err != nil {
    log.Printf("error: memberlist initialization failed. %s", err.Error())
    return err
  }
  if err := m.Join(); err != nil {
    log.Printf("error: cluster join error: %s", err.Error())
    return err
  }

  http        := floridalist.NewHttpServer(ctx, m)
  error_chan  := make(chan error, 0)
  stopService := func() error {
    if err := m.Leave(config.MemberlistLeaveTimeout); err != nil {
      log.Printf("error: %s", err.Error())
      return err
    }
    sctx, cancel := context.WithTimeout(ctx, 30 * time.Second);
    defer cancel()

    if err := http.Stop(sctx); err != nil {
      log.Printf("error: %s", err.Error())
      return err
    }
    return nil
  }

  dumpSeeds := func() {
    values := m.SeedList()
    buf := fmt.Sprintf(
      "seeds:%s\nat:%s\n",
      strings.Join(values, "|"),
      time.Now().Format("2006-01-02 15:04:05"),
    )
    // dump out to stdout
    os.Stdout.WriteString(buf)
  }

  go func(){
    if err := http.Start(context.TODO()); err != nil {
      error_chan <- err
    }
  }()

  signal_chan := make(chan os.Signal, 10)
  signal.Notify(signal_chan, syscall.SIGTERM)
  signal.Notify(signal_chan, syscall.SIGHUP)
  signal.Notify(signal_chan, syscall.SIGQUIT)
  signal.Notify(signal_chan, syscall.SIGINT)
  signal.Notify(signal_chan, syscall.SIGCONT)
  running := true
  var lastErr error
  for running {
    select {
    case err := <-error_chan:
      log.Printf("error: error has occurred: %s", err.Error())
      lastErr = err
      if e := stopService(); e != nil {
        lastErr = e
      }
      running = false
    case sig := <-signal_chan:
      if sig == syscall.SIGCONT {
        dumpSeeds()
        continue
      }
      log.Printf("info: signal trap(%s)", sig.String())
      if err := stopService(); err != nil {
        lastErr = err
      }
      running = false
    }
  }
  if lastErr == nil {
    log.Printf("info: shutdown successful")
    return nil
  }
  return lastErr
}

func main(){
  colog.SetDefaultLevel(colog.LDebug)
  colog.SetMinLevel(colog.LInfo)

  colog.SetFormatter(&colog.StdFormatter{
    Flag: log.Ldate | log.Ltime | log.Lshortfile,
  })
  colog.Register()

  app         := cli.NewApp()
  app.Version  = floridalist.Version
  app.Name     = floridalist.AppName
  app.Author   = ""
  app.Email    = ""
  app.Usage    = ""
  app.Action   = action
  app.Commands = Commands
  app.Flags    = []cli.Flag{
    cli.StringFlag{
      Name: "join, j",
      Usage: "join memberlist cluster address",
      Value: floridalist.DEFAULT_MEMBERLIST_JOIN_ADDR,
      EnvVar: "DYN_FLORIDALIST_JOIN_ADDR",
    },
    cli.StringFlag{
      Name: "ml-name",
      Usage: "memberlist name(defaults: `ml-ip`:`ml-port`)",
      EnvVar: "DYN_FLORIDALIST_NODE_NAME",
    },
    cli.StringFlag{
      Name: "ml-ip",
      Usage: "memberlist bind-ip",
      Value: floridalist.DEFAULT_MEMBERLIST_BIND_IP,
      EnvVar: "DYN_FLORIDALIST_BIND_IP",
    },
    cli.IntFlag{
      Name: "ml-port",
      Usage: "memberlist bind-port",
      Value: floridalist.DEFAULT_MEMBERLIST_BIND_PORT,
      EnvVar: "DYN_FLORIDALIST_BIND_PORT",
    },
    cli.StringFlag{
      Name: "conf, c",
      Usage: "path to dynomite.yml",
      EnvVar: "DYNOMITE_YAML_PATH",
    },
    cli.StringFlag{
      Name: "address",
      Usage: "Dynomite node listen address",
      EnvVar: "DYN_ADDRESS",
    },
    cli.StringFlag{
      Name: "datacenter",
      Usage: "Dynomite node datacenter name",
      EnvVar: "DYN_DC",
    },
    cli.StringFlag{
      Name: "rack",
      Usage: "Dynomite node rack name",
      EnvVar: "DYN_RACK",
    },
    cli.StringFlag{
      Name: "token",
      Usage: "Dynomite node owned token",
      EnvVar: "DYN_TOKEN",
    },
    cli.StringFlag{
      Name: "http-ip, i",
      Usage: "florida API http ip",
      Value: floridalist.DEFAULT_FLORIDA_API_IP,
      EnvVar: "DYNOMITE_FLORIDA_IP",
    },
    cli.IntFlag{
      Name: "http-port, p",
      Usage: "florida API http port",
      Value: floridalist.DEFAULT_FLORIDA_API_PORT,
      EnvVar: "DYNOMITE_FLORIDA_PORT",
    },
    cli.StringFlag{
      Name: "request, r",
      Usage: "florida API request endpoint",
      Value: floridalist.DEFAULT_FLORIDA_API_REQUEST,
      EnvVar: "DYNOMITE_FLORIDA_REQUEST",
    },
    cli.StringFlag{
      Name: "http-read-timeout",
      Usage: "florida API read timeout",
      Value: floridalist.DEFAULT_HTTP_READ_TIMEOUT,
    },
    cli.StringFlag{
      Name: "http-write-timeout",
      Usage: "florida API write timeout",
      Value: floridalist.DEFAULT_HTTP_WRITE_TIMEOUT,
    },
    cli.StringFlag{
      Name: "ml-leave-timeout",
      Usage: "memberlist cluster leave timeout",
      Value: floridalist.DEFAULT_MEMBERLIST_LEAVE_TIMEOUT,
    },
    cli.BoolFlag{
      Name: "wan",
      Usage: "use WANconfig (defaults LANConfig)",
    },
    cli.IntFlag{
      Name: "procs, P",
      Usage: "attach cpu(s)",
      Value: runtime.NumCPU(),
    },
    cli.BoolFlag{
      Name: "debug, d",
      Usage: "debug mode",
    },
    cli.BoolFlag{
      Name: "verbose, V",
      Usage: "verbose. more message",
    },
  }
  if err := app.Run(os.Args); err != nil {
    log.Printf("error: %s", err.Error())
    cli.OsExiter(1)
  }
}
