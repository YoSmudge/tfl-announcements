package language

import(
  "fmt"
  "gopkg.in/yaml.v2"
  "github.com/hoisie/mustache"
)

type H map[string]interface{}

func GetString(stringGroup string, stringName string) string{
  sf, err := Asset(fmt.Sprintf("language/%s.yml", stringGroup))
  if err != nil{
    panic(fmt.Sprintf("Tried to load language file %s but was not found", stringGroup))
  }

  stringMap := make(map[string]string)
  err = yaml.Unmarshal(sf, &stringMap)
  if err != nil{
    panic(fmt.Sprintf("Tried to load language file %s but was invalid", stringGroup))
  }

  pickedString := stringMap[stringName]
  if pickedString == ""{
    panic(fmt.Sprintf("String %s not found in language file %s", stringName, stringGroup))
  }

  return pickedString
}

func RenderString(stringGroup string, stringName string, ctx interface{}) string{
  baseTemplate := GetString(stringGroup, stringName)
  return mustache.Render(baseTemplate, ctx)
}
