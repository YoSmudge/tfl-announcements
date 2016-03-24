package fetcher

import(
  "github.com/samarudge/tfl-announcements/tfl"
  "github.com/samarudge/tfl-announcements/ivona"
  "github.com/samarudge/tfl-announcements/config"
  "time"
  "math"
  log "github.com/Sirupsen/logrus"
  "github.com/tuxychandru/pubsub"
  "github.com/satori/go.uuid"
)

type StatusUpdate struct{
  Id          uuid.UUID
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

func doStatus(a *tfl.Api, feed *pubsub.PubSub, isFull bool) {
  start := time.Now()
  statusUpdate, err := tfl.FetchStatus(a)
  if err != nil{
    log.WithFields(log.Fields{
      "error": err,
    }).Error("Error encountered getting statuses")
    return
  }

  statusText, err := statusUpdate.Generate(isFull)
  if err != nil{
    log.WithFields(log.Fields{
      "error": err,
    }).Error("Error encountered generating status text")
    return
  }

  if statusText != "" {
    i, err := ivona.New(config.Config.Ivona.Key, config.Config.Ivona.Secret)
    if err != nil{
      log.WithFields(log.Fields{
        "error": err,
      }).Error("Error creating Ivona connection")
      return
    }

    audio, err := i.GetSpeak(statusText)
    if err != nil{
      log.WithFields(log.Fields{
        "error": err,
      }).Error("Error talking to Ivona")
      return
    }

    duration := time.Since(start)

    u := StatusUpdate{}
    u.Id = uuid.NewV1()
    u.Created = time.Now()
    u.IsFull = true
    u.Text = statusText
    u.Duration = duration
    u.Audio = audio

    feed.Pub(&u, "updates")
  }
}

type nextRun struct{
  Time    time.Time
  IsFull  bool
}

func NextPeriod(period time.Duration) time.Time{
  unixStamp := float64(time.Now().Unix())
  secondsPerPeriod := period.Seconds()

  basePeriods := math.Floor(unixStamp/secondsPerPeriod)

  return time.Unix(int64(secondsPerPeriod*(basePeriods+1)), 0)
}

func NextRun() nextRun{
  p := nextRun{}
  p.Time = NextPeriod(config.Config.UpdatePeriod.Short)

  if p.Time == NextPeriod(config.Config.UpdatePeriod.Full){
    p.IsFull = true
  }

  return p
}

func RunStatus(runOnce bool, feed *pubsub.PubSub){
  a := tfl.NewApi(config.Config.Tfl.AppId, config.Config.Tfl.AppKey)
  np := NextRun()

  for {
    doStatus(a, feed, np.IsFull)

    if runOnce{
      <- time.After(5 * time.Second)
      return
    }

    np = NextRun()
    waitTime := time.Since(np.Time)*-1
    log.WithFields(log.Fields{
      "nextRun": np.Time,
      "wait": waitTime,
      "full": np.IsFull,
    }).Debug("Awaiting next run")
    <- time.After(waitTime)
  }
}
