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

// Struct for configuration file servers.
type Server struct {
	Hostname      string `json:"hostname"`
	Protocol      string `json:"protocol"`
	Port          string `json:"port"`
	KibanaVersion string `json:"kibana_version"`
	BasicAuth     string `json:"basic_auth"`
}

// Gets server port. If not provided it defaults to 443 to https and 80 to http protocols.
func (s *Server) GetPort() string {
	if len(s.Port) > 0 {
		return s.Port
	}
	if s.Protocol == "https" {
		return "443"
	}
	return "80"
}

// Struct for configuration file roles.
type Role struct {
	Index        string `json:"index"`
	WindowFilter string `json:"window_filter"`
}

// Struct for configuration file location.
type Location struct {
	path string
	name string
}

// Configuration file representation.
type ConfigurationFile struct {
	stdin         io.Reader
	location      Location
	Servers       map[string]Server `json:"servers"`
	Roles         map[string]Role   `json:"roles"`
	CurrentServer string            `json:"default_server"`
	CurrentRole   string            `json:"default_role"`
}

// Contract to manage the configuration file.
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

// Gets the path where config file is located.
func (c *ConfigurationFile) GetLocation() Location {
	return c.location
}

// Gets the current server definition.
func (c *ConfigurationFile) GetCurrentServer() Server {
	return c.Servers[c.CurrentServer]
}

// Gets the current server name.
func (c *ConfigurationFile) GetServer() string {
	return c.CurrentServer
}

// Finds a server by name.
func (c *ConfigurationFile) FindServer(name string) (Server, bool) {
	server, ok := c.Servers[name]
	return server, ok
}

// Sets the default server value.
func (c *ConfigurationFile) SetServer(name string) {
	c.CurrentServer = name
}

// Adds a server definition in the config file.
func (c *ConfigurationFile) AddServer(name string, server Server) {
	c.Servers[name] = server
}

// Gets the current tolr definition.
func (c *ConfigurationFile) GetCurrentRole() Role {
	return c.Roles[c.CurrentRole]
}

// Gets the current role name.
func (c *ConfigurationFile) GetRole() string {
	return c.CurrentRole
}

// Finds a role by name.
func (c *ConfigurationFile) FindRole(name string) (Role, bool) {
	role, ok := c.Roles[name]
	return role, ok
}

// Sets the default role.
func (c *ConfigurationFile) SetRole(name string) {
	c.CurrentRole = name
}

// Adds role in the config file.
func (c *ConfigurationFile) AddRole(name string, role Role) {
	c.Roles[name] = role
}

// Pretty prints the config file contents.
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

// Saves the config file in the files system in a JSON format.
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

// Checks if servers or roles definitions are empty. Returns an error if TRUE.
func (c *ConfigurationFile) CheckEmpty() error {
	if c.Servers == nil || len(c.Servers) <= 0 || c.Roles == nil || len(c.Roles) <= 0 {
		return errors.New("kishell is not configured. Use configure option before searching")
	}
	return nil
}

// Resets config file to empty values
func (c *ConfigurationFile) Reset() error {
	c.Servers = make(map[string]Server)
	c.Roles = make(map[string]Role)
	c.CurrentServer = ""
	c.CurrentRole = ""
	return c.Save()
}

// Where to read input from
func (c *ConfigurationFile) GetStdin() io.Reader {
	return c.stdin
}

func homeDir() string {
	homeDir, err := homedir.Dir()
	checkError(err)
	return homeDir
}

// Loads configuration from the default file config.
func LoadDefaultConfig() Configuration {
	return loadConfig(homeDir(), configFileName)
}

func loadConfig(path string, fileName string) Configuration {
	jsonFile, err := os.Open(path + fileName)
	if err != nil && os.IsNotExist(err) {
		return &ConfigurationFile{
			stdin: os.Stdin,
			location: Location{
				path: path,
				name: fileName,
			},
			Servers: map[string]Server{},
			Roles:   map[string]Role{},
		}
	}
	var configFile ConfigurationFile
	err = json.NewDecoder(jsonFile).Decode(&configFile)
	checkError(err)
	configFile.location = Location{
		path: path,
		name: fileName,
	}
	return &configFile
}
