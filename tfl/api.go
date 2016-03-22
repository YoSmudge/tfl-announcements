package tfl

import(
  "fmt"
  log "github.com/Sirupsen/logrus"
  "net/http"
  "encoding/json"
  "net/url"
  "io/ioutil"
  "strings"
  "time"
)

var modeTypes = []string{"tube", "overground", "dlr", "tflrail"}

type Api struct{
  AppId   string
  ApiKey  string
}

func (a *Api) DoCall(call string, reciever interface{}) error{
  callPath, _ := url.Parse("https://api.tfl.gov.uk/")
  parsedUrl, _ := url.Parse(call)
  callPath.Path = parsedUrl.Path
  logUrl := callPath.String()
  q := parsedUrl.Query()
  q.Set("app_id", a.AppId)
  q.Set("app_key", a.ApiKey)
  callPath.RawQuery = q.Encode()

  callStart := time.Now()
  rsp, err := http.Get(callPath.String())
  if err != nil{
    log.WithFields(log.Fields{
      "Error": err,
      "Status": rsp.StatusCode,
      "Target": logUrl,
    }).Error("Error calling TFL")
    return err
  }

  callDuration := time.Since(callStart)
  log.WithFields(log.Fields{
    "Target": logUrl,
    "Status": rsp.StatusCode,
    "Duration": callDuration,
  }).Debug("Call to TFL")

  defer rsp.Body.Close()
  responseRaw, _ := ioutil.ReadAll(rsp.Body)

  if err := json.Unmarshal(responseRaw, &reciever); err != nil {
    log.WithFields(log.Fields{
      "Error": err,
      "Status": rsp.StatusCode,
      "Target": callPath.String(),
    }).Error("Error decoding JSON response")
    return fmt.Errorf("Could not decode JSON", string(responseRaw))
  }

  return nil
}

func (a *Api) GetLines(lineModes []string) (LineList, error){
  lines := LineList{}

  err := a.DoCall(fmt.Sprintf("/Line/Mode/%s", strings.Join(lineModes, ",")), &lines.lines)
  if err != nil{
    return lines, err
  }
  log.WithFields(log.Fields{
    "lineCount": lines.Count(),
    "lineList": lines.List(),
  }).Debug("Found lines")

  return lines, nil
}

func (a *Api) GetLineStatus(line Line) (Status, error){
  status := Status{}
  var statusResponse interface{}
  err := a.DoCall(fmt.Sprintf("/Line/%s/Status", line.Id), &statusResponse)
  if err != nil{
    return status, err
  }

  status.Line = line
  ParseStatusStruct(&status, statusResponse)

  log.WithFields(log.Fields{
    "issues": status.Issues,
    "status": status.StatusLevel,
    "line": line.Id,
  }).Debug("Got line status")

  return status, nil
}

type asyncLineStatus struct{
  status    Status
  err       error
}

func (a *Api) AsyncGetLineStatus(line Line, response chan asyncLineStatus){
  st, err := a.GetLineStatus(line)
  response <- asyncLineStatus{st,err}
}

func (a *Api) GetLineStatuses(lines LineList) (StatusList, error){
  var statusList StatusList
  responses := make(chan asyncLineStatus)

  for _,line := range lines.lines{
    go a.AsyncGetLineStatus(line,responses)
  }

  for i := 1; i <= len(lines.lines); i++{
    rsp := <-responses
    if rsp.err != nil{
      log.Error(rsp.err)
    } else {
      statusList.Statuses = append(statusList.Statuses, rsp.status)
    }
  }

  return statusList, nil
}
