package options

import (
  "github.com/sidilabs/kishell/pkg/config"
  "github.com/stretchr/testify/mock"
  "io"
)

type ConfigurationMock struct{
  mock.Mock
}

func (c *ConfigurationMock) GetLocation() config.Location {
  args := c.Called()
  return args.Get(0).(config.Location)
}

func (c *ConfigurationMock) GetCurrentServer() config.Server {
  args := c.Called()
  return args.Get(0).(config.Server)
}

func (c *ConfigurationMock) GetServer() string {
  args := c.Called()
  return args.String(0)
}

func (c *ConfigurationMock) FindServer(name string) (config.Server, bool) {
  args := c.Called(name)
  return args.Get(0).(config.Server), args.Bool(1)
}

func (c *ConfigurationMock) SetServer(name string) {
  c.Called(name)
}

func (c *ConfigurationMock) AddServer(name string, server config.Server) {
  c.Called(name, server)
}

func (c *ConfigurationMock) GetCurrentRole() config.Role {
  args := c.Called()
  return args.Get(0).(config.Role)
}

func (c *ConfigurationMock) GetRole() string {
  args := c.Called()
  return args.String(0)
}

func (c *ConfigurationMock) FindRole(name string) (config.Role, bool) {
  args := c.Called(name)
  return args.Get(0).(config.Role), args.Bool(1)
}

func (c *ConfigurationMock) SetRole(name string) {
  c.Called(name)
}

func (c *ConfigurationMock) AddRole(name string, role config.Role) {
  c.Called(name, role)
}

func (c *ConfigurationMock) Save() error {
  args := c.Called()
  return args.Error(0)
}

func (c *ConfigurationMock) CheckEmpty() error {
  args := c.Called()
  return args.Error(0)
}

func (c *ConfigurationMock) Reset() error {
  args := c.Called()
  return args.Error(0)
}

func (c *ConfigurationMock) PrettyPrint() error {
  args := c.Called()
  return args.Error(0)
}

func (c *ConfigurationMock) GetStdin() io.Reader {
  args := c.Called()
  return args.Get(0).(io.Reader)
}
