package main

import(
  "time"
  "github.com/samarudge/homecontrol-tubestatus/tfl"
  log "github.com/Sirupsen/logrus"
)

const appId = "7813aa7e"
const apiKey = "fe928850b2f8a836ad1f6ffcf4768549"

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
  log.SetLevel(log.DebugLevel)

  a := tfl.NewApi(appId, apiKey)
  doStatus(a)

  /*
  statuses := tfl.StatusList{}
  statuses.Statuses = append(statuses.Statuses, tfl.Status{tfl.Line{"central","Central","tube"}, true, 9})
  */

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
