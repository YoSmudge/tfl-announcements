package tfl

import(
  //"fmt"
)

type StatusList struct{
  Statuses      []Status
}

type Status struct{
  Line          Line
  Issues        bool
  StatusLevel   float64
}

func (s *Status) Text() string{
  return "Moo"
}

func ParseStatusStruct(status *Status, statusStruct interface{}){
  statusResponses := statusStruct.([]interface{})
  statusArrays := statusResponses[0].(map[string]interface{})["lineStatuses"].([]interface{})
  var worstStatus float64 = 15

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

    if statusLevel < worstStatus && validNow {
      worstStatus = statusLevel
      status.Issues = true
    }
  }

  status.StatusLevel = worstStatus
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
    gsm[st.Line.Mode] = true
  }

  modes := []string{}
  for mode,_ := range gsm{
    modes = append(modes, mode)
  }

  return modes
}
