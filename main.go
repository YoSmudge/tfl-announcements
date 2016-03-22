package main

import(
  "github.com/samarudge/homecontrol-tubestatus/fetcher"
  "github.com/samarudge/homecontrol-tubestatus/afplay"
  log "github.com/Sirupsen/logrus"
  "github.com/voxelbrain/goptions"
  "github.com/tuxychandru/pubsub"
)

type options struct {
  Verbose   bool            `goptions:"-v, --verbose, description='Log verbosely'"`
  Once      bool            `goptions:"-o, --once, description='Run once then exit'"`
  Afplay    bool            `goptions:"--afplay, description='Play audio with OS/X afplay command'"`
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

  goptions.ParseAndFail(&parsedOptions)

  if parsedOptions.Verbose{
    log.SetLevel(log.DebugLevel)
  } else {
    log.SetLevel(log.InfoLevel)
  }

  log.SetFormatter(&log.TextFormatter{FullTimestamp:true})

  messageFeed := pubsub.New(1)
  go fetcher.SubscribeHandler(messageFeed, logOutput)

  if parsedOptions.Afplay{
    go fetcher.SubscribeHandler(messageFeed, afplay.Afplay)
  }

  fetcher.RunStatus(parsedOptions.Once, messageFeed)
}
