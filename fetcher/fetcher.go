package fetcher

import(
  "github.com/samarudge/homecontrol-tubestatus/tfl"
  "github.com/samarudge/homecontrol-tubestatus/ivona"
  "github.com/samarudge/homecontrol-tubestatus/config"
  "time"
  log "github.com/Sirupsen/logrus"
)

func doStatus(a *tfl.Api) {
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

  i := ivona.New(config.IvonaKey, config.IvonaSecret)
  err = i.Speak(statusTextFull)

  if err != nil{
    log.WithFields(log.Fields{
      "error": err,
    }).Warning("Error talking to Ivona")
  }
}

func RunStatus(runOnce bool){
  a := tfl.NewApi(config.TflAppId, config.TflApiKey)

  doStatus(a)

  if runOnce{
    return
  }

  statusTicker := time.NewTicker(10 * time.Minute)
  statusEnd := make(chan struct{})

  for {
    select{
      case <- statusTicker.C:
        go doStatus(a)
      case <- statusEnd:
        statusTicker.Stop()
        return
    }
  }
}
