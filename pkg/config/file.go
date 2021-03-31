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

// Server represents a server definition in the configuration file.
type Server struct {
	Hostname      string `json:"hostname"`
	Protocol      string `json:"protocol"`
	Port          string `json:"port"`
	KibanaVersion string `json:"kibana_version"`
	BasicAuth     string `json:"basic_auth"`
}

// GetPort gets server port. If not provided it defaults to 443 to https and 80 to http protocols.
func (s *Server) GetPort() string {
	if len(s.Port) > 0 {
		return s.Port
	}
	if s.Protocol == "https" {
		return "443"
	}
	return "80"
}

// Role represents a role definition in the configuration file.
type Role struct {
	Index        string `json:"index"`
	WindowFilter string `json:"window_filter"`
}

// Location represents the configuration file location.
type Location struct {
	path string
	name string
}

// ConfigurationFile represents the root structure for the configuration file.
type ConfigurationFile struct {
	stdin         io.Reader
	location      Location
	Servers       map[string]Server `json:"servers"`
	Roles         map[string]Role   `json:"roles"`
	CurrentServer string            `json:"default_server"`
	CurrentRole   string            `json:"default_role"`
}

// Configuration contract to manage the configuration file.
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

// GetLocation gets the path where config file is located.
func (c *ConfigurationFile) GetLocation() Location {
	return c.location
}

// GetCurrentServer gets the current server definition.
func (c *ConfigurationFile) GetCurrentServer() Server {
	return c.Servers[c.CurrentServer]
}

// GetServer gets the current server name.
func (c *ConfigurationFile) GetServer() string {
	return c.CurrentServer
}

// FindServer finds a server by name.
func (c *ConfigurationFile) FindServer(name string) (Server, bool) {
	server, ok := c.Servers[name]
	return server, ok
}

// SetServer sets the default server value.
func (c *ConfigurationFile) SetServer(name string) {
	c.CurrentServer = name
}

// AddServer adds a server definition in the config file.
func (c *ConfigurationFile) AddServer(name string, server Server) {
	c.Servers[name] = server
}

// GetCurrentRole gets the current role definition.
func (c *ConfigurationFile) GetCurrentRole() Role {
	return c.Roles[c.CurrentRole]
}

// GetRole gets the current role name.
func (c *ConfigurationFile) GetRole() string {
	return c.CurrentRole
}

// FindRole finds a role by name.
func (c *ConfigurationFile) FindRole(name string) (Role, bool) {
	role, ok := c.Roles[name]
	return role, ok
}

// SetRole sets the default role.
func (c *ConfigurationFile) SetRole(name string) {
	c.CurrentRole = name
}

// AddRole adds role in the config file.
func (c *ConfigurationFile) AddRole(name string, role Role) {
	c.Roles[name] = role
}

// PrettyPrint prints the config file contents prettier.
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

// Save saves the config file in the files system in a JSON format.
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

// CheckEmpty checks if servers or roles definitions are empty. Returns an error if TRUE.
func (c *ConfigurationFile) CheckEmpty() error {
	if c.Servers == nil || len(c.Servers) <= 0 || c.Roles == nil || len(c.Roles) <= 0 {
		return errors.New("kishell is not configured. Use configure option before searching")
	}
	return nil
}

// Reset resets config file to empty values
func (c *ConfigurationFile) Reset() error {
	c.Servers = make(map[string]Server)
	c.Roles = make(map[string]Role)
	c.CurrentServer = ""
	c.CurrentRole = ""
	return c.Save()
}

// GetStdin gets where to read input from.
func (c *ConfigurationFile) GetStdin() io.Reader {
	return c.stdin
}

func homeDir() string {
	homeDir, err := homedir.Dir()
	checkError(err)
	return homeDir
}

// LoadDefaultConfig loads configuration from the default file config.
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
