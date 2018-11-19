package floridalist

import(
  "log"
  "os"
  "strings"
  "strconv"

  "github.com/comail/colog"
)

const TAB string = "\t"

type HttpLogger interface {
  Write(host string, method string, uri string, status int, ua string)
}

type DefLogger struct {
  DebugMode  bool
}

func NewHttpLogger(c Config) HttpLogger {
  l := new(DefLogger)
  l.DebugMode = c.DebugMode
  return l
}
func (l *DefLogger) Write(host string, method string, uri string, status int, ua string) {
  msg := []string{
    "host:", host,
    TAB,
    "method:", method,
    TAB,
    "uri:", uri,
    TAB,
    "status:", strconv.Itoa(status),
    TAB,
    "ua:", ua,
  }
  m := strings.Join(msg, "")
  if l.DebugMode {
    log.Printf("debug: %s", m)
  } else {
    log.Printf("info: %s", m)
  }
}

type MemberLogger struct {
  c  *colog.CoLog
}
func NewMemberLogger(config Config) *MemberLogger {
  c := colog.NewCoLog(os.Stdout, "member ", log.Ldate | log.Ltime | log.Lshortfile)
  c.SetDefaultLevel(colog.LDebug)
  c.SetMinLevel(colog.LInfo)
  if config.DebugMode {
    c.SetMinLevel(colog.LDebug)
    if config.VerboseMode {
      c.SetMinLevel(colog.LTrace)
    }
  }

  ml := new(MemberLogger)
  ml.c = c
  return ml
}
func (ml *MemberLogger) NewLogger() *log.Logger {
  return ml.c.NewLogger()
}
