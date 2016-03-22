package ivona

import(
  iv "github.com/jpadilla/ivona-go"
  log "github.com/Sirupsen/logrus"
)

type Ivona struct{
  Client  *iv.Ivona
  Voice   *iv.Voice
}

func New(Access string, Secret string) (*Ivona, error){
  i := Ivona{}

  i.Client = iv.New(Access, Secret)

  voices, err := i.Client.ListVoices(iv.Voice{})
  if err != nil{
    return &i, err
  }

  for n,v := range voices.Voices{
    if v.Name == "Brian" && v.Language == "en-GB"{
      i.Voice = &voices.Voices[n]
    }
  }

  return &i, nil
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
