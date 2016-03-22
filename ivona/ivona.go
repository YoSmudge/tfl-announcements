package ivona

import(
  iv "github.com/jpadilla/ivona-go"
  log "github.com/Sirupsen/logrus"
)

type Ivona struct{
  Client  *iv.Ivona
  Voice   *iv.Voice
}

func New(Access string, Secret string) *Ivona{
  i := Ivona{}

  i.Client = iv.New(Access, Secret)

  voices, _ := i.Client.ListVoices(iv.Voice{})
  for n,v := range voices.Voices{
    if v.Name == "Brian" && v.Language == "en-GB"{
      i.Voice = &voices.Voices[n]
    }
  }

  return &i
}

func (i *Ivona) GetSpeak(text string) ([]byte, error){
  o := iv.NewSpeechOptions(text)
  o.Voice = i.Voice

  r, err := i.Client.CreateSpeech(o)
  if err != nil{
    log.WithFields(log.Fields{
      "error": err,
      "text": text,
    }).Error("Error talking to Ivona Cloud")
    return []byte{}, err
  }

  return r.Audio, nil
}
