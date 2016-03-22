package config

import(
  "time"
  "os"
  log "github.com/Sirupsen/logrus"
  "gopkg.in/yaml.v2"
  "io/ioutil"
)

type tfl struct{
  AppId   string  `yaml:"app_id"`
  AppKey  string  `yaml:"app_key"`
}

type ivona struct{
  Key     string
  Secret  string
}

type updatePeriod struct{
  Short   time.Duration
  Full    time.Duration
}

type config struct{
  Tfl           tfl
  Ivona         ivona
  UpdatePeriod  updatePeriod   `yaml:"update_period"`
}

var Config config

func Load(filePath string){
  Config = config{}
  if _, err := os.Stat(filePath); os.IsNotExist(err) {
    log.WithFields(log.Fields{
      "configFile": filePath,
      "error": "File not found",
    }).Error("Could not load config file")
    os.Exit(1)
  }

  configContent, err := ioutil.ReadFile(filePath)
  if err != nil {
    log.WithFields(log.Fields{
      "configFile": filePath,
      "error": err,
    }).Error("Could not load config file")
    os.Exit(1)
  }

  err = yaml.Unmarshal(configContent, &Config)
  if err != nil {
    log.WithFields(log.Fields{
      "configFile": filePath,
      "error": err,
    }).Error("Could not load config file")
    os.Exit(1)
  }

  if int64(Config.UpdatePeriod.Full.Seconds())%int64(Config.UpdatePeriod.Short.Seconds()) != 0{
    log.WithFields(log.Fields{
      "configFile": filePath,
      "error": "Short update period must be a factor of full update period",
    }).Error("Could not load config file")
    os.Exit(1)
  }
}
