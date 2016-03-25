package tfl

import(
  "fmt"
  "strings"
  "sort"
  "github.com/samarudge/tfl-announcements/language"
)

type StatusUpdate struct{
  Api       *Api
  Statuses  StatusList
}

type statusText []string

func LineName(lineName string, mode string) string{
  lineModes := language.LineModes()
  if lineModes[mode]{
    lineName = fmt.Sprintf("%s Line", lineName)
  }

  return lineName
}

func LineAliases(line Line) []string{
  aliases := []string{line.Name}

  ca := language.GetRawParam("config")["line_aliases"].(map[interface{}]interface{})[line.Id]

  if ca != nil{
    for _,alias := range(ca.([]interface{})){
      aliases = append(aliases, alias.(string))
    }
  }

  return aliases
}

func (s *StatusUpdate) CoerceStatusUpdate(status Status) string{
  statusMsg := status.StatusDetails
  lineAliases := LineAliases(status.Line)

  for _,alias := range lineAliases{
    linePrefix := strings.ToLower(fmt.Sprintf("%s:", LineName(alias, status.Line.Mode)))

    if strings.HasPrefix(strings.ToLower(statusMsg), linePrefix){
      statusMsg = strings.TrimLeft(statusMsg[len(linePrefix):len(statusMsg)], " ")
    }
  }

  severities, _ := s.Api.GetSeverities()
  fmt.Println(strings.Join(severities, ","))
  fmt.Println("--->", status.Line.Name, statusMsg)
  severities = append(severities, language.AdditionalSeverities()...)
  fixedMsg := statusMsg
  for _,sv := range severities{
    statusPrefix := strings.ToLower(sv)
    if strings.HasPrefix(strings.ToLower(statusMsg), statusPrefix){
      trimmedMsg := strings.TrimLeft(statusMsg[len(statusPrefix):len(statusMsg)], " ")
      if len(trimmedMsg) < len(fixedMsg){
        fixedMsg = trimmedMsg
      }
    }
  }
  statusMsg = fixedMsg
  fmt.Println(statusMsg)
  fmt.Println("---")

  statusMsg = strings.Split(statusMsg, ".")[0]

  return strings.TrimRight(statusMsg, "., ")
}

func (s *StatusUpdate) Generate(fullUpdate bool) (string, error){
  var statusDetails statusText

  goodServiceModes := []string{}
  for _,mode := range s.Statuses.GoodServiceModes(){
    goodServiceModes = append(goodServiceModes, language.GetString("modes", mode))
  }

  sort.Strings(goodServiceModes)

  if s.Statuses.HasDisruption(){
    for _,status := range s.Statuses.DisruptedLines(){
      statusDescription, err := s.Api.GetSeverityFromCode(status.Line.Mode, status.LineStatus)
      if err != nil{
        return "", err
      }

      entireDetailsString := "entire_line"
      partDetailsString := "line_details"
      if language.IsPrefix(statusDescription){
        entireDetailsString = fmt.Sprintf("%s_is", entireDetailsString)
        partDetailsString = fmt.Sprintf("%s_is", partDetailsString)
      }

      var lineDetails string
      if status.WholeLine{
        lineDetails = language.RenderString("strings", entireDetailsString, language.H{
          "line_name": LineName(status.Line.Name, status.Line.Mode),
          "line_status": statusDescription,
        })
      } else {
        lineDetails = language.RenderString("strings", partDetailsString, language.H{
          "line_name": LineName(status.Line.Name, status.Line.Mode),
          "line_status": statusDescription,
        })
      }

      statusDetails = statusDetails.Add(lineDetails)

      if fullUpdate{
        statusDetails = statusDetails.Add(language.RenderString("strings", "due_to", language.H{
          "reason": s.CoerceStatusUpdate(status),
        }))
      }
    }

    if fullUpdate{
      statusDetails = statusDetails.Add(language.RenderString("strings", "other_good", language.H{
        "good_modes": makeList(goodServiceModes),
      }))
    } else {
      statusDetails = statusDetails.Add(language.GetString("strings", "other_good_short"))
    }
  } else {
    if fullUpdate{
      statusDetails = statusDetails.Add(language.RenderString("strings", "all_good", language.H{
        "good_modes": makeList(goodServiceModes),
      }))
    }
  }

  if len(statusDetails) > 0{
    statusDetails = append([]string{language.GetString("strings", "prefix")}, statusDetails...)
    return statusDetails.Format(), nil
  } else {
    return "", nil
  }
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
