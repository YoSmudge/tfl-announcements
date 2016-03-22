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
  statuses, err := tfl.FetchStatus(a)
  if err != nil{
    log.WithFields(log.Fields{
      "error": err,
    }).Error("Error encountered getting statuses")
    return
  }

  statusUpdate := tfl.StatusUpdate{statuses}
  statusText := statusUpdate.Generate()
  duration := time.Since(start)

  log.WithFields(log.Fields{
    "duration": duration,
  }).Info(statusText)
}

func main() {
  log.SetLevel(log.DebugLevel)

  a := tfl.Api{appId, apiKey}
  doStatus(&a)

  /*
  statuses := tfl.StatusList{}
  statuses.Statuses = append(statuses.Statuses, tfl.Status{tfl.Line{"central","Central","tube"}, true, 9})
  */

  statusTicker := time.NewTicker(30 * time.Second)
  statusEnd := make(chan struct{})

  for {
    select{
      case <- statusTicker.C:
        doStatus(&a)
      case <- statusEnd:
        statusTicker.Stop()
        return
    }
  }
}
