package tfl

import(
  "strings"
  "html"
)

type LineList struct{
  lines   []Line
}

type Line struct{
  Id      string
  Name    string
  Mode    string  `json:"modeName"`
}

func (l *LineList) Count() int{
  return len(l.lines)
}

func (l *LineList) List() string{
  lineIds := []string{}

  for _,line := range l.lines{
    lineIds = append(lineIds, line.Id)
  }

  return strings.Join(lineIds, ",")
}

func (l *LineList) FixNames(){
  for _,line := range l.lines{
    line.Name = html.UnescapeString(line.Name)
  }
}
