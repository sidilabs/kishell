package config

import (
  "bytes"
  "encoding/json"
  "errors"
  "fmt"
  "github.com/mitchellh/go-homedir"
  "io"
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

type Location struct {
  path string
  name string
}

type ConfigurationFile struct {
  stdin io.Reader
  location Location
  Servers map[string]Server `json:"servers"`
  Roles map[string]Role `json:"roles"`
  CurrentServer string `json:"default_server"`
  CurrentRole string `json:"default_role"`
}

type Configuration interface {
  GetLocation() Location
  GetCurrentServer() Server
  GetServer() string
  FindServer(name string) (Server, bool)
  SetServer(name string)
  AddServer(name string, server Server)
  GetCurrentRole() Role
  GetRole() string
  FindRole(name string) (Role, bool)
  SetRole(name string)
  AddRole(name string, role Role)
  PrettyPrint() error
  Save() error
  CheckEmpty() error
  Reset() error
  GetStdin() io.Reader
}


func (l *Location) get() string {
  return l.path + l.name
}

func checkError(err error) {
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func (c *ConfigurationFile) GetLocation() Location {
  return c.location
}

func (c *ConfigurationFile) GetCurrentServer() Server {
  return c.Servers[c.CurrentServer]
}

func (c *ConfigurationFile) GetServer() string {
  return c.CurrentServer
}

func (c *ConfigurationFile) FindServer(name string) (Server, bool) {
  server, ok := c.Servers[name]
  return server, ok
}

func (c *ConfigurationFile) SetServer(name string) {
  c.CurrentServer = name
}

func (c *ConfigurationFile) AddServer(name string, server Server) {
  c.Servers[name] = server
}

func (c *ConfigurationFile) GetCurrentRole() Role {
  return c.Roles[c.CurrentRole]
}

func (c *ConfigurationFile) GetRole() string {
  return c.CurrentRole
}

func (c *ConfigurationFile) FindRole(name string) (Role, bool) {
  role, ok := c.Roles[name]
  return role, ok
}

func (c *ConfigurationFile) SetRole(name string) {
  c.CurrentRole = name
}

func (c *ConfigurationFile) AddRole(name string, role Role) {
  c.Roles[name] = role
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
  err = ioutil.WriteFile(c.location.get(), content, 0600)
  if err != nil {
    return err
  }
  return nil
}

func (c *ConfigurationFile) CheckEmpty() error {
  if c.Servers == nil || len(c.Servers) <= 0 || c.Roles == nil || len(c.Roles) <= 0 {
    return errors.New("kishell is not configured. Use configure option before searching")
  }
  return nil
}

func (c *ConfigurationFile) Reset() error {
  c.Servers = make(map[string]Server)
  c.Roles = make(map[string]Role)
  c.CurrentServer = ""
  c.CurrentRole = ""
  return c.Save()
}

func (c *ConfigurationFile) GetStdin() io.Reader {
  return c.stdin
}

func homeDir() string {
  homeDir, err := homedir.Dir()
  checkError(err)
  return homeDir
}

func LoadDefaultConfig() Configuration {
  return loadConfig(homeDir(), configFileName)
}

func loadConfig(path string, fileName string) Configuration {
  jsonFile, err := os.Open(path + fileName)
  if err != nil && os.IsNotExist(err) {
    return &ConfigurationFile {
      stdin: os.Stdin,
      location: Location {
        path: path,
        name: fileName,
      },
      Servers: map[string]Server{},
      Roles: map[string]Role{},
    }
  } else {
    checkError(err)
  }
  var configFile ConfigurationFile
  err = json.NewDecoder(jsonFile).Decode(&configFile)
  checkError(err)
  configFile.location = Location {
    path: path,
    name: fileName,
  }
  return &configFile
}
