package web

import(
  "fmt"
  "github.com/gorilla/websocket"
  "time"
  log "github.com/Sirupsen/logrus"
  "net/http"
  "net"
  "sync/atomic"
  "github.com/samarudge/tfl-announcements/fetcher"
  "github.com/tuxychandru/pubsub"
  "runtime"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 4096,
}

const pingInterval = time.Second*45
var connectedClients int64 = 0
var clientIdCounter uint64 = 0
var statusFeed = pubsub.New(1)

type socketClient struct{
  ClientId      uint64
  RemoteAddr    net.Addr
  Connection    *websocket.Conn
  Subscription  chan interface{}
}

type socketMsg struct{
  Type              string      `json:"message_type"`
  Body              interface{} `json:"body"`
  ClientId          uint64      `json:"client_id"`
  Time              time.Time   `json:"time"`
  PingInt           int64       `json:"ping_interval"`
  ConnectedClients  int64       `json:"connected_clients"`
  GrCount           int         `json:"goroutine_count"`
  MemUsage          uint64      `json:"memory_usage"`
  HeapObjects       uint64      `json:"heap_objects"`
}

type clientUpdate struct{
  Created     time.Time   `json:"created"`
  IsFull      bool        `json:"is_full"`
  Text        string      `json:"text"`
  UpdateUrl   string      `json:"audio_url"`
}

func formatUpdate(u *fetcher.StatusUpdate) clientUpdate{
  cu := clientUpdate{}
  cu.Created = u.Created
  cu.IsFull = u.IsFull
  cu.Text = u.Text
  cu.UpdateUrl = fmt.Sprintf("/UpdateAudio?update_id=%s", u.Id)

  return cu
}

func (c *socketClient) NewMessage(mt string, b interface{}) *socketMsg{
  m := socketMsg{}
  m.ClientId = c.ClientId
  m.Type = mt
  m.Time = time.Now().UTC()
  m.Body = b
  m.PingInt = int64(pingInterval.Seconds())
  m.ConnectedClients = connectedClients

  m.GrCount = runtime.NumGoroutine()
  s := runtime.MemStats{}
  runtime.ReadMemStats(&s)
  m.HeapObjects = s.HeapObjects
  m.MemUsage = s.Alloc

  return &m
}

func (c *socketClient) SendMessage(mt string, b interface{}){
  msg := c.NewMessage(mt, b)

  log.WithFields(log.Fields{
    "clientId": c.ClientId,
    "type": mt,
    "body": b,
  }).Debug("Sending message")

  err := c.Connection.WriteJSON(msg)
  if err != nil {
    log.WithFields(log.Fields{
      "address": c.RemoteAddr,
      "id": c.ClientId,
      "error": err,
      "type": mt,
      "body": b,
    }).Debug("Error sending message to client")
  }
}

func ping(c socketClient, close chan struct{}) {
  ticker := time.NewTicker(pingInterval)
  defer ticker.Stop()
  for {
    select{
      case <-ticker.C:
        c.SendMessage("ping", "")
      case <-close:
        return
    }
  }
}

func pushStatus(c socketClient, close chan struct{}){
  for {
    select{
      case u := <-c.Subscription:
        log.WithFields(log.Fields{
          "id": c.ClientId,
        }).Debug("Sending status update")
        c.SendMessage("status_update", formatUpdate(u.(*fetcher.StatusUpdate)))
      case <-close:
        return
    }
  }
}

func feedHandler(w http.ResponseWriter, r *http.Request) {
  ws, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.WithFields(log.Fields{
      "err": err,
    }).Warning("Failed to set websocket upgrade")
    return
  }

  defer ws.Close()

  c := socketClient{}
  c.RemoteAddr = ws.RemoteAddr()
  c.Connection = ws
  c.ClientId = atomic.AddUint64(&clientIdCounter, 1)
  c.Subscription = statusFeed.Sub("push_update")

  log.WithFields(log.Fields{
    "address": c.RemoteAddr,
    "id": c.ClientId,
  }).Debug("New socket connection")

  atomic.AddInt64(&connectedClients, 1)

  done := make(chan struct{})

  var heloMsg interface{}
  if currentStatus != nil{
    heloMsg = formatUpdate(currentStatus)
  }
  c.SendMessage("helo", heloMsg)

  go ping(c, done)
  go pushStatus(c, done)

  for{
    msgType, msg, err := ws.ReadMessage()
    if err != nil{
      if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
        log.WithFields(log.Fields{
          "id": c.ClientId,
        }).Debug("Client disconnect")
      } else {
        log.WithFields(log.Fields{
          "err": err,
        }).Warning("Error reading message")
      }
      break
    }

    log.WithFields(log.Fields{
      "type": msgType,
      "msg": msg,
    }).Debug("Got message from client")
  }

  log.WithFields(log.Fields{
    "id": c.ClientId,
  }).Debug("Closing client connection")

  close(done)
  atomic.AddInt64(&connectedClients, -1)
  statusFeed.Unsub(c.Subscription)
}

func SendUpdate(u *fetcher.StatusUpdate){
  currentStatus = u
  statusFeed.Pub(u, "push_update")
}
