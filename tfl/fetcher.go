package tfl

import(
)

func FetchStatus(api *Api) (StatusUpdate, error){
  var statuses StatusList
  var lines LineList
  var err error

  statusUpdate := StatusUpdate{}

  lines, err = api.GetLines(modeTypes)
  if err != nil{
    return statusUpdate, err
  }

  statuses, err = api.GetLineStatuses(lines)
  if err != nil{
    return statusUpdate, err
  }

  statusUpdate.Api = api
  statusUpdate.Statuses = statuses

  return statusUpdate, nil
}
