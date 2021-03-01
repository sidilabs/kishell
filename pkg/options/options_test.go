package options

import (
  "os"
  "testing"
)

func TestCLI(t *testing.T) {
  os.Args = []string{ "test", "list" }
  option := Parse()
  if option.Context == nil {
   t.Fatal("Option failed during creation")
  }
  option.Run()
}
