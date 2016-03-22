package tfl

import(
  "fmt"
  "strings"
  "github.com/samarudge/homecontrol-tubestatus/language"
)

type StatusUpdate struct{
  Statuses  StatusList
}

type statusText []string

func (s *StatusUpdate) Generate() string{
  var statusDetails statusText

  statusDetails = statusDetails.Add(language.GetString("strings", "prefix"))

  if s.Statuses.HasDisruption(){
    for _,status := range s.Statuses.DisruptedLines(){
      lineDetails := fmt.Sprintf("The %s Line has %f", status.Line.Name, status.StatusLevel)

      statusDetails = statusDetails.Add(lineDetails)
    }

    goodServiceModes := []string{}
    for _,mode := range s.Statuses.GoodServiceModes(){
      goodServiceModes = append(goodServiceModes, language.GetString("modes", mode))
    }

    statusDetails = statusDetails.Add(language.GetString("strings", "other_good"))
    statusDetails = statusDetails.Add(makeList(goodServiceModes))
  } else {
    statusDetails = statusDetails.Add(language.GetString("strings", "all_good"))
  }

  return statusDetails.Format()
}

func (t statusText) Add(msg string) statusText{
  t = append(t, msg)
  return t
}

func (t statusText) Format() string{
  return strings.Join(t, " ")
}


func makeList(l []string) string{
  wordsLen := len(l)
  outputs := []string{}
  for i,word := range l{
    if i == wordsLen-1{
      outputs = append(outputs, "and ")
    }

    outputs = append(outputs, word)

    if i < wordsLen{
      outputs = append(outputs, ", ")
    }
  }

  return strings.Join(outputs, "")
}
