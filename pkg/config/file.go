package config

import (
  "bytes"
  "encoding/json"
  "errors"
  "fmt"
  "github.com/mitchellh/go-homedir"
  "io/ioutil"
  "os"
)


const (
  configFileName = "/.kishell"
)

type Server struct {
  Hostname string `json:"hostname"`
  Protocol string `json:"protocol"`
  Port string `json:"port"`
  KibanaVersion string `json:"kibana_version"`
  BasicAuth string `json:"basic_auth"`
}

func (s *Server) GetPort() string {
  if len(s.Port) > 0 {
    return s.Port
  }
  if s.Protocol == "https" {
    return "443"
  }
  return "80"
}

type Role struct {
  Index string `json:"index"`
  WindowFilter string `json:"window_filter"`
}

type ConfigurationFile struct {
  Servers map[string]Server `json:"servers"`
  Roles map[string]Role `json:"roles"`
  CurrentServer string `json:"default_server"`
  CurrentRole string `json:"default_role"`
}

func checkError(err error) {
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func (c *ConfigurationFile) PrettyPrint() error {
  content, err := json.Marshal(c)
  if err != nil {
    return err
  }
  var prettyJSON bytes.Buffer
  err = json.Indent(&prettyJSON, content, "", "  ")
  if err != nil {
    return err
  }
  fmt.Println(string(prettyJSON.Bytes()))
  return nil
}

func (c *ConfigurationFile) Save() error {
  content, err := json.Marshal(c)
  if err != nil {
    return err
  }
  err = ioutil.WriteFile(homeDir() + configFileName, content, 0600)
  if err != nil {
    return err
  }
  return nil
}

func (c *ConfigurationFile) IsEmpty() error {
  if c.Servers == nil || len(c.Servers) <= 0 || c.Roles == nil || len(c.Roles) <= 0 {
    return errors.New("kishell is not configured. Use configure option before searching")
  }
  return nil
}

func (c *ConfigurationFile) Reset() error {
  configFile := ConfigurationFile {
    Servers: map[string]Server{},
    Roles: map[string]Role{},
  }
  return configFile.Save()
}

func homeDir() string {
  homeDir, err := homedir.Dir()
  checkError(err)
  return homeDir
}

func Load() ConfigurationFile {
  jsonFile, err := os.Open(homeDir() + configFileName)
  if err != nil && os.IsNotExist(err) {
    return ConfigurationFile{
      Servers: map[string]Server{},
      Roles: map[string]Role{},
    }
  } else {
    checkError(err)
  }
  var configFile ConfigurationFile
  err = json.NewDecoder(jsonFile).Decode(&configFile)
  checkError(err)
  return configFile
}
