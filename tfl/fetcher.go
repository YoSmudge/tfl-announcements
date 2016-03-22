package tfl

import(
)

func FetchStatus(api *Api) (StatusList, error){
  var statuses StatusList
  var lines LineList
  var err error

  lines, err = api.GetLines(modeTypes)
  if err != nil{
    return statuses, err
  }

  statuses, err = api.GetLineStatuses(lines)
  if err != nil{
    return statuses, err
  }

  return statuses, nil
}
