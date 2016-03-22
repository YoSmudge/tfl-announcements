package main

import(
  "time"
  "github.com/samarudge/homecontrol-tubestatus/tfl"
  "github.com/samarudge/homecontrol-tubestatus/ivona"
  log "github.com/Sirupsen/logrus"
  "github.com/voxelbrain/goptions"
)


type options struct {
  Verbose   bool            `goptions:"-v, --verbose, description='Log verbosely'"`
  Once      bool            `goptions:"-o, --once, description='Run once then exit'"`
  Speak     bool            `goptions:"-s, --speak, description='Speak the output'"`
  Help      goptions.Help   `goptions:"-h, --help, description='Show help'"`
}

const tflAppId = "7813aa7e"
const tflApiKey = "fe928850b2f8a836ad1f6ffcf4768549"
const ivonaKey = "GDNAIRN6TKHGQVKHI6CQ"
const ivonaSecret = "jZnhDUwO5NuTj5THDjfYzR4KD99+fNuM+HBNIEoS"

func doStatus(a *tfl.Api, speak bool) {
  start := time.Now()
  statusUpdate, err := tfl.FetchStatus(a)
  if err != nil{
    log.WithFields(log.Fields{
      "error": err,
    }).Error("Error encountered getting statuses")
    return
  }

  statusTextShort, err := statusUpdate.Generate(false)

  if err != nil{
    log.WithFields(log.Fields{
      "error": err,
    }).Error("Error encountered generating short status text")
    return
  }

  statusTextFull, err := statusUpdate.Generate(true)

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

  if speak{
    i := ivona.New(ivonaKey, ivonaSecret)
    err := i.Speak(statusTextFull)

    if err != nil{
      log.WithFields(log.Fields{
        "error": err,
      }).Warning("Error speaking text")
    }
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

  a := tfl.NewApi(tflAppId, tflApiKey)
  go doStatus(a, parsedOptions.Speak)

  if parsedOptions.Once{
    return
  }

  statusTicker := time.NewTicker(10 * time.Minute)
  statusEnd := make(chan struct{})

  for {
    select{
      case <- statusTicker.C:
        go doStatus(a, parsedOptions.Speak)
      case <- statusEnd:
        statusTicker.Stop()
        return
    }
  }
}
