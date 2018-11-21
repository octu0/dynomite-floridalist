package main

import(
  "os"
  "fmt"
  "strings"
  "text/template"

  "gopkg.in/urfave/cli.v1"

  "github.com/octu0/dynomite-floridalist"
)

type TemplateValue struct {
  Datacenter      string
  Rack            string
  Token           string
  Address         string
  BackendServers []string
}

func getTemplate(filename string) (*template.Template, error) {
  if filename != "" {
    return template.ParseFiles(filename)
  }
  s := strings.Trim(floridalist.DYNOMITE_YML_TMPL, "\n")
  return template.New("default.dynomite.yml.tmpl").Parse(s)
}

func generate_dynomite_yml_action(c *cli.Context) error {
  datacenter := c.String("datacenter")
  rack       := c.String("rack")
  token      := c.String("token")
  backendSvr := c.String("backend-server")
  address    := c.String("address")

  if datacenter == "" || rack == "" || token == "" || backendSvr == "" || address == "" {
    err := fmt.Errorf(
      `
      error: requires all value
      datacenter: '%s'
      rack: '%s'
      token: '%s'
      backend-server: '%s'
      address: '%s'
      `,
      datacenter,
      rack,
      token,
      backendSvr,
      address,
    )
    return err
  }

  tmpl, err := getTemplate(c.String("tmpl"))
  if err != nil {
    return err
  }

  var file *os.File
  output := c.String("output")
  if output == "" {
    file = os.Stdout
  } else {
    f, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE, 0644)
    if err != nil {
      return err
    }
    file = f
  }
  defer file.Close()

  args := TemplateValue{
    Datacenter: datacenter,
    Rack: rack,
    Token: token,
    Address: address,
    BackendServers: []string{ backendSvr },
  }
  if err := tmpl.Execute(file, args); err != nil {
    return err
  }

  return nil
}

func init(){
  AddCommand(cli.Command{
    Name: "generate",
    Usage: "generate `dynomite.yml`",
    Flags: []cli.Flag{
      cli.StringFlag{
        Name: "tmpl, t",
        Usage: "path to `dynomite.yml.tmpl`. if Empty using default value",
        EnvVar: "DYNOMITE_YAML_TMPL_PATH",
      },
      cli.StringFlag{
        Name: "output, o",
        Usage: "write output to `FILE`. defaults output to stdout",
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
        Name: "backend-server",
        Usage: "`servers` format 'ip:port:weight'",
        Value: "127.0.0.1:6379:100",
        EnvVar: "DYN_BACKEND_SERVER",
      },
    },
    Action: generate_dynomite_yml_action,
  })
}
