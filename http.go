package floridalist

import(
  "log"
  "context"
  "net"
  "net/http"
  "strings"
  "strconv"

	"github.com/gorilla/mux"
)

type HttpController struct {
  config   Config
  member   *Member
}
func NewHttpController(config Config, member *Member) *HttpController {
  c := new(HttpController)
  c.config = config
  c.member = member
  return c
}
func (c *HttpController) HttpHandler() http.Handler {
  r := mux.NewRouter()
  r.StrictSlash(true)

  r.HandleFunc(c.config.FloridaEndpoint, c.GetSeeds).Methods("GET")
  r.HandleFunc("/_version", c.Version).Methods("GET")
  r.HandleFunc("/_chk", c.CheckStatus).Methods("GET")
  r.HandleFunc("/", c.CheckStatus).Methods("GET")

  return r
}
func (c *HttpController) GetSeeds(res http.ResponseWriter, req *http.Request) {
  if c.member.IsJoined() != true {
    res.Header().Set("Content-Type", "text/plain")
    res.WriteHeader(http.StatusNotAcceptable)
    res.Write([]byte("cluster not joined"))
    return
  }

  sl := c.member.SeedList()
  v  := strings.Join(sl, "|")
  res.Header().Set("Content-Type", "application/json")
  res.WriteHeader(http.StatusOK)
  res.Write([]byte(v))
}
func (c *HttpController) CheckStatus(res http.ResponseWriter, req *http.Request) {
  res.Header().Set("Content-Type", "text/plain")
  res.WriteHeader(http.StatusOK)

  v := strings.Join([]string{"OK", "\n"}, "")
  res.Write([]byte(v))
}
func (c *HttpController) Version(res http.ResponseWriter, req *http.Request) {
  res.Header().Set("Content-Type", "text/plain")
  res.WriteHeader(http.StatusOK)

  v := strings.Join([]string{UA, "\n"}, "")
  res.Write([]byte(v))
}

type WrapWriter struct {
  Writer      http.ResponseWriter
  LastStatus  int
}
func (w *WrapWriter) Header() http.Header {
  return w.Writer.Header()
}
func (w *WrapWriter) Write(b []byte) (int, error) {
  return w.Writer.Write(b)
}
func (w *WrapWriter) WriteHeader(status int) {
  w.LastStatus = status
  w.Writer.WriteHeader(status)
}
func WrapAccessLog(next http.Handler, logger HttpLogger) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    wrap := new(WrapWriter)
    wrap.Writer = w

    next.ServeHTTP(wrap, r)

    logger.Write(
      r.Host,
      r.Method,
      r.RequestURI,
      wrap.LastStatus,
      r.Header.Get("User-Agent"),
    )
  })
}

type HttpServer struct {
  config     Config
  Server     *http.Server
  Controller *HttpController
}
func NewHttpServer(ctx context.Context, m *Member) *HttpServer {
  config := ctx.Value("config").(Config)
  logger := ctx.Value("logger.http").(HttpLogger)

  ctr := NewHttpController(config, m)

  svr := new(HttpServer)
  svr.config = config
  svr.Controller = ctr
  svr.Server = &http.Server {
    Handler: WrapAccessLog(ctr.HttpHandler(), logger),
    ReadTimeout: config.HttpReadTimeout,
    WriteTimeout: config.HttpWriteTimeout,
  }
  return svr
}
func (s *HttpServer) Start(sctx context.Context) error {
  config      := s.config
  listenAddr  := net.JoinHostPort(config.FloridaBindIP, strconv.Itoa(config.FloridaBindPort))
  log.Printf("info: florida starting %s", listenAddr)

  listener, err := net.Listen("tcp4", listenAddr)
  if err != nil {
    log.Printf("error: addr '%s' listen error: %s", listenAddr, err.Error())
    return err
  }

  if err := s.Server.Serve(listener); err != nil && err != http.ErrServerClosed {
    log.Printf("error: serv error: %s", err.Error())
    return err
  }
  return nil
}
func (s *HttpServer) Stop(sctx context.Context) error {
  log.Printf("info: florida stoping")
  if err := s.Server.Shutdown(sctx); err != nil {
    log.Printf("error: shutdown error: %s", err.Error())
    return err
  }
  return nil
}
