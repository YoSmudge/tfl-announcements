package tfl

import(
  "html"
)

type StatusList struct{
  Statuses      []Status
}

type Status struct{
  Line          Line
  Issues        bool
  WholeLine     bool
  LineStatus    float64
  StatusDetails string
}

func ParseStatusStruct(status *Status, statusStruct interface{}){
  statusResponses := statusStruct.([]interface{})

  for _,sr := range statusResponses{
    statusArrays := sr.(map[string]interface{})["lineStatuses"].([]interface{})

    for _,stIn := range statusArrays{
      statusLevel := stIn.(map[string]interface{})["statusSeverity"].(float64)
      validNow := false

      validityArray := stIn.(map[string]interface{})["validityPeriods"].([]interface{})
      for _,vp := range validityArray{
        validity := vp.(map[string]interface{})
        isNow := validity["isNow"].(bool)
        if isNow{
          validNow = true
        }
      }

      if statusLevel < status.LineStatus && validNow {
        status.Issues = true
      }

      disruptionDetails := stIn.(map[string]interface{})["disruption"]

      if disruptionDetails != nil {
        wholeLine := disruptionDetails.(map[string]interface{})["isWholeLine"]
        if wholeLine != nil && wholeLine.(bool) == true{
          status.Issues = true
          status.LineStatus = statusLevel
          status.WholeLine = true
          status.StatusDetails = html.UnescapeString(stIn.(map[string]interface{})["reason"].(string))
        }
      }
    }
  }
}

func (s *StatusList) DisruptedLines() []Status{
  disrupted := []Status{}

  for _,st := range s.Statuses{
    if st.Issues{
      disrupted = append(disrupted, st)
    }
  }

  return disrupted
}

func (s *StatusList) HasDisruption() bool{
  if len(s.DisruptedLines()) > 0{
    return true
  } else {
    return false
  }
}

func (s *StatusList) GoodServiceModes() []string{
  gsm := make(map[string]bool)
  for _, st := range s.Statuses{
    if !st.Issues{
      gsm[st.Line.Mode] = true
    }
  }

  modes := []string{}
  for mode,_ := range gsm{
    modes = append(modes, mode)
  }

  return modes
}
