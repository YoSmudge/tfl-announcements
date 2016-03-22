package main

import(
  "fmt"
  "github.com/samarudge/homecontrol-tubestatus/tfl"
  log "github.com/Sirupsen/logrus"
)

const appId = "7813aa7e"
const apiKey = "fe928850b2f8a836ad1f6ffcf4768549"

func main() {
  log.SetLevel(log.DebugLevel)

  a := tfl.Api{appId, apiKey}
  statuses, err := tfl.FetchStatus(&a)
  if err != nil{
    log.WithFields(log.Fields{
      "error": err,
    }).Error("Error encountered getting statuses")
    return
  }

  /*
  statuses := tfl.StatusList{}
  statuses.Statuses = append(statuses.Statuses, tfl.Status{tfl.Line{"central","Central","tube"}, true, 9})
  */

  statusUpdate := tfl.StatusUpdate{statuses}

  fmt.Println(statusUpdate.Generate())
}
