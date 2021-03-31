package options

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestAddServerConfig(t *testing.T) {
	serverName := "ut-server"
	var stdin bytes.Buffer
	stdin.Write([]byte(serverName + "\n")) // server name
	stdin.Write([]byte("http\n"))          // protocol
	stdin.Write([]byte("ut.test\n"))       // hostname
	stdin.Write([]byte("8080\n"))          // port
	stdin.Write([]byte("ut-user\n"))       // username
	stdin.Write([]byte("ut-passwd\n"))     // passwd
	stdin.Write([]byte("6.4.3\n"))         // kbn version
	stdin.Write([]byte("y\n"))             // is default?

	configuration := new(ConfigurationMock)
	configuration.On("GetStdin").Return(&stdin)
	configuration.On("AddServer", serverName, mock.Anything)
	configuration.On("GetServer").Return(serverName)
	configuration.On("SetServer", serverName)
	configuration.On("Save").Return(nil)

	context := Context{
		Debug:         true,
		Configuration: configuration,
	}

	cmd := ConfigureCmd{
		Server: true,
	}

	err := cmd.Run(&context)
	if err != nil {
		t.Fatal(err)
	}
	configuration.AssertExpectations(t)
	fmt.Printf(lineBreak)
}

func TestAddRoleConfig(t *testing.T) {
	roleName := "ut-role"

	var stdin bytes.Buffer
	stdin.Write([]byte(roleName + "\n")) // role name
	stdin.Write([]byte("ut-*\n"))        // index pattern
	stdin.Write([]byte("@timestamp\n"))  // window filter
	stdin.Write([]byte("y\n"))           // is default?

	configuration := new(ConfigurationMock)
	configuration.On("GetStdin").Return(&stdin)
	configuration.On("AddRole", roleName, mock.Anything)
	configuration.On("GetRole").Return(roleName)
	configuration.On("SetRole", roleName)
	configuration.On("Save").Return(nil)

	context := Context{
		Debug:         true,
		Configuration: configuration,
	}

	cmd := ConfigureCmd{
		Role: true,
	}
	err := cmd.Run(&context)
	if err != nil {
		t.Fatal(err)
	}
	configuration.AssertExpectations(t)
	fmt.Printf(lineBreak)
}

func TestInvalidConfigurationOption(t *testing.T) {
	context := Context{
		Debug: true,
	}
	cmd := ConfigureCmd{}

	err := cmd.Run(&context)
	if err == nil {
		t.Fatal(err)
	}
}
