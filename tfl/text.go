package tfl

import(
  "fmt"
  "strings"
  "github.com/samarudge/homecontrol-tubestatus/language"
)

type StatusUpdate struct{
  Api       *Api
  Statuses  StatusList
}

type statusText []string

func (s *StatusUpdate) Generate() (string, error){
  var statusDetails statusText

  statusDetails = statusDetails.Add(language.GetString("strings", "prefix"))
  lineModes := language.LineModes()

  if s.Statuses.HasDisruption(){
    for _,status := range s.Statuses.DisruptedLines(){
      statusDescription, err := s.Api.GetSeverityFromCode(status.Line.Mode, status.StatusLevel)
      if err != nil{
        return "", err
      }

      lineName := status.Line.Name
      if lineModes[status.Line.Mode]{
        lineName = fmt.Sprintf("%s Line", status.Line.Name)
      }

      lineDetails := language.RenderString("strings", "line_status", language.H{
        "line_name": lineName,
        "line_status": statusDescription,
      })

      statusDetails = statusDetails.Add(lineDetails)
    }

    goodServiceModes := []string{}
    for _,mode := range s.Statuses.GoodServiceModes(){
      goodServiceModes = append(goodServiceModes, language.GetString("modes", mode))
    }

    statusDetails = statusDetails.Add(language.RenderString("strings", "other_good", language.H{
      "good_modes": makeList(goodServiceModes),
    }))
  } else {
    statusDetails = statusDetails.Add(language.GetString("strings", "all_good"))
  }

  return statusDetails.Format(), nil
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

    if i < wordsLen-1{
      outputs = append(outputs, ", ")
    }
  }

  return strings.Join(outputs, "")
}
