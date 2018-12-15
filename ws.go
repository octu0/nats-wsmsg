package wsmsg

import(
  "log"
  "time"
  "context"
  "bytes"
  "net/http"
  "encoding/gob"

  "github.com/gorilla/websocket"
  "github.com/nats-io/go-nats"
)

const maxMessageSize = 1024 * 1024
var upgrader = websocket.Upgrader{
  ReadBufferSize:  10 * maxMessageSize,
  WriteBufferSize: 10 * maxMessageSize,
}

type Message struct {
  MsgType  int
  Data     []byte
}

func CreateNatsClient(ctx context.Context) (*nats.Conn, error) {
  connUrl := ctx.Value("nats.url").(string)
  nc, err := nats.Connect(connUrl,
    nats.DontRandomize(),
    nats.NoEcho(),
    nats.Name(UA),
    nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error){
      log.Printf("error: %v", err)
    }),
  )
  if err != nil {
    log.Printf("error: failed to connection: %v", err)
    return nil, err
  }
  return nc, nil
}
func DecodeMessage(nc *nats.Conn, m *nats.Msg) (*Message, error) {
  data := m.Data
  msg  := new(Message)
  dec  := gob.NewDecoder(bytes.NewBuffer(data))
  if err := dec.Decode(&msg); err != nil {
    return nil, err
  }
  return msg, nil
}
func EncodeMessage(nc *nats.Conn, msg *Message) ([]byte, error) {
  b   := new(bytes.Buffer)
  enc := gob.NewEncoder(b)
  if err := enc.Encode(msg); err != nil {
    return nil, err
  }
  return b.Bytes(), nil
}
func publish(nc *nats.Conn, topic string, msg *Message) error {
  b, err := EncodeMessage(nc, msg)
  if err != nil {
    return err
  }
  nc.Publish(topic, b)
  nc.Flush()
  return nil
}
func PublishText(nc *nats.Conn, topic string, data []byte) error {
  return publish(nc, topic, &Message{
    MsgType: websocket.TextMessage,
    Data:    data,
  })
}
func PublishBinary(nc *nats.Conn, topic string, data []byte) error {
  return publish(nc, topic, &Message{
    MsgType: websocket.BinaryMessage,
    Data:    data,
  })
}

type SendQueue chan *Message
type SubQueue  chan *nats.Msg
type WebsocketHandler struct {
  conn         *websocket.Conn
  nc           *nats.Conn
  send         SendQueue
  subs         SubQueue
  subscription *nats.Subscription
  running      bool
}
func (ws *WebsocketHandler) readLoop() {
  ws.conn.SetReadLimit(maxMessageSize)
  ws.conn.SetReadDeadline(time.Now().Add(10 * time.Second))
  ws.conn.SetPongHandler(func(string) error {
    ws.conn.SetReadDeadline(time.Now().Add(10 * time.Second));
    return nil
  })
  for ws.running {
    _, message, err := ws.conn.ReadMessage()
    if err != nil {
      if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
        log.Printf("error: %v", err)
      }
      break
    }
    log.Printf("debug: message = %s", message)
  }
}
func (ws *WebsocketHandler) writeLoop() {
  ticker := time.NewTicker(5 * time.Second)
  defer func() {
    ticker.Stop()
  }()

  for ws.running {
    select {
    case message, ok := <-ws.send:
      ws.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
      if !ok {
        ws.conn.WriteMessage(websocket.CloseMessage, []byte{})
        return
      }

      w, err := ws.conn.NextWriter(message.MsgType)
      if err != nil {
        return
      }
      w.Write(message.Data)

      if err := w.Close(); err != nil {
        return
      }
    case <-ticker.C:
      ws.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
      if err := ws.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
        return
      }
    }
  }
}
func (ws *WebsocketHandler) subExchangeLoop(){
  for ws.running {
    select {
    case m := <-ws.subs:
      msg, err := DecodeMessage(ws.nc, m);
      if err != nil {
        log.Printf("warn: gob decode failure: %s", err.Error())
        continue
      }
      ws.send <- msg
    }
  }
}
func (ws *WebsocketHandler) Close() error {
  ws.subscription.Unsubscribe()
  ws.nc.Drain()
  ws.nc.Close()
  ws.running = false
  close(ws.send)
  close(ws.subs)
  return ws.conn.Close()
}
func (ws *WebsocketHandler) RunSubscribe(topic string) error {
  sub, err := ws.nc.ChanSubscribe(topic, ws.subs)
  if err != nil {
    log.Printf("error: nats topic(%s) subscription failure: %s", topic, err.Error())
    ws.Close()
    return err
  }
  ws.nc.Flush()
  ws.subscription = sub

  ws.runloop()
  return nil
}
func (ws *WebsocketHandler) RunSubscribeWithGroup(topic, group string) error {
  sub, err := ws.nc.ChanQueueSubscribe(topic, group, ws.subs)
  if err != nil {
    log.Printf("error: nats topic(%s) subscription failure: %s", topic, err.Error())
    ws.Close()
    return err
  }
  ws.nc.Flush()
  ws.subscription = sub

  ws.runloop()
  return nil
}
func (ws *WebsocketHandler) runloop() {
  ws.running = true
  go ws.readLoop()
  go ws.writeLoop()
  go ws.subExchangeLoop()
}

func CreateWebsocketHandler(ctx context.Context, res http.ResponseWriter, req *http.Request) (*WebsocketHandler, error) {
  conn, err := upgrader.Upgrade(res, req, nil)
  if err != nil {
    log.Printf("warn: ws upgrade failure: %s", err.Error())
    return nil, err
  }

  nc, err := CreateNatsClient(ctx)
  if err != nil {
    log.Printf("error: nats client creation failure: %s", err.Error())
    return nil, err
  }

  ws := new(WebsocketHandler)
  ws.conn  = conn
  ws.nc    = nc
  ws.send  = make(SendQueue, 0)
  ws.subs  = make(SubQueue, 0)
  return ws, nil
}
