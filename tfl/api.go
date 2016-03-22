package tfl

import(
  "fmt"
  log "github.com/Sirupsen/logrus"
  "net/http"
  "encoding/json"
  "net/url"
  "io/ioutil"
  "strings"
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

  rsp, err := http.Get(callPath.String())
  if err != nil{
    log.WithFields(log.Fields{
      "Error": err,
      "Status": rsp.StatusCode,
      "Target": logUrl,
    }).Error("Error calling TFL")
    return err
  }

  log.WithFields(log.Fields{
    "Target": logUrl,
    "Status": rsp.StatusCode,
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

func (a *Api) GetLineStatuses(lines LineList) (StatusList, error){
  var statusList StatusList
  for _,line := range lines.lines{
    st, err := a.GetLineStatus(line)
    if err != nil{
      log.WithFields(log.Fields{
        "line": line.Id,
        "error": err,
      }).Error("Could not get line status")
    } else {
      statusList.Statuses = append(statusList.Statuses, st)
    }
  }

  return statusList, nil
}
