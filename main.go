package main

import(
  "github.com/samarudge/tfl-announcements/fetcher"
  "github.com/samarudge/tfl-announcements/afplay"
  "github.com/samarudge/tfl-announcements/config"
  "github.com/samarudge/tfl-announcements/web"
  log "github.com/Sirupsen/logrus"
  "github.com/voxelbrain/goptions"
  "github.com/tuxychandru/pubsub"
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
