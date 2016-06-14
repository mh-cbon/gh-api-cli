package local

import (
  "encoding/json"
  "io/ioutil"
  "errors"
  "strings"
  "os"

	"github.com/google/go-github/github"
	"github.com/mitchellh/go-homedir"
)

func Add (name string, auth *github.Authorization) error {
  auths := Read()
  auths[name] = *auth
  return Save(auths)
}

func Remove (name string) error {
  auths := Read()
  delete(auths, name)
  return Save(auths)
}

func Save (auths map[string]github.Authorization) error {
  p, err := filePath()
  if err!=nil {
    return err
  }
  jsonContent, err := json.MarshalIndent(auths, "", "    ")
  if _, err := os.Stat(p); os.IsNotExist(err) {
    f, err2 := os.Create(p)
    f.Close()
    if err2!=nil {
      return err
    }
  }
  return ioutil.WriteFile(p, jsonContent, 0744)
}

func Read () map[string]github.Authorization {
  ret := make(map[string]github.Authorization)
  p, err := filePath()
  if err!=nil {
    return ret
  }
  if _, err := os.Stat(p); os.IsNotExist(err) {
    if err!=nil {
      return ret
    }
  }
  data, err := ioutil.ReadFile(p)
  if err!=nil {
    return ret
  }
  if len(data)>0 {
    err = json.Unmarshal(data, &ret);
  }
  return ret
}

func Get (name string) (github.Authorization, error) {
  ret := github.Authorization{}
  auths := Read ()
  if _, ok := auths[name]; ok==false {
    return ret, errors.New("Authorization '"+name+"' not found")
  }
  return auths[name], nil
}

func filePath () (string, error){
  home, err := homedir.Dir()
  if err!=nil {
    return "", err
  }
  filePath := strings.Join([]string{home, "gh-auths.json"}, string(os.PathSeparator))
  return filePath, nil
}
