package main

import(
  "time"
  "github.com/samarudge/homecontrol-tubestatus/tfl"
  log "github.com/Sirupsen/logrus"
  "github.com/voxelbrain/goptions"
)


type options struct {
  Verbose   bool            `goptions:"-v, --verbose, description='Log verbosely'"`
  Once      bool            `goptions:"-o, --once, description='Run once then exit'"`
  Help      goptions.Help   `goptions:"-h, --help, description='Show help'"`
}

const tflAppId = "7813aa7e"
const tflApiKey = "fe928850b2f8a836ad1f6ffcf4768549"

func doStatus(a *tfl.Api) {
  start := time.Now()
  statusUpdate, err := tfl.FetchStatus(a)
  if err != nil{
    log.WithFields(log.Fields{
      "error": err,
    }).Error("Error encountered getting statuses")
    return
  }

  statusTextShort, err := statusUpdate.Generate(true)

  if err != nil{
    log.WithFields(log.Fields{
      "error": err,
    }).Error("Error encountered generating short status text")
    return
  }

  statusTextFull, err := statusUpdate.Generate(false)

  if err != nil{
    log.WithFields(log.Fields{
      "error": err,
    }).Error("Error encountered generating full status text")
    return
  }

  duration := time.Since(start)

  log.WithFields(log.Fields{
    "duration": duration,
    "type": "short",
  }).Info(statusTextShort)

  log.WithFields(log.Fields{
    "duration": duration,
    "type": "full",
  }).Info(statusTextFull)
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

  a := tfl.NewApi(tflAppId, tflApiKey)
  doStatus(a)

  if parsedOptions.Once{
    return
  }

  statusTicker := time.NewTicker(30 * time.Second)
  statusEnd := make(chan struct{})

  for {
    select{
      case <- statusTicker.C:
        doStatus(a)
      case <- statusEnd:
        statusTicker.Stop()
        return
    }
  }
}
