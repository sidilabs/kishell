package options

import (
  "errors"
  "testing"
)

func TestRunListOptionSuccess(t *testing.T) {
  configuration := new(ConfigurationMock)
  configuration.On("PrettyPrint").Return(nil)

  context := Context{
   Debug:         true,
   Configuration: configuration,
  }
  cmd := ListCmd{}
  err := cmd.Run(&context)
  if err != nil {
    t.Fatal(err)
  }
  configuration.AssertExpectations(t)
}

func TestRunListOptionFailure(t *testing.T) {
  configuration := new(ConfigurationMock)
  configuration.On("PrettyPrint").Return(errors.New("mocked error response"))

  context := Context{
   Debug:         true,
   Configuration: configuration,
  }
  cmd := ListCmd{}
  err := cmd.Run(&context)
  if err == nil {
    t.Fatal(err)
  }
  configuration.AssertExpectations(t)
}
