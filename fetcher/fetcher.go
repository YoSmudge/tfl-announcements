package fetcher

import(
  "github.com/samarudge/homecontrol-tubestatus/tfl"
  "github.com/samarudge/homecontrol-tubestatus/ivona"
  "github.com/samarudge/homecontrol-tubestatus/config"
  "time"
  log "github.com/Sirupsen/logrus"
  "github.com/tuxychandru/pubsub"
)

type StatusUpdate struct{
  Created     time.Time
  IsFull      bool
  Text        string
  Duration    time.Duration
  Audio       []byte
}

func SubscribeHandler(feed *pubsub.PubSub, handler func(u *StatusUpdate)){
  c := feed.Sub("updates")
  for uc := range c{
    u := uc.(*StatusUpdate)
    handler(u)
  }
}

func doStatus(a *tfl.Api, feed *pubsub.PubSub) {
  start := time.Now()
  statusUpdate, err := tfl.FetchStatus(a)
  if err != nil{
    log.WithFields(log.Fields{
      "error": err,
    }).Error("Error encountered getting statuses")
    return
  }

  statusText, err := statusUpdate.Generate(false)
  if err != nil{
    log.WithFields(log.Fields{
      "error": err,
    }).Error("Error encountered generating status text")
    return
  }

  i := ivona.New(config.IvonaKey, config.IvonaSecret)
  audio, err := i.GetSpeak(statusText)
  if err != nil{
    log.WithFields(log.Fields{
      "error": err,
    }).Warning("Error talking to Ivona")
  }

  duration := time.Since(start)

  u := StatusUpdate{}
  u.Created = time.Now()
  u.IsFull = true
  u.Text = statusText
  u.Duration = duration
  u.Audio = audio

  feed.Pub(&u, "updates")
}

func RunStatus(runOnce bool, feed *pubsub.PubSub){
  a := tfl.NewApi(config.TflAppId, config.TflApiKey)

  doStatus(a, feed)

  if runOnce{
    return
  }

  statusTicker := time.NewTicker(10 * time.Minute)
  statusEnd := make(chan struct{})

  for {
    select{
      case <- statusTicker.C:
        doStatus(a, feed)
      case <- statusEnd:
        statusTicker.Stop()
        return
    }
  }
}
