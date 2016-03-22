package language

import(
  "fmt"
  "gopkg.in/yaml.v2"
  "github.com/hoisie/mustache"
)

type H map[string]interface{}

func GetRawParam(stringGroup string) map[string]interface{}{
  sf, err := Asset(fmt.Sprintf("language/%s.yml", stringGroup))
  if err != nil{
    panic(fmt.Sprintf("Tried to load language file %s but was not found", stringGroup))
  }

  data := make(map[string]interface{})
  err = yaml.Unmarshal(sf, &data)

  if err != nil{
    panic(fmt.Sprintf("Tried to load language file %s but was invalid", stringGroup))
  }

  return data
}

func GetString(stringGroup string, stringName string) string{
  stringMap := GetRawParam(stringGroup)

  pickedString := stringMap[stringName]
  if pickedString == ""{
    panic(fmt.Sprintf("String %s not found in language file %s", stringName, stringGroup))
  }

  return pickedString.(string)
}

func RenderString(stringGroup string, stringName string, ctx interface{}) string{
  baseTemplate := GetString(stringGroup, stringName)
  return mustache.Render(baseTemplate, ctx)
}

func LineModes() map[string]bool{
  rsp := make(map[string]bool)

  for _,mode := range GetRawParam("config")["line_modes"].([]interface{}){
    rsp[mode.(string)] = true
  }

  return rsp
}
