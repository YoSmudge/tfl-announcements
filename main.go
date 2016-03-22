package main

import(
  "github.com/samarudge/homecontrol-tubestatus/fetcher"
  log "github.com/Sirupsen/logrus"
  "github.com/voxelbrain/goptions"
)

type options struct {
  Verbose   bool            `goptions:"-v, --verbose, description='Log verbosely'"`
  Once      bool            `goptions:"-o, --once, description='Run once then exit'"`
  Speak     bool            `goptions:"-s, --speak, description='Speak the output'"`
  Help      goptions.Help   `goptions:"-h, --help, description='Show help'"`
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

  fetcher.RunStatus(parsedOptions.Once)
}
