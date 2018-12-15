package wsmsg

import(
  "log"
  "context"
  "io"
  "net"
  "net/http"
  "strings"
  "strconv"
  "bufio"
  "errors"

  "github.com/gorilla/mux"
)

type HttpController struct {
  ctx context.Context
}
func NewHttpController(ctx context.Context) *HttpController {
  c := new(HttpController)
  c.ctx = ctx
  return c
}
func (c *HttpController) HttpHandler() http.Handler {
  r := mux.NewRouter()
  r.StrictSlash(true)

  ws := r.PathPrefix("/ws").Subrouter()
  ws.HandleFunc("/sub/{topic}", c.TopicSubscription)
  ws.HandleFunc("/pub/{topic}", c.TopicPublish).Methods("POST")
  ws.HandleFunc("/queue/{topic}", c.TopicEnqueue).Methods("POST")

  r.HandleFunc("/_version", c.Version).Methods("GET")
  r.HandleFunc("/_chk", c.CheckStatus).Methods("GET")
  r.HandleFunc("/", c.CheckStatus).Methods("GET")

  return r
}
func (c *HttpController) readBody(req *http.Request) ([]byte, error) {
  length, err := strconv.Atoi(req.Header.Get("Content-Length"))
  if err != nil {
    log.Printf("error: read content-length error %s", err.Error())
    return nil, err
  }

  body := make([]byte, length)
  length, err = req.Body.Read(body)
  if err != nil && err != io.EOF {
    log.Printf("error: read body error %s", err.Error())
    return nil, err
  }
  req.Body.Close()

  return body, nil
}
func (c *HttpController) writeln(text string, res http.ResponseWriter) {
  res.Write([]byte(strings.Join([]string{text, "\n"}, "")))
}
func (c *HttpController) na(res http.ResponseWriter, req *http.Request) {
  res.Header().Set("Content-Type", "application/json")
  res.WriteHeader(http.StatusNotAcceptable)

  c.writeln(`{"status":"not acceptable"}`, res)
}
func (c *HttpController) json(text string, res http.ResponseWriter, req *http.Request) {
  res.Header().Set("Content-Type", "application/json")
  res.WriteHeader(http.StatusOK)

  c.writeln(text, res)
}
func (c *HttpController) ok(text string, res http.ResponseWriter, req *http.Request) {
  res.Header().Set("Content-Type", "text/plain")
  res.WriteHeader(http.StatusOK)

  c.writeln(text, res)
}
func (c *HttpController) Version(res http.ResponseWriter, req *http.Request) {
  c.ok(UA, res, req)
}
func (c *HttpController) CheckStatus(res http.ResponseWriter, req *http.Request) {
  c.ok("OK", res, req)
}
func (c *HttpController) TopicSubscription(res http.ResponseWriter, req *http.Request) {
  vars  := mux.Vars(req)
  topic := vars["topic"]

  ws, err := CreateWebsocketHandler(c.ctx, res, req)
  if err != nil {
    c.na(res, req)
    return
  }
  ws.RunSubscribe(topic)
}
func (c *HttpController) TopicPublish(res http.ResponseWriter, req *http.Request) {
  vars  := mux.Vars(req)
  topic := vars["topic"]

  body, err := c.readBody(req);
  if err != nil {
    c.na(res, req)
    return
  }

  nc, err := CreateNatsClient(c.ctx)
  if err != nil {
    c.na(res, req)
    return
  }
  defer nc.Close()

  PublishText(nc, topic, body)
  c.json(`{"message":"success"}`, res, req)
}
func (c *HttpController) TopicEnqueue(res http.ResponseWriter, req *http.Request) {
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
// allow hijack
func (w *WrapWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
  hijacker, ok := w.Writer.(http.Hijacker)
  if !ok {
		return nil, nil, errors.New("WrapWriter doesn't support the Hijacker interface")
	}
  return hijacker.Hijack()
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
func NewHttpServer(ctx context.Context) *HttpServer {
  config := ctx.Value("config").(Config)
  logger := ctx.Value("logger.http").(HttpLogger)

  ctr := NewHttpController(ctx)

  svr := new(HttpServer)
  svr.config = config
  svr.Controller = ctr
  svr.Server = &http.Server {
    Handler: WrapAccessLog(ctr.HttpHandler(), logger),
  }
  return svr
}
func (s *HttpServer) Start(sctx context.Context) error {
  config      := s.config
  listenAddr  := net.JoinHostPort(config.BindIP, strconv.Itoa(config.BindPort))
  log.Printf("info: http server starting %s", listenAddr)

  listener, err := net.Listen("tcp", listenAddr)
  if err != nil {
    log.Printf("error: addr '%s' listen error: %s", listenAddr, err.Error())
    return err
  }

  if err := s.Server.Serve(listener); err != nil && err != http.ErrServerClosed {
    log.Printf("error: http server serv error: %s", err.Error())
    return err
  }
  return nil
}
func (s *HttpServer) Stop(sctx context.Context) error {
  log.Printf("info: http server stoping")
  if err := s.Server.Shutdown(sctx); err != nil {
    log.Printf("error: shutdown error: %s", err.Error())
    return err
  }
  return nil
}
