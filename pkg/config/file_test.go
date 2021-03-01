package config

import (
  "bytes"
  "io/ioutil"
  "os"
  "testing"
)

const (
  testConfigPath = "../../testdata"
  testConfigFile = "/kishell-config.json"
)

func TestLoadConfigFile(t *testing.T) {
  file := loadConfig(testConfigPath, testConfigFile)
  _ = file.PrettyPrint()
  _ = file.CheckEmpty()
  if file.GetRole() != "local" {
    t.Errorf("Invalid current role %s", file.GetRole())
  }
  if file.GetServer() != "local" {
    t.Errorf("Invalid current server %s", file.GetServer())
  }
  role := file.GetCurrentRole()
  if role.Index != "local-*" {
    t.Errorf("Invalid index for role %s", role.Index)
  }
  role, _ = file.FindRole("another")
  if role.Index != "another-*" {
    t.Errorf("Invalid index for role %s", role.Index)
  }
  server := file.GetCurrentServer()
  if server.GetPort() != "8080" {
    t.Errorf("Invalid server port %s", server.GetPort())
  }
  server, _ = file.FindServer("httpServer")
  if server.GetPort() != "80" {
    t.Errorf("Invalid server port %s", server.GetPort())
  }
  server, _ = file.FindServer("httpsServer")
  if server.GetPort() != "443" {
    t.Errorf("Invalid server port %s", server.GetPort())
  }
}

func TestResetConfigFile(t *testing.T) {
  file, err := createTempConfig()
  if err != nil {
    t.Fatal("Unable to create temp config file", err)
  }
  _ = file.Reset()
  err = file.CheckEmpty()
  if err == nil {
    t.Fatal("Config file is supposed to be empty after resetting")
  }
  location := file.GetLocation()
  actual, _ := ioutil.ReadFile(location.get())
  expected, _ := ioutil.ReadFile("../../testdata/empty-kishell-config.golden")

  if !bytes.Equal(actual, expected) {
    t.Fatalf("Expected config to be %s but was %s", string(expected), string(actual))
  }

  defer os.RemoveAll(file.location.path)
}

func TestConfigFileEmpty(t *testing.T) {
  file := loadConfig(testConfigPath, "/file-does-not-exist.json")
  err := file.CheckEmpty()
  if err == nil {
    t.Fatal("Config file should be empty byt it is not", err)
  }
}

func TestAddSetRoleNServer(t *testing.T) {
  file := loadConfig(testConfigPath, testConfigFile)
  file.AddRole("new-role", Role{ Index: "new-index" })
  file.SetRole("new-role")
  if file.GetCurrentRole().Index != "new-index" {
    t.Fatal("Invalid role definition")
  }
  file.AddServer("new-server", Server{ Hostname: "new.index.is" })
  file.SetServer("new-server")
  if file.GetCurrentServer().Hostname != "new.index.is" {
    t.Fatal("Invalid server definition")
  }
}

func createTempConfig() (*ConfigurationFile, error) {
  dir, err := ioutil.TempDir("/tmp", "ut-")
  if err != nil {
    return nil, err
  }
  file := ConfigurationFile {
    location: Location {
      name: "/kishell-config",
      path: dir,
    },
    CurrentRole: "test",
    CurrentServer: "teest",
    Servers: map[string]Server {
      "test": {
        Hostname:      "test.local",
        Protocol:      "http",
        Port:          "443",
        KibanaVersion: "1.1.1",
        BasicAuth:     "dGVzdDpsb2NhbHBhc3N3ZAo=",
      },
    },
    Roles: map[string]Role {
      "test": {
        Index:        "test-*",
        WindowFilter: "@timestamp",
      },
    },
  }
  err = file.Save()
  if err != nil {
    return nil, err
  }
  return &file, nil
}
