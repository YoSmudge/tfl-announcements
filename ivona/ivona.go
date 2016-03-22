package ivona

import(
  iv "github.com/jpadilla/ivona-go"
  log "github.com/Sirupsen/logrus"
  "io/ioutil"
)

type Ivona struct{
  Client  *iv.Ivona
}

func New(Access string, Secret string) *Ivona{
  i := Ivona{}

  i.Client = iv.New(Access, Secret)

  return &i
}

func (i *Ivona) Speak(text string) error{
  o := iv.NewSpeechOptions(text)

  r, err := i.Client.CreateSpeech(o)
  if err != nil{
    log.WithFields(log.Fields{
      "error": err,
      "text": text,
    }).Error("Error talking to Ivona Cloud")
    return err
  }

  err = ioutil.WriteFile("./test.mp3", r.Audio, 0644)
  return nil
}
