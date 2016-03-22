package afplay

import(
  log "github.com/Sirupsen/logrus"
  "github.com/samarudge/homecontrol-tubestatus/fetcher"
  "io/ioutil"
  "os"
  "os/exec"
)

func Afplay(u *fetcher.StatusUpdate){
  tmpFile, err := ioutil.TempFile("", "tfl-afplay")
  if err != nil{
    log.WithFields(log.Fields{
      "err": err,
    }).Warning("Could not create temp file for afplay")
    return
  }
  tmpFilePath := tmpFile.Name()
  defer tmpFile.Close()

  _, err = tmpFile.Write(u.Audio)
  if err != nil{
    log.WithFields(log.Fields{
      "err": err,
      "file": tmpFilePath,
    }).Warning("Could not write temp file for afplay")
    return
  }

  log.WithFields(log.Fields{
    "tmpFile": tmpFilePath,
  }).Debug("Playing with AFPlay")

  playCmd := exec.Command("afplay", tmpFilePath)
  err = playCmd.Start()
  if err != nil{
    log.WithFields(log.Fields{
      "err": err,
      "file": tmpFilePath,
    }).Warning("Could not run afplay")
  }

  err = playCmd.Wait()
  if err != nil{
    log.WithFields(log.Fields{
      "err": err,
      "file": tmpFilePath,
    }).Warning("Non-zero exit for afplay")
  }

  if os.Remove(tmpFilePath) != nil{
    log.WithFields(log.Fields{
      "err": err,
      "file": tmpFilePath,
    }).Warning("Could not remove temp file for afplay")
  }
}
