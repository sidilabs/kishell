package options

import (
  "bufio"
  "encoding/base64"
  "errors"
  "fmt"
  "github.com/sidilabs/kishell/pkg/config"
  "strings"
)

const (
  lineBreak = "\n"
  lineBreakAsByte = '\n'
)

func (c *ConfigureCmd) Run(ctx *Context) error {
  if c.Server {
    addServer(ctx.Configuration)
    return ctx.Configuration.Save()
  } else if c.Role {
    addRole(ctx.Configuration)
    return ctx.Configuration.Save()
  } else if c.Reset {
    return ctx.Configuration.Reset()
  }
  return errors.New("missing parameter. One of the following is expected: --server | --role")
}

func addServer(configuration config.Configuration) {
  reader := bufio.NewReader(configuration.GetStdin())
  fmt.Print("Server name: ")
  serverName, _ := reader.ReadString(lineBreakAsByte)
  serverName = strings.TrimSuffix(serverName, lineBreak)
  configuration.AddServer(serverName, buildServer(reader))
  fmt.Print("Set as default? [Y/n]: ")
  defaultServer, _ := reader.ReadString(lineBreakAsByte)
  defaultServer = strings.TrimSuffix(defaultServer, lineBreak)
  if len(configuration.GetServer()) <=  0 || len(defaultServer) <= 0 || (defaultServer == "Y" || defaultServer == "y") {
    configuration.SetServer(serverName)
  }
}

func buildServer(reader *bufio.Reader) config.Server {
  fmt.Print("Protocol: ")
  protocol, _ := reader.ReadString(lineBreakAsByte)
  fmt.Print("Hostname: ")
  hostname, _ := reader.ReadString(lineBreakAsByte)
  fmt.Print("Port: ")
  port, _ := reader.ReadString(lineBreakAsByte)
  fmt.Print("Username: ")
  username, _ := reader.ReadString(lineBreakAsByte)
  fmt.Print("Password: ")
  password, _ := reader.ReadString(lineBreakAsByte)
  fmt.Print("Kibana Version: ")
  kibanaVersion, _ := reader.ReadString(lineBreakAsByte)
  basicAuth := strings.TrimSuffix(username, lineBreak) + ":" + strings.TrimSuffix(password, lineBreak)
  if len(basicAuth) <= 1 {
    basicAuth = ""
  }
  server := config.Server {
    Protocol:  strings.TrimSuffix(protocol, lineBreak),
    Hostname:  strings.TrimSuffix(hostname, lineBreak),
    Port:      strings.TrimSuffix(port, lineBreak),
    BasicAuth: base64.StdEncoding.EncodeToString([]byte(basicAuth)),
    KibanaVersion: strings.TrimSuffix(kibanaVersion, lineBreak),
  }
  return server
}

func addRole(configuration config.Configuration)  {
  reader := bufio.NewReader(configuration.GetStdin())
  fmt.Print("Role name: ")
  roleName, _ := reader.ReadString(lineBreakAsByte)
  roleName = strings.TrimSuffix(roleName, lineBreak)
  configuration.AddRole(roleName, buildRole(reader))
  fmt.Print("Set as default? [Y/n]: ")
  defaultRole, _ := reader.ReadString(lineBreakAsByte)
  defaultRole = strings.TrimSuffix(defaultRole, lineBreak)
  if len(configuration.GetRole()) <=  0 || len(defaultRole) <= 0 || (defaultRole == "Y" || defaultRole == "y") {
    configuration.SetRole(roleName)
  }
}

func buildRole(reader *bufio.Reader) config.Role {
  fmt.Print("Index name: ")
  index, _ := reader.ReadString(lineBreakAsByte)
  fmt.Print("Window filter time (e.g. @timestamp, modified_date): ")
  windowFilter, _ := reader.ReadString(lineBreakAsByte)
  role := config.Role {
    Index: strings.TrimSuffix(index, lineBreak),
    WindowFilter: strings.TrimSuffix(windowFilter, lineBreak),
  }
  return role
}
