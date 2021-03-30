package options

import (
	"github.com/sidilabs/kishell/pkg/config"
	"testing"
)

func TestUseServerOption(t *testing.T) {
	serverName := "ut-server"
	configuration := new(ConfigurationMock)
	configuration.On("FindServer", serverName).Return(config.Server{}, true)
	configuration.On("SetServer", serverName)
	configuration.On("Save").Return(nil)

	context := Context{
		Debug:         true,
		Configuration: configuration,
	}

	cmd := UseCmd{
		Server: serverName,
	}
	err := cmd.Run(&context)
	if err != nil {
		t.Fatal("Defining current server failed", err)
	}
	configuration.AssertExpectations(t)
}

func TestServerNotFound(t *testing.T) {
	serverName := "ut-server"
	configuration := new(ConfigurationMock)
	configuration.On("FindServer", serverName).Return(config.Server{}, false)

	context := Context{
		Debug:         true,
		Configuration: configuration,
	}

	cmd := UseCmd{
		Server: serverName,
	}
	err := cmd.Run(&context)
	if err == nil {
		t.Fatal("Defining current server failed", err)
	}
	configuration.AssertExpectations(t)
}

func TestUseRoleOption(t *testing.T) {
	roleName := "ut-role"
	configuration := new(ConfigurationMock)
	configuration.On("FindRole", roleName).Return(config.Role{}, true)
	configuration.On("SetRole", roleName)
	configuration.On("Save").Return(nil)

	context := Context{
		Debug:         true,
		Configuration: configuration,
	}

	cmd := UseCmd{
		Role: roleName,
	}
	err := cmd.Run(&context)
	if err != nil {
		t.Fatal("Defining current role failed", err)
	}
	configuration.AssertExpectations(t)
}

func TestRoleNotFound(t *testing.T) {
	serverName := "ut-server"
	configuration := new(ConfigurationMock)
	configuration.On("FindRole", serverName).Return(config.Role{}, false)

	context := Context{
		Debug:         true,
		Configuration: configuration,
	}

	cmd := UseCmd{
		Role: serverName,
	}
	err := cmd.Run(&context)
	if err == nil {
		t.Fatal("Defining current role failed", err)
	}
	configuration.AssertExpectations(t)
}

func TestInvalidUseOption(t *testing.T) {
	configuration := new(ConfigurationMock)
	context := Context{
		Debug:         true,
		Configuration: configuration,
	}
	cmd := UseCmd{}
	err := cmd.Run(&context)
	if err == nil {
		t.Fatal("Defining current role failed", err)
	}
}
