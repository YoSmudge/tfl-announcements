package main

import(
  "github.com/samarudge/tfl-announcements/fetcher"
  "github.com/samarudge/tfl-announcements/afplay"
  "github.com/samarudge/tfl-announcements/config"
  "github.com/samarudge/tfl-announcements/web"
  log "github.com/Sirupsen/logrus"
  "github.com/voxelbrain/goptions"
  "github.com/tuxychandru/pubsub"
  "time"
  "runtime"
  "fmt"
  "github.com/quipo/statsd"
)

type options struct {
  Config    string          `goptions:"-c, --config, description='Config Yaml file to use'"`
  Verbose   bool            `goptions:"-v, --verbose, description='Log verbosely'"`
  Once      bool            `goptions:"-o, --once, description='Run once then exit'"`
  Afplay    bool            `goptions:"--afplay, description='Play audio with OS/X afplay command'"`
  Web       bool            `goptions:"--web, description='Run web player'"`
  WebBind   string          `goptions:"--web-bind, description='Bind address for web player'"`
  Help      goptions.Help   `goptions:"-h, --help, description='Show help'"`
}

func logOutput(u *fetcher.StatusUpdate){
  log.WithFields(log.Fields{
    "full":u.IsFull,
    "duration":u.Duration,
    "created":u.Created,
  }).Info(u.Text)
}

func main() {
  parsedOptions := options{}
  parsedOptions.Config = "./config.yml"
  parsedOptions.WebBind = ":8001"

  goptions.ParseAndFail(&parsedOptions)

  if parsedOptions.Verbose{
    log.SetLevel(log.DebugLevel)
  } else {
    log.SetLevel(log.InfoLevel)
  }

  log.SetFormatter(&log.TextFormatter{FullTimestamp:true})

  config.Load(parsedOptions.Config)

  go reportStats()

  messageFeed := pubsub.New(1)
  go fetcher.SubscribeHandler(messageFeed, logOutput)

  if parsedOptions.Afplay{
    go fetcher.SubscribeHandler(messageFeed, afplay.Afplay)
  }

  if parsedOptions.Web{
    go fetcher.SubscribeHandler(messageFeed, web.SendUpdate)
    go web.StartServer(parsedOptions.WebBind)
  }

  fetcher.RunStatus(parsedOptions.Once, messageFeed)
}

func reportStats(){
  var appName = "tfl_announcements"

  if config.Config.Statsd == ""{
    return
  }

  statsClient := statsd.NewStatsdClient(config.Config.Statsd, "")
  err := statsClient.CreateSocket()
  if err != nil{
    log.WithFields(log.Fields{
      "host": config.Config.Statsd,
    }).Warning("Could not connect to statsd")
    return
  }

  statsInterval := 10*time.Second
  statsSender := statsd.NewStatsdBuffer(statsInterval, statsClient)
  statsTicker := time.NewTicker(statsInterval)
  defer statsTicker.Stop()
  defer statsSender.Close()

  for {
    s := runtime.MemStats{}
    runtime.ReadMemStats(&s)
    metrics := map[string]int64{
      "goroutines.running": int64(runtime.NumGoroutine()),
      "memory.objects.HeapObjects": int64(s.HeapObjects),
      "memory.Alloc": int64(s.Alloc),
      "web.connections": int64(web.ConnectedClients),
    }

    for k,v := range metrics{
      mName := fmt.Sprintf("%s.%s", appName, k)
      statsSender.Gauge(mName, v)
    }

    <-statsTicker.C
  }
}
