package config

import(
  "time"
)

type tfl struct{
  AppId   string
  AppKey  string
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
  Tfl   tfl
  Ivona ivona
  UpdatePeriod updatePeriod
}

// These should be loaded dynamically
var Config = config{
  Tfl: tfl{
    AppId:  "7813aa7e",
    AppKey: "fe928850b2f8a836ad1f6ffcf4768549",
  },
  Ivona: ivona{
    Key:    "GDNAIRN6TKHGQVKHI6CQ",
    Secret: "jZnhDUwO5NuTj5THDjfYzR4KD99+fNuM+HBNIEoS",
  },
  UpdatePeriod: updatePeriod{
    Short:  time.Minute*10,
    Full:   time.Hour,
  },
}
