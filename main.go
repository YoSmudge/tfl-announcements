package main

import(
  "github.com/samarudge/homecontrol-tubestatus/fetcher"
  log "github.com/Sirupsen/logrus"
  "github.com/voxelbrain/goptions"
  "github.com/tuxychandru/pubsub"
)

type options struct {
  Verbose   bool            `goptions:"-v, --verbose, description='Log verbosely'"`
  Once      bool            `goptions:"-o, --once, description='Run once then exit'"`
  Speak     bool            `goptions:"-s, --speak, description='Speak the output'"`
  Help      goptions.Help   `goptions:"-h, --help, description='Show help'"`
}

func logOutput(c chan interface{}){
  for us := range c{
    u := us.(fetcher.StatusUpdate)
    log.WithFields(log.Fields{
      "full":u.IsFull,
      "duration":u.Duration,
      "created":u.Created,
    }).Info(u.Text)
  }
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
  go logOutput(messageFeed.Sub("updates"))

  fetcher.RunStatus(parsedOptions.Once, messageFeed)
}
